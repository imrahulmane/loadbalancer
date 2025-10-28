const http = require('http');


const server = http.createServer((req, res) => {
    console.log('Request received on Server 2');
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('Hello from Server 2\n');
});

server.listen(9002, () => {
    console.log('Server 2 is listening on port 9002');
});