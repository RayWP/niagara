const commander = require('commander');
const { Command } = require("commander");

const parser = new Command();

parser
    .name('string-util')
    .description('CLI to some JavaScript string utilities')
    .version('0.8.0')
    .option('-h, --host <type>', 'Specify the host', 'localhost')  // Optional with a default value
    .option('-p, --port <number>', 'Specify the port', '3000')
    .option('-r, --rate <number>', 'How many data per minute', '60');

module.exports = { parser };