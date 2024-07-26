package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	plot "github.com/tsenart/vegeta/v12/lib/plot"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "Hello from api",
		})
	})

	attackUrl := os.Getenv("ATTACK_URL")
	attackDuration := os.Getenv("ATTACK_DURATION")
	attackFreq := os.Getenv("ATTACK_FREQ")

	app.Post("/attack/:mode", func(c *fiber.Ctx) error {
		freq, err := strconv.Atoi(attackFreq)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Freq not correct")
		}

		duration_sec, err := strconv.Atoi(attackDuration)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Duration not correct")
		}

		var port string
		var title string
		if c.Params("mode") == "single" {
			port = "31555"
			title = "Single Cluster Plot"
		} else {
			port = "31554"
			title = "Multi Cluster Plot"
		}

		rate := vegeta.Rate{Freq: freq, Per: time.Second}
		duration := time.Duration(duration_sec) * time.Second
		targeter := vegeta.NewStaticTargeter(vegeta.Target{
			Method: "GET",
			URL:    fmt.Sprintf("%s:%s/fibonacci/25", attackUrl, port),
		})
		attacker := vegeta.NewAttacker()

		p := plot.New(plot.Title(title), plot.Downsample(4000))
		for res := range attacker.Attack(targeter, rate, duration, title) {
			p.Add(res)
		}

		filename := fmt.Sprintf("%s.html", time.Now().Format("20060102150405"))
		f, err := os.Create(filename)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "File creation error")
		}

		p.WriteTo(f)

		return c.JSON(map[string]string{
			"message": "Plot created successfully",
			"file":    filename,
		})
	})

	app.Get("/plot/:filename", func(c *fiber.Ctx) error {
		return c.SendFile(c.Params("filename"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(fmt.Sprintf("0.0.0.0:%s", port))
}
