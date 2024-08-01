/*
Copyright 2023-2024 Shahmir Ejaz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import Fastify from 'fastify'

const port = process.env.PORT || 3000

const fastify = Fastify({
  logger: true
})

const cached_data = {
    0: 0,
    1: 1,
    2: 1,
    3: 2,
    4: 3,
    5: 5,
    6: 8,
    7: 13,
    8: 21,
    9: 34,
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

const fibonacci = async (number) => {
    if (number === 0) return 0
    if (number === 1) return 1

    const r1 = await fibonacci(number - 1)
    const r2 = await fibonacci(number - 2)
    return r1 + r2
}

fastify.get('/', async function handler (request, reply) {
  return { hello: 'world' }
})

fastify.get('/fibonacci/:number', async function handler (request, reply) {
    const number = request.params.number
    const result = await fibonacci(number)
    return { "message": result }
})

fastify.get('/fibonacci_cached/:number', async function handler (request, reply) {
    const number = request.params.number
    const result = cached_data[number]
    return { "message": result }
})

try {
  await fastify.listen({ host: '0.0.0.0', port: port })
} catch (err) {
  fastify.log.error(err)
  process.exit(1)
}
