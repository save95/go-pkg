package maker

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransNo(t *testing.T) {
	userId := 18922911
	transNo := TransNo(uint(userId))

	userIdStr := strconv.Itoa(userId)

	assert.Equal(t, len(transNo), 16)
	assert.Equal(t, transNo[0:2], fmt.Sprintf("%d", time.Now().Year()-2000))
	assert.Equal(t, transNo[2:5], fmt.Sprintf("%d", time.Now().YearDay()))
	assert.Equal(t, transNo[7:11], fmt.Sprintf(userIdStr[len(userIdStr)-4:]))
}

func TestTransNoWith(t *testing.T) {
	userId := 18922911

	transNo := TransNoWith(uint(userId), 16)
	assert.Equal(t, len(transNo), 16)

	userIdStr := fmt.Sprintf("%05d", userId)
	transNo2 := TransNoWith(uint(userId), 30)
	assert.Equal(t, len(transNo2), 30)
	assert.Equal(t, transNo2[0:2], fmt.Sprintf("%d", time.Now().Year()-2000))
	assert.Equal(t, transNo2[2:5], fmt.Sprintf("%d", time.Now().YearDay()))
	assert.Equal(t, transNo[len(transNo)-9:len(transNo)-5], fmt.Sprintf(userIdStr[len(userIdStr)-4:]))
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
