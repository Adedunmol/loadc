# loadc
A tool to simulate load on your webiste or API.

# Usage Example

# Compile the project
```console
$ cd loadc
$ go run main.go
```
# Make requests sequentially
```console
$ ./main.exe -u https://google.com/ -n 10
Results:
 Total Requests  (2xx)..........:  10
 Failed Requests (5xx)..........:  0

Total request time (min, max, mean)...:  0.34 0.96 0.65
```
# Make requests concurrently
```console
$ ./main.exe -u https://google.com/ -n 10 -c 2
Results:
 Total Requests  (2xx)..........:  10
 Failed Requests (5xx)..........:  0

Total request time (min, max, mean)...:  0.34 0.96 0.65
```