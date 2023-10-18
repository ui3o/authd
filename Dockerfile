FROM docker.io/golang:1.21 as golang
COPY . /opt
WORKDIR /opt
RUN go version && go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w' -o authd *.go

FROM docker.io/node:20-alpine as node
COPY . /opt
WORKDIR /opt/html
RUN yarn && yarn run fix && yarn run build

FROM alpine:latest

ENV AUTH_D_FE_HTML_PATH=/opt/index.html
ENV AUTH_D_PASS_JSON_PATH=/opt/pass.json

COPY pass.json /opt/pass.json
COPY --from=node /opt/assets/index.html /opt/index.html
COPY --from=golang /opt/authd /opt/authd

ENTRYPOINT ["/opt/authd"]