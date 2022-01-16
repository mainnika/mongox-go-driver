# mongox-go-driver

## testing locally

reqs:
- mongodb v4.0 or newer run on localhost
- golang v1.13 or newer

test it by calling go tests
```sh
$ go test ./...
```

## testing by using dockerfile

reqs:
- docker with buildkit

```sh
$ DOCKER_BUILDKIT=1 docker build -t mongox-testing -f testing.Dockerfile .
```