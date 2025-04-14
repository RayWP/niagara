const WebSocket = require('ws');
const { parser } = require('./flag-parser');

const options = parser.parse();

const wss = new WebSocket.Server(
    { 
        port: 3000, 
    });

wss.on('connection', (ws) => {
    console.log('Client connected');

    ws.on('message', (message) => {
        console.log(`Received message: ${message}`);
        // Handle incoming message here
    });

    ws.on('close', () => {
        console.log('Client disconnected');
        // Handle client disconnection here
    });
});

