const http = require('http');


const server = http.createServer((req, res) => {
    console.log('Request received on Server 1');
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('Hello from Server 1\n');
});

server.listen(9001, () => {
    console.log('Server 1 is listening on port 9001');
});