package maya

/*
why yaml?
https://github.com/blog/1647-viewing-yaml-metadata-in-your-documents
*/

import (
	"errors"
	"fmt"
)

type Dict struct {
	m map[interface{}]interface{}
}

func NewDict(m map[interface{}]interface{}) *Dict {
	return &Dict{m}
}

// 대부분의 요소는 string-string이라서 접근하기 쉽도록
func (d *Dict) GetStr(key string) (string, error) {
	raw, ok := d.m[key]
	if !ok {
		return "", errors.New("not found")
	}
	val, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("invalid type: %v", raw)
	}
	return val, nil
}

func (d *Dict) GetInt(key string) (int, error) {
	raw, ok := d.m[key]
	if !ok {
		return 0, errors.New("not found")
	}
	val, ok := raw.(int)
	if !ok {
		return 0, fmt.Errorf("invalid type: %v", raw)
	}
	return val, nil
}

func (d *Dict) GetStrList(key string) ([]string, error) {
	raw, ok := d.m[key]
	if !ok {
		return nil, errors.New("not found")
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type: %#v", raw)
	}

	val := []string{}
	for _, el := range list {
		v, ok := el.(string)
		if ok {
			val = append(val, v)
		} else {
			return nil, fmt.Errorf("contain invalid type: %#v", v)
		}
	}
	return val, nil
}

func (d *Dict) GetIntList(key string) ([]int, error) {
	raw, ok := d.m[key]
	if !ok {
		return nil, errors.New("not found")
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type: %#v", raw)
	}

	val := []int{}
	for _, el := range list {
		v, ok := el.(int)
		if ok {
			val = append(val, v)
		} else {
			return nil, fmt.Errorf("contain invalid type: %#v", v)
		}
	}
	return val, nil
}
