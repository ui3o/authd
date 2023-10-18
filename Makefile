# VERSION ?= 0.3.

build-local:
	@echo [INFO] start build 
	cd html && yarn && yarn run fix && yarn run build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w' -o authd *.go
	@echo [DONE] end build



