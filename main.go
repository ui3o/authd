package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	SERVER_PORT    = os.Getenv("AUTH_D_SERVER_PORT")
	FE_HTML        = os.Getenv("AUTH_D_FE_HTML_PATH")
	COOKIE_TIMEOUT = os.Getenv("AUTH_D_COOKIE_TIMEOUT")
	PASS_JSON_PATH = os.Getenv("AUTH_D_PASS_JSON_PATH")
	COOKIE_NAME    = os.Getenv("AUTH_D_COOKIE_NAME")
)

func setDefault(val *string, defaultVal string) {
	if len(*val) == 0 {
		*val = "8080"
	}
}

func Init() {
	setDefault(&SERVER_PORT, "8080")
	setDefault(&FE_HTML, "assets/index.html")
	setDefault(&COOKIE_TIMEOUT, "20")
	setDefault(&PASS_JSON_PATH, "pass.json")
	setDefault(&COOKIE_NAME, "auth_d_info")
}

func isUserValid(user, pass string) bool {
	var users Users

	b, err := os.ReadFile(PASS_JSON_PATH)
	if err != nil {
		log.Println(err)
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

	e.GET("/switch-space", func(c echo.Context) (err error) {
		// todo redirect
		s := new(Space)
		cookie, err := c.Cookie(COOKIE_NAME)
		if err == nil {
			s.Name = cookie.Value
		}
		return c.JSON(http.StatusOK, s)
	})

	e.POST("/switch-space", func(c echo.Context) (err error) {
		cookieBody := []string{}
		url := ""
		user := ""
		pass := ""
		my_data := echo.Map{}
		if err := c.Bind(&my_data); err != nil {
			return err
		} else {
			for k, v := range my_data {
				// log.Println("type: ", reflect.TypeOf(v), " k ", k, " v ", fmt.Sprintf("%v", v))
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

		cookie := new(http.Cookie)
		cookie.Name = COOKIE_NAME
		cookie.Value = strings.Join(cookieBody, "_")
		timeout, err := strconv.Atoi(COOKIE_TIMEOUT)
		if err != nil {
			panic(err)
		}
		cookie.Expires = time.Now().Add(time.Duration(timeout) * time.Second)
		c.SetCookie(cookie)

		l := new(LoginAccept)
		l.Url = url
		if isUserValid(user, pass) {
			return c.JSON(http.StatusOK, l)
		} else {
			return c.JSON(http.StatusUnauthorized, nil)
		}
	})

	e.File("/", FE_HTML)
	e.File("/index.html", FE_HTML)

	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	e.HideBanner = true
	e.Logger.Fatal(e.Start(":" + SERVER_PORT))
}
