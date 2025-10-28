const http = require('http');


const server = http.createServer((req, res) => {
    console.log('Request received on Server 3');
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('Hello from Server 3\n');
});

server.listen(9003, () => {
    console.log('Server 3 is listening on port 9003');
});