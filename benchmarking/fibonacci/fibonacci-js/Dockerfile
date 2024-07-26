FROM node:18.16-slim

RUN apt-get update \
    && apt-get install curl --no-install-recommends -y \
    && rm -rf /var/lib/apt/lists/*

ENV PORT 8000
EXPOSE 8000

COPY package.json /opt/microservices/
COPY index.js /opt/microservices/
WORKDIR /opt/microservices

RUN npm install

CMD ["npm", "start"]
