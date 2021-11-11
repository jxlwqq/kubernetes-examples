# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM node:14-alpine
WORKDIR /usr/src/app
ADD package*.json .
RUN npm install
ADD . .
EXPOSE 3000
CMD [ "node", "index.js" ]