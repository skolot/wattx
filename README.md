# wattx
WATTx Code Challenges

# How to run
## Dependicies

* lastes version of GoLang (may work on previous as well)
* docker + docker-compose
* make

### GoLang dependicies
* github.com/jinzhu/configor - configuaration utility, that support setting config from shell enviropment
* github.com/stretchr/testify/require - assertion library for tests

## Start/Stop application

to start application run next command in root of repository

```
$ make start
```

to stop

```
$ make stop
```

alternatively you can use docker-compose directly

start

```
$ docker-compose up --build --detach
```

stop

```
$ docker-compose down
```

## Services
* cctop - Ranking Service
* cmcprices - Pricing Service
* collector - HTTP-API Service

# How to use

Application bind http://127.0.0.1:2525 by default.<br>
HTTP API support next query parameters<br>
* limit - count of output lines (default 200)
* format - output formats, can be plain, csv or json (default plain)

example:

```
$ curl "http://localhost:2525?format=plain&limit=10" 
RANK	NAME	FULLNAME	PRICE	CURRENCY
1	BTC	Bitcoin	9558.551758	USD
2	ETH	Ethereum	168.688629	USD
3	ADA	Cardano	0.044599	USD
4	BCH	Bitcoin Cash	281.028992	USD
5	XRP	XRP	0.254449	USD
6	KNC	Kyber Network	0.162719	USD
7	LINK	Chainlink	1.814763	USD
8	ETC	Ethereum Classic	6.422846	USD
9	TRX	TRON	0.015601	USD
10	EOS	EOS	3.215316	USD
```

# Tests

For application testing I decided to use intergarion tests.<br>
Pros - we want check result, if some of components doesn't work properly test will fail.<br>
Cons of this method - we can't be certain what component cause error.<br>

## Run tests

Tests run with next command

```
$ make test
```

or manualy

```
$ cd test
$ docker-compose -p wattx_test -f ../docker-compose-test.yml up --build --detach
$ go test -v -count=1 main_test.go
$ docker-compose -p wattx_test -f ../docker-compose-test.yml down
```

# Problem and solution

Sometimes CoinMarketCap return error for some symbol values. To avoid it I pass option 'skip_invalid' to CoinMarketCap.
Sadly this option doesnt work in sandbox and API return error. Previous code that handle this error for sandbox can be found in commit e5a39091472b617f4cb2990336c81f8d47cbd0d1

# What to improve

* Better error handling between internal and foreign services
* Move all reusable code to library
* Cover error situation in tests
