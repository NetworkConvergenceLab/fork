# Copyright 2023-2024 Shahmir Ejaz

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

cache_data = {
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

class FibonacciResult(BaseModel):
    message: str

async def fibonacci(number: int) -> int:
    if number == 0:
        return 0
    elif number == 1:
        return 1
    else:
        r1 = await fibonacci(number - 1)
        r2 = await fibonacci(number - 2)
        return r1 + r2

@app.get("/")
async def root():
    return {"message": "Hello World"}

@app.get("/fibonacci/{number}", response_model=FibonacciResult)
async def get_fibonacci(number: int):
    result = await fibonacci(number)
    return {"message": str(result)}

@app.get("/fibonacci_cached/{number}", response_model=FibonacciResult)
async def get_fibonacci_cached(number: int):
    if number in cache_data:
        return {"message": str(cache_data[number])}
    else:
        result = await fibonacci(number)
        return {"message": str(result)}
