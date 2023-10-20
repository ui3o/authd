package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	LoginAccept struct {
		Url string `json:"url"`
	}

	Space struct {
		Name string `json:"cookie"`
	}

	Users struct {
		Store []string
	}
)

var (
	XForwardedUriCookieName  = "Auth-X-Forwarded-Uri"
	SERVER_PORT              = os.Getenv("AUTH_D_SERVER_PORT")
	FE_HTML                  = os.Getenv("AUTH_D_FE_HTML_PATH")
	COOKIE_TIMEOUT           = os.Getenv("AUTH_D_COOKIE_TIMEOUT")
	PASS_JSON_PATH           = os.Getenv("AUTH_D_PASS_JSON_PATH")
	COOKIE_NAME              = os.Getenv("AUTH_D_COOKIE_NAME")
	LOGIN_CHECK_URL          = os.Getenv("AUTH_D_LOGIN_CHECK_URL")
	REDIRECT_TO_LOGIN_URL    = os.Getenv("AUTH_D_REDIRECT_TO_LOGIN_URL")
	REDIRECT_AFTER_LOGIN_URL = os.Getenv("AUTH_D_REDIRECT_AFTER_LOGIN_URL")
	DEBUG                    = os.Getenv("AUTH_D_DEBUG")
)

func setDefault(val *string, defaultVal string) {
	if len(*val) == 0 {
		*val = defaultVal
	}
}

func Init() {
	setDefault(&SERVER_PORT, "8080")
	setDefault(&FE_HTML, "assets/index.html")
	setDefault(&COOKIE_TIMEOUT, "20")
	setDefault(&PASS_JSON_PATH, "pass.json")
	setDefault(&COOKIE_NAME, "auth_d_info")
	setDefault(&LOGIN_CHECK_URL, "/switch-space")
	setDefault(&REDIRECT_TO_LOGIN_URL, "/login")
	setDefault(&REDIRECT_AFTER_LOGIN_URL, "")
	setDefault(&DEBUG, "")
}

func debugLog(v ...any) {
	if len(DEBUG) > 0 {
		log.Println(v[:])
	}
}

func isUserValid(user, pass string) bool {
	var users Users

	b, err := os.ReadFile(PASS_JSON_PATH)
	if err != nil {
		debugLog(err)
	}
	err = json.Unmarshal(b, &users)

	if err == nil {
		for _, v := range users.Store {
			p, e := base64.StdEncoding.DecodeString(v)
			if e == nil {
				u := strings.Split(string(p[:]), ":")
				if u[0] == user && u[1] == pass {
					return true
				}
			}

		}
	}
	return false
}

func main() {
	Init()
	e := echo.New()

	loginGet := func(c echo.Context) (err error) {
		s := new(Space)

		if c.Request().Header["X-Forwarded-Uri"] != nil {
			uri := c.Request().Header["X-Forwarded-Uri"][0]
			debugLog("loginGet X-Forwarded-Uri", uri)

			cookie := new(http.Cookie)
			cookie.Name = XForwardedUriCookieName
			cookie.Value = uri[1:]
			cookie.HttpOnly = true
			cookie.Secure = true
			c.SetCookie(cookie)
		}

		cookie, err := c.Cookie(COOKIE_NAME)
		if err == nil {
			s.Name = cookie.Value
			debugLog("cookie found", s.Name)
			return c.JSON(http.StatusOK, s)
		} else {
			debugLog("no cookie set")
			if c.QueryParam("auth") == "fe" {
				return c.JSON(http.StatusOK, s)
			} else {
				return c.Redirect(http.StatusMovedPermanently, REDIRECT_TO_LOGIN_URL)
			}
		}
	}

	loginPost := func(c echo.Context) (err error) {
		cookieBody := []string{}
		url := ""
		user := ""
		pass := ""
		my_data := echo.Map{}
		if err := c.Bind(&my_data); err != nil {
			return err
		} else {
			for k, v := range my_data {
				// debugLog("type: ", reflect.TypeOf(v), " k ", k, " v ", fmt.Sprintf("%v", v))
				switch v.(type) {
				case string:
					if k == "url" {
						url = fmt.Sprintf("%v", v)
					}
				default:
					field := v.(map[string]interface{})
					boolValue, err := strconv.ParseBool(fmt.Sprintf("%v", field["skip"]))
					v, e := base64.StdEncoding.DecodeString(fmt.Sprintf("%v", field["v"]))

					if err == nil && !boolValue {
						if e == nil {
							cookieBody = append(cookieBody, k+":"+string(v[:]))
						}
					}
					if k == "p" && e == nil {
						pass = string(v[:])
					}
					if k == "u" && e == nil {
						user = string(v[:])
					}
				}
			}
		}

		// https cookie
		cookie := new(http.Cookie)
		cookie.Name = COOKIE_NAME
		cookie.Value = strings.Join(cookieBody, "_")
		timeout, err := strconv.Atoi(COOKIE_TIMEOUT)
		if err != nil {
			panic(err)
		}
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.Expires = time.Now().Add(time.Duration(timeout) * time.Second)
		c.SetCookie(cookie)
		debugLog(cookie)

		l := new(LoginAccept)
		if len(REDIRECT_AFTER_LOGIN_URL) > 0 {
			l.Url = REDIRECT_AFTER_LOGIN_URL
		} else {
			uriCooke, err := c.Cookie(XForwardedUriCookieName)
			if err == nil {
				l.Url = uriCooke.Value
			} else {
				l.Url = url
			}
		}
		debugLog("json response", l)
		if isUserValid(user, pass) {
			return c.JSON(http.StatusOK, l)
		} else {
			return c.JSON(http.StatusUnauthorized, nil)
		}
	}

	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	e.HideBanner = true

	u, err := url.Parse(REDIRECT_TO_LOGIN_URL)
	if err != nil {
		debugLog("can not get path from ", REDIRECT_TO_LOGIN_URL)

	}

	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		// res := c.Response()
		switch req.Method {
		case "POST":
			switch req.RequestURI {
			case LOGIN_CHECK_URL:
				return loginPost(c)
			}
		case "GET":
			switch req.URL.Path {
			case "/":
			case "/index.html":
			case u.Path:
				b, err := os.ReadFile(FE_HTML)
				if err != nil {
					debugLog(err)
				}
				c.HTML(http.StatusOK, string(b))
			case LOGIN_CHECK_URL:
				return loginGet(c)
			}

		}
		return
	})
	e.Logger.Fatal(e.Start(":" + SERVER_PORT))
}
