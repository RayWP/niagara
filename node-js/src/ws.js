const WebSocket = require('ws');
const { parser } = require('./flag-parser');
const { DataHub } = require('./DataHub');

try {
    // Parse arguments with process.argv
    parser.parse(process.argv);
    const options = parser.opts();

    const port = parseInt(options.port) || 8080;
    const host = options.host || 'localhost';
    const rate = parseInt(options.rate) || 60;

    if (isNaN(port) || port < 0 || port > 65535) {
        throw new Error('Port must be a number between 0 and 65535');
    }

    console.log('Server configuration:');
    console.log(`Host: ${host}`);
    console.log(`Port: ${port}`);
    console.log(`Rate: ${rate} messages per minute`);

    const wss = new WebSocket.Server({
        port,
        host
    });

    const dataHub = new DataHub(rate);

    wss.on('connection', (ws) => {
        console.log('Client connected');
        dataHub.addConnection(ws);
        
        ws.on('close', () => {
            console.log('Client disconnected');
            dataHub.removeConnection(ws);
        });

    });
    
    console.log(`WebSocket server is running on ws://${host}:${port}`);
} catch (error) {
    console.error('Error:', error.message);
    process.exit(1);
}

