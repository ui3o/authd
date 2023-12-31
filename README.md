# intro

Simple authentication daemon with **golang** and **reactjs**. The password store is a json file.

# environment variable for configuration

* **AUTH_D_SERVER_PORT**: set the sever listener port (default: 8080)
* **AUTH_D_FE_HTML_PATH**: single html file path (default: ./assets/index.html)
* **AUTH_D_COOKIE_TIMEOUT**: authentication cookie timeout in sec (default: 20)
* **AUTH_D_PASS_JSON_PATH**: path to password.json (default: ./pass.json)
* **AUTH_D_COOKIE_NAME**: cookie name for accept authentication (default: auth_d_info)
* **AUTH_D_LOGIN_CHECK_URL**: url path which response 200 if user logged in (default: /switch-space)
* **AUTH_D_REDIRECT_TO_LOGIN_URL**: url path where authd redirect if user not logged in (default: /login)
* **AUTH_D_REDIRECT_AFTER_LOGIN_URL**: override react page url with this env variable where authd redirect user after logged in (default: empty)
* **AUTH_D_DEBUG**: add debug option for debug server (default: empty) 

# pass generation

Example in **pass.json**. For generation use this command: `echo -n "bar:foo" | base64`
