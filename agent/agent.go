package agent

import (
	pb "github.com/go-ignite/ignite-agent/protos"

	"google.golang.org/grpc"
)

type Client struct {
	pb.AgentServiceClient
	conn    *grpc.ClientConn
	address string
}

func NewClient(address string) *Client {
	return &Client{address: address}
}

func (client *Client) Dial() error {
	var err error
	client.conn, err = grpc.Dial(client.address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client.AgentServiceClient = pb.NewAgentServiceClient(client.conn)
	return nil
}

func (client *Client) Close() error {
	return client.conn.Close()
}

func Dial(address string) (*Client, error) {
	client := NewClient(address)
	if err := client.Dial(); err != nil {
		return nil, err
	}
	return client, nil
}
