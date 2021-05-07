'use strict';

const express = require('express');
const path = require('path');

// Constants
const PORT = process.env.PORT || 8080;
const HOST = '0.0.0.0';

const CLIENT_BUILD_PATH = path.join(__dirname, '../../client/build');

// App
const app = express();
const redis = require('redis');
const client = redis.createClient({host:'redis'});

// Static files
app.use(express.static(CLIENT_BUILD_PATH));
app.use(express.json());
app.use(express.urlencoded({ extended: true}));

// API
app.get('/product', (req, res) => {
  res.set('Content-Type', 'application/json');
  const data = client.hkeys('item1', function(err, values){
      res.send(JSON.stringify(values, null, 2));
  })
});

app.post('/product', (req, res) => {
  console.log("body =", req.body)
  res.set('Content-Type', 'application/json');
  client.hset('item1', req.body.title, req.body.metadata, redis.print)
  res.send(JSON.stringify(null, 2));
});

// All remaining requests return the React app, so it can handle routing.
app.get('*', function(request, response) {
  response.sendFile(path.join(CLIENT_BUILD_PATH, 'index.html'));
});

app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
