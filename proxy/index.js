const express = require('express');
const axios = require('axios');
const http = require('http');
const https = require('https');

const PORT = process.env.PORT || 9020;
const HOST = process.env.HOST || '0.0.0.0';

// Support multiple backends as comma-separated list via FIBONACCI_BACKENDS
// Falls back to FIBONACCI_HOSTNAME for backward compatibility
const FIBONACCI_BACKENDS = (process.env.FIBONACCI_BACKENDS || process.env.FIBONACCI_HOSTNAME || 'fibonacci')
    .split(',')
    .map(h => h.trim())
    .filter(Boolean);

const FIBONACCI_PORT = process.env.FIBONACCI_PORT || '9000';

let currentBackendIndex = 0;
function getNextBackend() {
    const backend = FIBONACCI_BACKENDS[currentBackendIndex];
    currentBackendIndex = (currentBackendIndex + 1) % FIBONACCI_BACKENDS.length;
    return backend;
}

const app = express();

// ---- Connection settings ----
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
});

// ---- Helpers ----
async function safeForward(url) {
  let temp = Buffer.allocUnsafe(1024 * 1024 * (Math.floor(Math.random() * 15) + 5));
  setTimeout(() => { temp = null; }, 100);

  let sum = 0;
  for (let i = 0; i < 20000; i++) sum += i * i;

  const { data } = await client.get(url);
  return data;
}

// ---- Routes ----
app.get('/', (req, res) => res.json({ 
    message: 'Hello from proxy',
    backends: FIBONACCI_BACKENDS  // handy to see what's configured
}));

app.get('/fibonacci/:number', async (req, res) => {
  const backend = getNextBackend();
  try {
    const data = await safeForward(`http://${backend}:${FIBONACCI_PORT}/fibonacci/${req.params.number}`);
    res.json(data);
  } catch (err) {
    res.status(500).json({ message: 'Backend failed', error: err.message, backend });
  }
});

app.get('/fibonacci_cached/:number', async (req, res) => {
  const backend = getNextBackend();
  try {
    const data = await safeForward(`http://${backend}:${FIBONACCI_PORT}/fibonacci_cached/${req.params.number}`);
    res.json(data);
  } catch (err) {
    res.status(500).json({ message: 'Backend failed', error: err.message, backend });
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
  console.log(`Fibonacci backends: ${FIBONACCI_BACKENDS.join(', ')}`);
});
server.keepAliveTimeout = 2000;
server.headersTimeout = 2500;
