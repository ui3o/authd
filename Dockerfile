FROM alpine:latest

ENV AUTH_D_FE_HTML_PATH=/opt/index.html
ENV AUTH_D_PASS_JSON_PATH=/opt/pass.json

COPY pass.json /opt/pass.json
COPY assets/index.html /opt/index.html
COPY authd /opt/authd

ENTRYPOINT ["/opt/authd"]