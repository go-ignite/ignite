package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-ignite/ignite/config"
)

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func RandString(length int) string {
	return randChar(length, StdChars)
}

func randChar(length int, chars []byte) string {
	new_pword := make([]byte, length)
	random_data := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, random_data); err != nil {
			panic(err)
		}
		for _, c := range random_data {
			if c >= maxrb {
				continue
			}
			new_pword[i] = chars[c%clen]
			i++
			if i == length {
				return string(new_pword)
			}
		}
	}
}

func ServiceURL(serviceType, host string, port int, method, password string) string {
	var protocol, base64Link string
	switch serviceType {
	case "SS", "":
		protocol = "ss"
		//method:password@server:port
		base64Link = ssbase64Encode(fmt.Sprintf("%s:%s@%s:%d", method, password, host, port))
	case "SSR":
		protocol = "ssr"
		//server:port:protocol:method:obfs:password_base64/?suffix_base64
		base64Pwd := ssbase64Encode(password)
		suffix := fmt.Sprintf("protoparam=%s", ssbase64Encode("32"))
		base64Link = ssbase64Encode(fmt.Sprintf("%s:%d:%s:%s:%s:%s/?%s", host, port, "auth_aes128_md5", method, "tls1.2_ticket_auth_compatible", base64Pwd, suffix))
	default:
		return ""
	}
	return fmt.Sprintf("%s://%s", protocol, base64Link)
}

func ssbase64Encode(s string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(s))
	return strings.TrimRight(encoded, "=")
}

func GetAvailablePort(usedPorts *[]int) (int, error) {
	portMap := map[int]int{}

	for _, p := range *usedPorts {
		portMap[p] = p
	}

	from, to := config.C.Host.From, config.C.Host.From
	for port := from; port <= to; port++ {
		if _, exists := portMap[port]; exists {
			continue
		}
		conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return port, nil
		}
		conn.Close()
	}
	return 0, errors.New("no port available")
}

func CreateToken(secret string, id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return "Bearer " + tokenStr, nil
}

func VerifyToken(tokenString string, isAdmin *bool) bool {
	tokenString = tokenString[7:]
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected siging method")
		}
		return []byte(config.C.App.Secret), nil
	})
	if err != nil {
		return false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	if !token.Valid {
		return false
	}
	if isAdmin != nil {
		id, ok := claims["id"].(float64)
		if !ok {
			return false
		}
		if (*isAdmin && id != -1) || (!*isAdmin && id <= 0) {
			return false
		}
	}
	return claims.VerifyExpiresAt(time.Now().Unix(), true)
}
