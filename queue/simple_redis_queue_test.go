package queue

import (
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"
)

func TestNewSimpleRedis(t *testing.T) {
	ctx := context.Background()
	redisCnf := &RedisQueueConfig{
		Addr:     "127.0.0.1:6379",
		Password: "",
		Timeout:  0,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 30; i++ {
			err := NewSimpleRedis(redisCnf, "test").Push(ctx, fmt.Sprintf("testNo:%d", i))
			d := math.Ceil(float64(i / 3))
			time.Sleep(time.Duration(d) * time.Second)
			t.Logf("push %d, sleep: %0.f: %+v", i, d, err)
		}
	}()

	wg.Add(1)
	go func() {
		for {
			str, err := NewSimpleRedis(redisCnf, "test").Pop(ctx)
			t.Logf("pop queue: %s, err: %+v", str, err)
		}
	}()

	wg.Wait()
}
