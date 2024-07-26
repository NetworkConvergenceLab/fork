package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	app := fiber.New()

	cached_data := map[int]int{
		0:  0,
		1:  1,
		2:  1,
		3:  2,
		4:  3,
		5:  5,
		6:  8,
		7:  13,
		8:  21,
		9:  34,
		10: 55,
		11: 89,
		12: 144,
		13: 233,
		14: 377,
		15: 610,
		16: 987,
		17: 1597,
		18: 2584,
		19: 4181,
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{
			"hello": "world",
		})
	})

	app.Get("/fibonacci/:number", func(c *fiber.Ctx) error {
		number, err := strconv.Atoi(c.Params("number"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid number")
		}

		result := fibonacci(number)

		return c.JSON(map[string]interface{}{
			"message": result,
		})
	})

	app.Get("/fibonacci_cached/:number", func(c *fiber.Ctx) error {
		number, err := strconv.Atoi(c.Params("number"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid number")
		}

		if val, ok := cached_data[number]; ok {
			return c.JSON(map[string]interface{}{
				"message": val,
			})
		} else {
			result := fibonacci(number)
			cached_data[number] = result
			return c.JSON(map[string]interface{}{
				"message": result,
			})
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(fmt.Sprintf("0.0.0.0:%s", port))
}
