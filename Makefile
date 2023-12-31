# VERSION ?= 0.3.

build-go:
	@echo [INFO] start build go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w' -o authd *.go
	@echo [DONE] end build go


build-local: build-go
	@echo [INFO] start build local 
	cd html && yarn && yarn run fix && yarn run build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w' -o authd *.go
	@echo [DONE] end build local

cert:
	openssl genrsa 2048 | openssl pkcs8 -topk8 -inform PEM -out id_rsa
	openssl rsa -in id_rsa -pubout -out id_rsa.pub