const express = require('express');
const axios = require('axios');
const http = require('http');
const https = require('https');

const PORT = process.env.PORT || 9020;
const HOST = process.env.HOST || '0.0.0.0';

const FIBONACCI_HOSTNAME = process.env.FIBONACCI_HOSTNAME || 'fibonacci';
const FIBONACCI_PORT = process.env.FIBONACCI_PORT || '9000';

const app = express();

// ---- Connection settings ----
// Disable keepAlive so cross-cluster always pays connect cost
const agentOptions = {
  keepAlive: false,
  maxSockets: 50,
  timeout: 5000
};
const httpAgent = new http.Agent(agentOptions);
const httpsAgent = new https.Agent(agentOptions);

const client = axios.create({
  httpAgent,
  httpsAgent,
  // ❌ remove validateStatus => true so axios throws on error
});

// ---- Helpers ----
async function safeForward(url) {
  // Oscillating memory: 5–20 MB
  let temp = Buffer.allocUnsafe(1024 * 1024 * (Math.floor(Math.random() * 15) + 5));
  setTimeout(() => { temp = null; }, 100);

  // CPU jitter
  let sum = 0;
  for (let i = 0; i < 20000; i++) sum += i * i;

  // Forward request and let real latency show
  const { data } = await client.get(url);
  return data;
}

// ---- Routes ----
app.get('/', (req, res) => res.json({ message: 'Hello from proxy' }));

app.get('/fibonacci/:number', async (req, res) => {
  try {
    const data = await safeForward(`http://${FIBONACCI_HOSTNAME}:${FIBONACCI_PORT}/fibonacci/${req.params.number}`);
    res.json(data);
  } catch (err) {
    res.status(500).json({ message: 'Backend failed', error: err.message });
  }
});

app.get('/fibonacci_cached/:number', async (req, res) => {
  try {
    const data = await safeForward(`http://${FIBONACCI_HOSTNAME}:${FIBONACCI_PORT}/fibonacci_cached/${req.params.number}`);
    res.json(data);
  } catch (err) {
    res.status(500).json({ message: 'Backend failed', error: err.message });
  }
});

// ---- Background oscillation ----
if (global.gc) {
  setInterval(() => {
    let junk = Buffer.allocUnsafe(1024 * 1024 * (Math.floor(Math.random() * 20) + 30));
    setTimeout(() => { junk = null; global.gc(); }, 500);
  }, 3000);
}

const server = app.listen(PORT, HOST, () => {
  console.log(`Running on http://${HOST}:${PORT}`);
});
server.keepAliveTimeout = 2000;
server.headersTimeout = 2500;
