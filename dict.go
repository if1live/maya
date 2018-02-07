package maya

/*
why yaml?
https://github.com/blog/1647-viewing-yaml-metadata-in-your-documents
*/

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type valueType int

const (
	valueTypeUnknown valueType = 0
	valueTypeStr               = 1
	valueTypeInt               = 2
	valueTypeStrList           = 3
	valueTypeIntList           = 4
)

type Dict struct {
	m yaml.MapSlice
}

func NewDict(m yaml.MapSlice) *Dict {
	return &Dict{m}
}

func (d *Dict) GetValueType(key string) valueType {
	_, err := d.GetStr(key)
	if err == nil {
		return valueTypeStr
	}

	_, err = d.GetInt(key)
	if err == nil {
		return valueTypeInt
	}

	_, err = d.GetStrList(key)
	if err == nil {
		return valueTypeStrList
	}

	_, err = d.GetIntList(key)
	if err == nil {
		return valueTypeIntList
	}

	return valueTypeUnknown
}

// 대부분의 요소는 string-string이라서 접근하기 쉽도록
func (d *Dict) getRootValue(key string) (interface{}, error) {
	for _, item := range d.m {
		if k, ok := item.Key.(string); ok {
			if k == key {
				return item.Value, nil
			}
		}
	}
	return nil, fmt.Errorf("not found: %s", key)
}

func (d *Dict) GetStr(key string) (string, error) {
	raw, err := d.getRootValue(key)
	if err != nil {
		return "", err
	}

	val, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("invalid type: %v", raw)
	}
	return val, nil
}

func (d *Dict) GetInt(key string) (int, error) {
	raw, err := d.getRootValue(key)
	if err != nil {
		return 0, err
	}
	val, ok := raw.(int)
	if !ok {
		return 0, fmt.Errorf("invalid type: %v", raw)
	}
	return val, nil
}

func (d *Dict) GetStrList(key string) ([]string, error) {
	raw, err := d.getRootValue(key)
	if err != nil {
		return nil, err
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
	raw, err := d.getRootValue(key)
	if err != nil {
		return nil, err
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

func (d *Dict) GetStrKeys() []string {
	keys := []string{}
	for _, item := range d.m {
		k := item.Key
		if key, ok := k.(string); ok {
			keys = append(keys, key)
		}
	}
	return keys
}
