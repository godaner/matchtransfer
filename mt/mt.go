package mt

import (
	"errors"
	"reflect"
	"strings"
)

var (
	ErrAbnormalEntry = errors.New("abnormal entry error")
)

const (
	MaxDeep = 10
)

//  MatchTransfer For any incoming object, according to the multi-group matching rules, the attribute value of the unit conversion, return the incoming object
//  data is the original input data and cannot be nil
//  In matchTransferRules, key and value are matching rules and conversion relationships, and cannot be nil, where key is matching rule and value is conversion relationship
//  The maximum nest level of objects hierarchy is 10
func MatchTransfer(matchTransferRules map[string]float64, data interface{}) (interface{}, error) {
	if nil == data || nil == matchTransferRules {
		return nil, ErrAbnormalEntry
	}
	deep := 0
	for k, v := range matchTransferRules {
		onceMatchTransfer(&deep, k, v, data)
	}

	return data, nil
}

// onceMatchTransfer For any input, according to the matching rules and conversion relationship, the value of the attribute is converted
// data is input original data, cannot be nil
// key indicates the matching rule in the format of "key1.key2.key3", and rule indicates the conversion relation
func onceMatchTransfer(deep *int, match string, rule float64, data interface{}) bool {
	if *deep > MaxDeep {
		return false
	}
	*deep++
	split := strings.SplitN(match, ".", 2)

	if data == nil {
		return false
	}

	m, ok := data.(map[string]interface{})
	if ok {
		var firstMatch string
		if len(split) < 2 {
			firstMatch = match
		} else {
			firstMatch = split[0]
			match = split[1]
		}

		value := m[firstMatch]
		if value == nil {
			return false
		}
		kind := reflect.TypeOf(value).Kind()
		if kind == reflect.Float64 {
			m[firstMatch] = m[firstMatch].(float64) * rule
			return true
		}

		return onceMatchTransfer(deep, match, rule, value)
	}

	arrays, ok := data.([]interface{})
	if ok {
		for _, oneArr := range arrays {
			onceMatchTransfer(deep, match, rule, oneArr)
		}
	}

	return false
}
