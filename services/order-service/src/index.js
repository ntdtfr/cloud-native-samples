// Modules dependencies (Express, File System)
const express = require('express');
const fs = require('fs');

// Express app
const app = express();

// HTTP port from NODE_PORT environment variable
const DEFAULT_PORT = 8080;
const port = parseInt(process.env.NODE_PORT || DEFAULT_PORT);

// Say Hello
app.get('/', (req, res) => {
  console.log('Hello!!!')
});

// start HTTP server
app.listen(port, () => {
  console.log('Example app listening on port ' + port + '!')
});
