package client

import (
	"bytes"
	"context"
	"io"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr,
	}
}

func (c *Client) Set(ctx context.Context, key, val string) error {
	conn, err := net.Dial("tcp", c.addr)

	if err != nil {
		return err
	}

	var buf *bytes.Buffer
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(val),
	})

	_, err = conn.Write(buf.Bytes())
	io.Copy(conn, buf)

	return err
}