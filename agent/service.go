package agent

import (
	pb "github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/models"
)

func GetServiceConfigs() []*models.ServiceConfig {
	return []*models.ServiceConfig{
		{
			Type:      "ss-libev",
			TypeProto: pb.ServiceType_SS_LIBEV,
			Methods:   []string{"aes-256-cfb", "aes-128-gcm", "aes-192-gcm", "aes-256-gcm", "chacha20-ietf-poly1305"},
		},
		{
			Type:      "ssr",
			TypeProto: pb.ServiceType_SSR,
			Methods:   []string{"aes-256-cfb", "aes-256-ctr", "chacha20", "chacha20-ietf"},
		},
	}
}
