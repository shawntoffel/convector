# Convector ☁

Sends GET requests to specified endpoints on a recurring interval with jitter. An example use case may be to avoid cold starts of serverless applications. 

## Usage:
```
Usage of convector:
  -d    enable debug logging
  -e string
        comma-delimited list of endpoints
  -i duration
        base interval between requests (default 10m0s)
  -j duration
        standard deviation for a normal distribution jitter between requests (default 5m0s)
  -t duration
        http client timeout for requests (default 10s)
  -v    print convector version
```

## Example:
```sh
./convector -d -i 10m -j 5m -e "https://myip.shawntoffel.com"
```
