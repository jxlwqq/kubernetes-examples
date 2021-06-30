FROM node:14
WORKDIR /usr/src/app
ADD package*.json .
RUN npm install
ADD . .
EXPOSE 3000
CMD [ "node", "index.js" ]