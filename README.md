# loadc
A tool to simulate load on your website or API.

# Usage Example

# Compile the project
```console
$ cd loadc
$ go build main.go
```

# Make requests sequentially
```console
$ ./main.exe -u <url> -n <number of requests>
Results:
 Total Requests  (2xx)..........:  10
 Failed Requests (5xx)..........:  0

Total request time (min, max, mean)...:  0.34 0.96 0.65
```

# Make requests concurrently
```console
$ ./main.exe -u <url> -n <number of requests> -c <number of concurrent requests>
Results:
 Total Requests  (2xx)..........:  10
 Failed Requests (5xx)..........:  0

Total request time (min, max, mean)...:  0.34 0.96 0.65
```

# Read URLs from a file
```console
$ ./main.exe -f <file path> -n <number of requests> -c <number of concurrent requests>
Results:
 Total Requests  (2xx)..........:  10
 Failed Requests (5xx)..........:  0

Total request time (min, max, mean)...:  0.34 0.96 0.65
```