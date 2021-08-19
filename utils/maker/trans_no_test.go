package maker

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenTransNo(t *testing.T) {
	userId := 18922911
	transNo := TransNo(uint(userId))

	userIdStr := strconv.Itoa(userId)

	assert.Equal(t, len(transNo), 15)
	assert.Equal(t, transNo[0:2], fmt.Sprintf("%d", time.Now().Year()-2000))
	assert.Equal(t, transNo[2:5], fmt.Sprintf("%d", time.Now().YearDay()))
	assert.Equal(t, transNo[7:10], fmt.Sprintf(userIdStr[len(userIdStr)-3:]))
}

func TestGenTransNo_Sync(t *testing.T) {
	var wg sync.WaitGroup
	var set sync.Map
	max := 1000
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(uid int) {
			defer wg.Done()

			for j := 0; j < max; j++ {
				transNo := TransNo(uint(uid))
				v, ok := set.Load(transNo)
				if ok {
					set.Store(transNo, v.(int)+1)
				} else {
					set.Store(transNo, 1)
				}
			}
		}(i)
	}

	wg.Wait()

	set.Range(func(key, value interface{}) bool {
		v, ok := value.(int)
		if !ok {
			return false
		}

		if v > 1 {
			fmt.Printf("%s \t %d\n", key, v)
		}
		return true
	})
}
