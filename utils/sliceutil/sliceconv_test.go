package sliceutil

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToAny(t *testing.T) {
	source := []int{1, 2, 3, 4, 5}
	res, ok := ToAny(source)

	assert.True(t, ok)
	assert.Equal(t, res, []interface{}{1, 2, 3, 4, 5})
}

func TestToFloat64(t *testing.T) {
	sourceStr := `{"float64":[1,2,3,4,5,6]}`
	var source map[string]interface{}
	err := json.Unmarshal([]byte(sourceStr), &source)
	if nil != err {
		t.Error(err)
		return
	}
	if v, ok := source["float64"].([]interface{}); ok {
		fmt.Printf("%#v\n", ToFloat64(v))
	}
}

func TestToInt(t *testing.T) {
	source := []interface{}{1, 2, 3, 4, 5, "1", "b", "c"}
	assert.Equal(t, ToInt(source), []int{1, 2, 3, 4, 5})
}

func TestToPossibleInt(t *testing.T) {
	source := []interface{}{1, 2, 3, 4, 5, "6", "b", "c", "", "7.0", "8.53"}
	assert.Equal(t, ToPossibleInt(source), []int{1, 2, 3, 4, 5, 6, 7, 8})
}

func TestToString(t *testing.T) {
	source := []interface{}{1, 2, 3, 4, 5, "1", "b", "c"}
	assert.Equal(t, ToString(source), []string{"1", "b", "c"})
}
