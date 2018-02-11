package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func NewPasswd(length int) string {
	return rand_char(length, StdChars)
}

func rand_char(length int, chars []byte) string {
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

	for port := HOST_From; port <= HOST_To; port++ {
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
