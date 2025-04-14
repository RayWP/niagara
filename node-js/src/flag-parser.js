const { Command } = require("commander");

const parser = new Command();

parser
    .name('niagara')
    .description('WebSocket server for stock trading data simulation')
    .version('1.0.0')
    .option('-h, --host <string>', 'host address of the machine', 'localhost')
    .option('-p, --port <number>', 'port number for the service', '8080')
    .option('-r, --rate <number>', 'how many data messages per minute', '60')
    .option('-w, --worker <number>', 'how many core workers (not implemented)', '1')
    .showHelpAfterError()
    .addHelpText('after', `
Example:
  $ node ws.js --host=localhost --port=8080 --rate=10`);

module.exports = { parser };