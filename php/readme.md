## Installation
```bash
composer install
```

## run
```bash
php ws.php
```
## run options
```bash
--host=string (the host of the machine)
--port=int (the port of the service)
--rate=int (how many data per minute to be sent)
--worker=int (how many core worker)
```
## run example
```bash
php ws.php --host=localhost --port=8080 --rate=10 --worker=10
```
> open postman and try to connect with ws://{host}:{port}
> Known issues: need to wait around 10 seconds to be able to connect to the server -_-