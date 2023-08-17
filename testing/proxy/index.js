const express = require('express');
const axios = require('axios');

const PORT = process.env.PORT || 9020;
const HOST = process.env.HOST || '0.0.0.0';

const DETAILS_HOSTNAME = process.env.DETAILS_HOSTNAME || 'details';
const DETAILS_PORT = process.env.DETAILS_PORT || '9080';
const REVIEWS_HOSTNAME = process.env.REVIEWS_HOSTNAME || 'reviews';
const REVIEWS_PORT = process.env.REVIEWS_PORT || '9080';
const RATINGS_HOSTNAME = process.env.RATINGS_HOSTNAME || 'ratings';
const RATINGS_PORT = process.env.RATINGS_PORT || '9080';

const app = express();

app.get('/', (req, res) => {
    res.json({
        message: 'Hello from proxy'
    });
});

app.get('/details', async (req, res) => {
    try {
        console.log(`Accessing route: http://${DETAILS_HOSTNAME}:${DETAILS_PORT}/details/0`);
        const { data } = await axios.get(`http://${DETAILS_HOSTNAME}:${DETAILS_PORT}/details/0`);
        res.json(data);
    }
    catch (error) {
        console.log(error);
        res.status(500).json({
            message: 'Something went wrong'
        });
    }
});

app.get('/reviews', async (req, res) => {
    try {
        console.log(`Accessing route: http://${REVIEWS_HOSTNAME}:${REVIEWS_PORT}/reviews/0`);
        const { data } = await axios.get(`http://${REVIEWS_HOSTNAME}:${REVIEWS_PORT}/reviews/0`);
        res.json(data);
    } catch (error) {
        console.log(error);
        res.status(500).json({
            message: 'Something went wrong'
        });
    }
})

app.get('/ratings', async (req, res) => {
    try {
        console.log(`Accessing route: http://${RATINGS_HOSTNAME}:${RATINGS_PORT}/ratings/0`);
        const { data } = await axios.get(`http://${RATINGS_HOSTNAME}:${RATINGS_PORT}/ratings/0`);
        res.json(data);
    } catch (error) {
        console.log(error);
        res.status(500).json({
            message: 'Something went wrong'
        });
    }
})

app.listen(PORT, HOST, () => {
    console.log(`Running on http://${HOST}:${PORT}`);
});
