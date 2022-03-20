# myhttp
Tool makes http requests and prints the address of the request along with the MD5 hash of the response.

## Build tool
```
go build ./myhttp.go
```

## Run tool
```
./myhttp <host1> <host2> <host3>
```
Example:
```
./myhttp google.com facebook.com yahoo.com
```

## Parallel 
Tool perform requests in parallel. To limit the number of parallel requests use flag parallel. Default limit is 10.

```
./myhttp -parallel <parallel limit> <host1> <host2> <host3>
```
Example:
```
./myhttp -parallel 3 google.com facebook.com yahoo.com
```