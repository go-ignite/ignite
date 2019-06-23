package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
)

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
