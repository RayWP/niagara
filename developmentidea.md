# Language target
- C++
- C#
- GoLang
- Java
- Php
- Php-Laravel

# Specs
## 1. This software should have several flags for condition<br>
### Rate Flag

--rate(int) <br>
how many data should be send in one minute <br>
ignored if random-rate flag is `true`

### Random-rate flag
--random-rate(bool) <br>
randomnized the number of data we send in one minute

## 2. Supported ticks dan format
### Data Format to be sent
> StockCode|Price|B/S|Amount <br>
> eg. AAAA|100|B|1
