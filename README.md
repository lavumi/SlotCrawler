# Slot Crawler

crawling tool for pragmaticplay slot game 

### version
- [go v1.19.5](https://go.dev/dl/)
- IDE [GoLand 2022.3.4](https://www.jetbrains.com/go/download/#section=windows)

### Install Go package manually
```bash
go get
```
```bash
//On windows, Set Enviromental Values ( GOPATH, GOBIN )
```

### Run Server
Need to set Environment variable
```shell
set MONGO_URL=mongo:27017;
set MONGO_USER=root;
set MONGO_PASS=example
```

```bash
go run cmd/main.go 
```

### Build Source

```bash
env GOOS=linux go build -o build/slot-crawler cmd/main.go
```

```
//vscode & widnows
go env -w GOOS=linux
go build -o build/slot-cralwer cmd/main.go
```


### Export Docker Image
```bash
docker build -t slot-crawler:x.y.z .
docker save -o slot-crawler.tar slot-crawler
```

### Run Docker
```bash
docker build -t slot-crawler:x.y.z .
docker compose up -d
```

### Execution URL
```bash
http://localhost:{port}
http://localhost:{port}/download.html
```

### Dev Mode 
```bash
//Docker가 아닌 자체 실행 모드
//windows & vscode 기반, terminal window에서 crawling 결과 값 출력 
go env -w GOOS=windows
go run ./cmd/crawler.go 
```
