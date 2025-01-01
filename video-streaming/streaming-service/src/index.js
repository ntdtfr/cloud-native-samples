const express = require('express');
const fs = require('fs');

const app = express();

const DEFAULT_PORT = 3000;

if (!process.env.PORT) {
  process.env.PORT = DEFAULT_PORT;
}

const PORT = parseInt(process.env.PORT);

app.get('/video', (req, res) => {
  //res.send('Hello World!')
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

app.listen(PORT, () => {
  console.log('Example app listening on port ' + PORT + '!')
});
