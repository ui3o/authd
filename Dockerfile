FROM docker.io/golang:1.21 as golang
COPY . /opt/authd
WORKDIR /opt/authd
RUN go version && go get ./...
RUN ls -al
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w' -o authd *.go

FROM docker.io/node:20-alpine as node
COPY . /opt/authd
WORKDIR /opt/authd/html
RUN rm -rf ../assets && ls -al ../
RUN npm install --global yarn
RUN yarn && yarn run fix && yarn run build

FROM alpine:latest

ENV AUTH_D_FE_HTML_PATH=/opt/index.html
ENV AUTH_D_PASS_JSON_PATH=/opt/pass.json

COPY pass.json /opt/pass.json
COPY --from=node /opt/authd/assets/index.html /opt/index.html
COPY --from=golang /opt/authd /opt/authd

ENTRYPOINT ["/opt/authd"]