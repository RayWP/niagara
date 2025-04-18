# Language target
- C++
- C#
- GoLang
- Java
- [NodeJS](./node-js/)
- [Php](./php/)
- Php-Laravel
- Rust

# Specs
## 1. This software should have several flags for condition
### Host Flag
--host(string) <br> (default: localhost)
the host address of the machine

### Port Flag
--port(int) <br> (default: 8080)
the port number for the service

### Rate Flag
--rate(int) <br> (default: 60)
how many data should be send in one minute <br>
ignored if random-rate flag is `true`

### Worker Flag
--worker(int) <br>
how many core workers to use


## 2. Supported ticks dan format
### Data Format to be sent
> StockCode|Price|B/S|Amount <br>
> eg. AAAA|100|B|1

# future feature
## Chaos (randomly increase or decrease rate)
