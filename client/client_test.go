package client

import (
	"context"
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("localhost:6969")

	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Println("SET =>", fmt.Sprintf("bar_%d", i))
		if err := client.Set(context.Background(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			t.Error(err)
		}

		val, err := client.Get(context.Background(), fmt.Sprintf("foo_%d", i))

		if err != nil {
			t.Error(err)
		}

		if val != fmt.Sprintf("bar_%d", i) {
			t.Errorf("GET: expected bar_%d, got %v", i, val)
		}
	}
}
