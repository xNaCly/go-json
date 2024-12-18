package libjson

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkLibJson(b *testing.B) {
	data := strings.Repeat(`{"key1": "value","array": [],"obj": {},"atomArray": [11201,1e112,true,false,null,"str"]},`, 500_000)
	d := []byte("[" + data[:len(data)-1] + "]")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := New(d)
		assert.NoError(b, err)
	}
}

func BenchmarkEncodingJson(b *testing.B) {
	data := strings.Repeat(`{"key1": "value","array": [],"obj": {},"atomArray": [11201,1e112,true,false,null,"str"]},`, 500_000)
	d := []byte("[" + data[:len(data)-1] + "]")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := []struct {
			Key1      string
			Array     []any
			Obj       any
			AtomArray []any
		}{}
		err := json.Unmarshal(d, &v)
		assert.NoError(b, err)
		fmt.Println(v[0])
	}
}
