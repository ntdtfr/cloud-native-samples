// Modules dependencies (Express, File System)
const express = require('express');
const fs = require('fs');

// Express app
const app = express();

// HTTP port from NODE_PORT environment variable
const DEFAULT_PORT = 3000;
const port = parseInt(process.env.NODE_PORT || DEFAULT_PORT);

// serve static video file
app.get('/video', (req, res) => {
  const path = "./videos/file_example_MP4_480_1_5MG.mp4";
  fs.stat(path, (err, stats) => {
    if (err) {
      console.error(err);
      res.sendStatus(500);
      return;
    }
    res.writeHead(200, {
      'Content-Length': stats.size,
      'Content-Type': 'video/mp4',
    });
    fs.createReadStream(path).pipe(res);
  })
});

// start HTTP server
app.listen(port, () => {
  console.log('Example app listening on port ' + port + '!')
});
