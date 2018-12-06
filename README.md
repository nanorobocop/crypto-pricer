## Features

* Collector
  * Collect prices from crypo currency exchanges
  * Support for multiple crypto currencies
  * Support for multiple exchanges
  * Support for separate schedule for each exchange
  * Permanent store data in database

* Web server
  * Show price difference between exchanges
  * Support custom period and threshold in queries

## Technical Features

* Common interface for exchanges (easy to add new exchanges and currencies)
* Concurrent requests to exchanges
* Go modules for dependencies management
* Store data in Mongodb
* Agregation pipelines in Mongodb for fast queries
* Lightweight [gin](https://github.com/gin-gonic/gin) web framework

## Run locally

```bash
$ docker-compose up --build
```

## Usage

### Check collector logs

```bash
collector_1_3856a8a926a6 | time="2018-12-06T14:29:01Z" level=info msg="Fetching price" exchange=bitflyer pair=BTC_JPY
collector_1_3856a8a926a6 | time="2018-12-06T14:29:01Z" level=info msg="Fetching price" exchange=zaif pair=btc_jpy
collector_1_3856a8a926a6 | time="2018-12-06T14:29:01Z" level=info msg="Fetching price" exchange=zaif pair=eth_jpy
collector_1_3856a8a926a6 | time="2018-12-06T14:29:01Z" level=info msg="Fetching price" exchange=zaif pair=eth_btc
collector_1_3856a8a926a6 | time="2018-12-06T14:29:01Z" level=info msg="Fetching price" exchange=bitflyer pair=ETH_BTC
collector_1_3856a8a926a6 | time="2018-12-06T14:29:01Z" level=info msg="Commiting result" error="<nil>" exchange=bitflyer pair=BTC_JPY price=414313 timestamp="2018-12-06 14:29:01 +0000 UTC"
collector_1_3856a8a926a6 | time="2018-12-06T14:29:01Z" level=info msg="Commiting result" error="<nil>" exchange=bitflyer pair=ETH_BTC price=0.02635 timestamp="2018-12-06 14:29:01 +0000 UTC"
collector_1_3856a8a926a6 | time="2018-12-06T14:29:02Z" level=info msg="Commiting result" error="<nil>" exchange=zaif pair=eth_btc price=0.0265 timestamp="2018-12-06 14:29:02.047872 +0000 UTC m=+660.167063801"
collector_1_3856a8a926a6 | time="2018-12-06T14:29:02Z" level=info msg="Commiting result" error="<nil>" exchange=zaif pair=eth_jpy price=11155 timestamp="2018-12-06 14:29:02.0490139 +0000 UTC m=+660.168206001"
collector_1_3856a8a926a6 | time="2018-12-06T14:29:02Z" level=info msg="Commiting result" error="<nil>" exchange=zaif pair=btc_jpy price=415190 timestamp="2018-12-06 14:29:02.0638324 +0000 UTC m=+660.183024401"
```

### Check recent prices

http://localhost:8080

```bash
{"message":[{"id":"5c089246a4e653b37bcf5a85","exchange":"zaif","pair":"btc_jpy","price":419280,"datetime":"2018-12-06T12:06:46.741+09:00"},{"id":"5c08925da4e653b37bcf5a9a","exchange":"bitflyer","pair":"btc_jpy","price":419042,"datetime":"2018-12-06T12:07:12+09:00"},{"id":"5c08925da4e653b37bcf5a9c","exchange":"zaif","pair":"btc_jpy","price":419265,"datetime":"2018-12-06T12:07:09.526+09:00"},{"id":"5c08926ba4e653b37bcf5aa6","exchange":"zaif","pair":"btc_jpy","price":419000,"datetime":"2018-12-06T12:07:23.276+09:00"},{"id":"5c089287a4e653b37bcf5ab5","exchange":"bitflyer","pair":"btc_jpy","price":419163,"datetime":"2018-12-06T12:07:55+09:00"}]}
```

### Check diff between exchanges and alerts

http://localhost:8080/alerts

```bash
{"message":[{"pair":"eth_jpy","min":10970,"max":11170,"price_between_max_min":200,"datetime":"2018-12-06T22:39:50.268+09:00","alert":true},{"pair":"eth_btc","min":0.02635,"max":0.0265,"price_between_max_min":0.00015000000000000083,"datetime":"2018-12-06T22:39:50.268+09:00","alert":false},{"pair":"btc_jpy","min":407001,"max":410558,"price_between_max_min":3557,"datetime":"2018-12-06T22:39:50.268+09:00","alert":true}]}
```

### Set threshold

http://localhost:8080/alerts?threshold=1000

```bash
{"message":[{"pair":"eth_jpy","min":10970,"max":11170,"price_between_max_min":200,"datetime":"2018-12-06T22:40:47.612+09:00","alert":false},{"pair":"eth_btc","min":0.02635,"max":0.0265,"price_between_max_min":0.00015000000000000083,"datetime":"2018-12-06T22:40:47.612+09:00","alert":false},{"pair":"btc_jpy","min":407001,"max":410558,"price_between_max_min":3557,"datetime":"2018-12-06T22:40:47.612+09:00","alert":true}]}
```

### Set threshold and period

http://localhost:8080/alerts?threshold=1000&period=1h

```bash
{"message":[{"pair":"eth_jpy","min":10970,"max":11385,"price_between_max_min":415,"datetime":"2018-12-06T22:41:17.143+09:00","alert":false},{"pair":"eth_btc","min":0.02635,"max":0.0265,"price_between_max_min":0.00015000000000000083,"datetime":"2018-12-06T22:41:17.143+09:00","alert":false},{"pair":"btc_jpy","min":407001,"max":410558,"price_between_max_min":3557,"datetime":"2018-12-06T22:41:17.143+09:00","alert":true}]}
```
