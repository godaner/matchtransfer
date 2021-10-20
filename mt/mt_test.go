package mt

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestMatchTransfer(t *testing.T) {
	data := `{
    "key":{
        "key2":[
            {
                "key3":[{"key6":15.4},{"key6":8},{"key6":36}],
				  "key4":6.3
            },
            {
                "key3":[{"key6":12},{"key6":3},{"key6":25}],
				  "key4":8
            }
        ],
        "foo":123.2,
        "bar":{
            "foo":456
        },
        "no":123
    }
}`

	rules := map[string]float64{"key.foo": 2, "key.bar.foo": 3, "key.key2.key3.key6": 1.5, "key.key2.key4": 2.5}
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatal(err)
	}
	MatchTransfer(rules, m)
	fmt.Println(m)
}

// pressure test
func BenchmarkMatchTransfer(b *testing.B) {
	data := `{
    "key":{
        "key2":[
            {
                "key3":[{"key6":4},{"key6":5}],
				  "key4":6
            }
        ],
        "foo":123,
        "bar":{
            "foo":456
        },
        "no":123
    }
}`

	rules := map[string]float64{"key.foo": 2, "key.key2.key4": 0.5, "key.bar.foo": 3, "key.key2.key3.key6": 1.5}

	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatal(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// MatchTransfer(rules,m)
		// MatchTransfer(rules, DeepCopy(m))
		fmt.Println(MatchTransfer(rules, DeepCopy(m)))
	}
}

func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}
		return newSlice
	}

	return value
}
