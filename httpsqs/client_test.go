package httpsqs

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_client = NewClient(&Config{
		Addr:     "127.0.0.1:1218",
		Timeout:  0,
		Password: "",
	})
	name = "queue_test"
	ctx  = context.Background()
)

func TestClient_Put(t *testing.T) {
	bs, _ := json.Marshal(struct {
		Name  string
		Title string
		Age   int
		Money int
	}{
		"张三", "无产阶级", 30, 10000,
	})
	pos, err := _client.Put(ctx, name, string(bs))
	assert.Nil(t, err)
	assert.Less(t, int64(0), pos)
}

func TestClient_Get(t *testing.T) {
	body, pos, err := _client.Get(ctx, name)
	assert.Nil(t, err)
	assert.Less(t, int64(0), pos)
	assert.Less(t, 0, len(body))
	t.Logf("queue pos: %d, data: %s", pos, body)
}

func TestClient_Status(t *testing.T) {
	status, err := _client.Status(ctx, name)
	assert.Nil(t, err)
	t.Logf("queue status: %#v", status)
	t.Logf("queue status: %s", status)
}

func TestClient_View(t *testing.T) {
	body, err := _client.View(ctx, name, 1)
	assert.Nil(t, err)
	t.Logf("queue view data: %s", body)
}

func TestClient_SetMaxQueue(t *testing.T) {
	max := 100000000
	err := _client.SetMaxQueue(ctx, name, max)
	assert.Nil(t, err)

	status, err := _client.Status(ctx, name)
	assert.Nil(t, err)
	assert.Equal(t, int64(max), status.MaxQueue)
}

func TestClient_Reset(t *testing.T) {
	err := _client.Reset(ctx, name)
	assert.Nil(t, err)
}
