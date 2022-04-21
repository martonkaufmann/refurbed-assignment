## Developement

Build docker image:
```
docker build --no-cache -f Dockerfile.dev -t refurbed/assignment-dev:latest .
```

Run dev container:
```
docker run -d -v [PATH]:/go/src/refurbed/assignment refurbed/assignment-dev:latest
```

Generating test data:
```
tr -dc "A-Za-z 0-9" < /dev/urandom | fold -w100|head -n 100000 >> messages.txt
```

Exec:
```
go build -o notify main.go
```
```
./notify -url http://httpbin.org/post < ../../messages.txt
```