# intro

Simple authentication daemon with **golang** and **reactjs**. The password store is a json file.

# environment variable for configuration

* **AUTH_D_SERVER_PORT**: set the sever listener port (default: 8080)
* **AUTH_D_FE_HTML_PATH**: single html file path (default: ./assets/index.html)
* **AUTH_D_COOKIE_TIMEOUT**: authentication cookie timeout in sec (default: 20)
* **AUTH_D_PASS_JSON_PATH**: path to password.json (default: ./pass.json)
* **COOKIE_NAME**: cookie name for accept authentication (default: auth_d_info)

# pass generation

Example in **pass.json**. For generation use this command: `echo -n "bar:foo" | base64`
