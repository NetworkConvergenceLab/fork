FROM node:18.16-slim

WORKDIR /usr/src/app

COPY package*.json ./
COPY index.js ./

RUN npm ci

EXPOSE 9020

CMD ["node", "./index.js"]
