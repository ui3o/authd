package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kardianos/osext"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type JWT struct {
	privateKeyPath string
	privateKey     []byte
	publicKey      []byte
}

func InitJwt() {
	filename, _ := osext.Executable()
	path := strings.Split(filename, "/")
	path = path[:len(path)-1]
	home, err := os.UserHomeDir()
	if err != nil {
		os.Exit(1)
	}

	AUTH_D_PRIV_IDRSA := os.Getenv("AUTH_D_PRIV_IDRSA")
	if len(AUTH_D_PRIV_IDRSA) == 0 {
		Jwt.privateKeyPath = "id_rsa"
		Jwt.privateKey, err = os.ReadFile(Jwt.privateKeyPath)
		if err != nil {
			Jwt.privateKeyPath = "/etc/authd/id_rsa"
			Jwt.privateKey, err = os.ReadFile(Jwt.privateKeyPath)
			if err != nil {
				Jwt.privateKeyPath = strings.Join(append(path, "id_rsa"), "/")
				Jwt.privateKey, err = os.ReadFile(Jwt.privateKeyPath)
				if err != nil {
					Jwt.privateKeyPath = home + "/.ssh/id_rsa"
					Jwt.privateKey, err = os.ReadFile(Jwt.privateKeyPath)
					if err != nil {
						os.Exit(1)
					}
				}
			}
		}
	} else {
		Jwt.privateKey = []byte(AUTH_D_PRIV_IDRSA)
	}

	AUTH_D_PUB_IDRSA := os.Getenv("AUTH_D_PUB_IDRSA")
	if len(AUTH_D_PUB_IDRSA) == 0 {
		Jwt.publicKey, err = os.ReadFile("id_rsa.pub")
		if err != nil {
			Jwt.publicKey, err = os.ReadFile("/etc/authd/id_rsa.pub")
			if err != nil {
				Jwt.publicKey, err = os.ReadFile(strings.Join(append(path, "id_rsa.pub"), "/"))
				if err != nil {
					Jwt.publicKey, err = os.ReadFile(home + "/.ssh/id_rsa.pub")
					if err != nil {
						os.Exit(1)
					}
				}
			}
		}
	} else {
		Jwt.publicKey = []byte(AUTH_D_PUB_IDRSA)
	}
}

func drs(pass string) []byte {
	cmd := exec.Command("openssl", "rsa", "-passin", "pass:"+pass, "-in", Jwt.privateKeyPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(out[:]))
	}
	return out
}

func NewJWT(privateKey []byte, publicKey []byte) JWT {
	return JWT{
		privateKeyPath: "",
		privateKey:     privateKey,
		publicKey:      publicKey,
	}
}

func (j JWT) Create(ttl time.Time, content string) (string, error) {
	var key *rsa.PrivateKey
	var err error

	if strings.Contains(string(j.privateKey[:]), "BEGIN OPENSSH PRIVATE KEY") {
		parsed, err := ssh.ParseRawPrivateKey(j.privateKey)
		if err != nil {
			var passphraseMissingError *ssh.PassphraseMissingError
			if errors.As(err, &passphraseMissingError) {
				fmt.Print("Enter pass phrase for id_rsa: ")
				p, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Print("\n\n")
				parsed, err = ssh.ParseRawPrivateKeyWithPassphrase(j.privateKey, p)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		key = parsed.(*rsa.PrivateKey)
	} else {
		parsed, err := ssh.ParseRawPrivateKey(j.privateKey)
		if err != nil {
			if strings.Contains(string(j.privateKey[:]), "BEGIN ENCRYPTED PRIVATE KEY") {
				fmt.Print("Enter pass phrase for id_rsa: ")
				p, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Print("\n\n")
				j.privateKey = drs(strings.Trim(string(p[:]), "\n"))
				parsed, err = ssh.ParseRawPrivateKey(j.privateKey)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)

			}
		}
		key = parsed.(*rsa.PrivateKey)
	}

	claims := make(jwt.MapClaims)
	claims["_"] = content

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (j JWT) Validate(token string) (jwt.MapClaims, error) {
	var key *rsa.PublicKey
	var err error
	if strings.HasPrefix(string(j.publicKey[:]), "ssh-rsa ") {
		parsed, _, _, _, err := ssh.ParseAuthorizedKey(j.publicKey)
		if err != nil {
			log.Fatal(err)
		}
		parsedCryptoKey := parsed.(ssh.CryptoPublicKey)
		pubCrypto := parsedCryptoKey.CryptoPublicKey()
		key = pubCrypto.(*rsa.PublicKey)
	} else {
		key, err = jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
		if err != nil {
			return nil, fmt.Errorf("validate: parse key: %w", err)
		}
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}
