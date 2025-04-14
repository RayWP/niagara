## Installation
```bash
npm install
```

## run
```bash
node src/ws.js --help
```
## run options
```bash
--host=string (the host of the machine)
--port=int (the port of the service)
--rate=int (how many data per minute to be sent)
--worker=int (how many core worker) (not implemented yet)
```
## run example
```bash
node src/ws.js --host=localhost --port=8080 --rate=10 --worker=10
```
> open postman and try to connect with ws://{host}:{port}
> Known issues: xx -_-