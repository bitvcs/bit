package snow

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MyStruct struct {
	ID ID `json:"id"`
}

func TestSnow(t *testing.T) {
	node, err := NewNode(1)
	assert.NoError(t, err)
	assert.NotNil(t, node)

	id := node.Generate()
	fmt.Println(id)
	fmt.Println(id.Base36())
	id += 1
	fmt.Println(id.Base36())
	time.Sleep(2 * time.Millisecond)
	id = node.Generate()
	fmt.Println(id.Base36())

	s := MyStruct{ID: node.Generate()}
	jsonByte, err := json.Marshal(&s)
	assert.NoError(t, err)
	fmt.Println(string(jsonByte))

	result := &MyStruct{}
	err = json.Unmarshal(jsonByte, result)
	assert.NoError(t, err)
	assert.Equal(t, s.ID, result.ID)
}

func TestTimestamp(t *testing.T) {
	before, err := time.Parse("2006", "2003")
	assert.NoError(t, err)
	fmt.Println(before.UnixMilli())
	now := time.Since(before).Milliseconds()
	fmt.Println(now)
	second := time.Since(before).Milliseconds() / 1000
	fmt.Println(second)
}
