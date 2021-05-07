FROM node:15.14.0-alpine

WORKDIR /usr/app

COPY package*.json ./
RUN npm ci -qy

COPY . .

EXPOSE 8080

CMD ["npm", "start"]
