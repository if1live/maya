package maya

import (
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func Test_Dict_GetStr(t *testing.T) {
	cases := []struct {
		data     string
		expected string
		ok       bool
	}{
		{"key: foo", "foo", true},
		{"key: 123", "", false},
		{"key: foo bar", "foo bar", true},
		{"key: [foo, bar]", "", false},
	}
	for _, c := range cases {
		m := yaml.MapSlice{}
		err := yaml.Unmarshal([]byte(c.data), &m)
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		dict := NewDict(m)
		actual, err := dict.GetStr("key")
		if c.ok && err != nil {
			t.Errorf("should success, but fail - data: %s, err %q", c.data, err)
		}
		if !c.ok && err == nil {
			t.Errorf("should fail, but succes - data: %s, err %q", c.data, err)
		}

		if c.ok && err == nil {
			if c.expected != actual {
				t.Errorf("expected %q, got %q", c.expected, actual)
			}
		}
	}
}

func Test_Dict_GetInt(t *testing.T) {
	cases := []struct {
		data     string
		expected int
		ok       bool
	}{
		{"key: 123", 123, true},
		{"key: foo", 0, false},
	}
	for _, c := range cases {
		m := yaml.MapSlice{}
		err := yaml.Unmarshal([]byte(c.data), &m)
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		dict := NewDict(m)
		actual, err := dict.GetInt("key")
		if c.ok && err != nil {
			t.Errorf("should success, but fail - data: %s, err %q", c.data, err)
		}
		if !c.ok && err == nil {
			t.Errorf("should fail, but succes - data: %s, err %q", c.data, err)
		}

		if c.ok && err == nil {
			if c.expected != actual {
				t.Errorf("expected %q, got %q", c.expected, actual)
			}
		}
	}
}

func Test_Dict_GetStrList(t *testing.T) {
	cases := []struct {
		data     string
		expected []string
		ok       bool
	}{
		{"key: [foo, bar]", []string{"foo", "bar"}, true},
		{"key: [1, 2]", nil, false},
		{"key: foo, bar", nil, false},
	}
	for _, c := range cases {
		m := yaml.MapSlice{}

		err := yaml.Unmarshal([]byte(c.data), &m)
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		dict := NewDict(m)
		actual, err := dict.GetStrList("key")
		if c.ok && err != nil {
			t.Errorf("should success, but fail - data: %s, err %q", c.data, err)
		}
		if !c.ok && err == nil {
			t.Errorf("should fail, but succes - data: %s, err %q", c.data, err)
		}

		if c.ok && err == nil {
			if !reflect.DeepEqual(c.expected, actual) {
				t.Errorf("data: %s, expected %q, got %q", c.data, c.expected, actual)
			}
		}
	}
}

func Test_Dict_GetIntList(t *testing.T) {
	cases := []struct {
		data     string
		expected []int
		ok       bool
	}{
		{"key: [foo, bar]", nil, false},
		{"key: [1, 2]", []int{1, 2}, true},
		{"key: foo, bar", nil, false},
	}
	for _, c := range cases {
		m := yaml.MapSlice{}
		err := yaml.Unmarshal([]byte(c.data), &m)
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		dict := NewDict(m)
		actual, err := dict.GetIntList("key")
		if c.ok && err != nil {
			t.Errorf("should success, but fail - data: %s, err %q", c.data, err)
		}
		if !c.ok && err == nil {
			t.Errorf("should fail, but succes - data: %s, err %q", c.data, err)
		}

		if c.ok && err == nil {
			if !reflect.DeepEqual(c.expected, actual) {
				t.Errorf("data: %s, expected %q, got %q", c.data, c.expected, actual)
			}
		}
	}
}

func Test_Dict_GetStrKeys(t *testing.T) {
	data := `
foo: bar
spam: dummy
`
	expected := []string{"foo", "spam"}
	m := yaml.MapSlice{}
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	dict := NewDict(m)
	actual := dict.GetStrKeys()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("data: %s, expected %q, got %q", data, expected, actual)
	}
}

func Test_Dict_GetValueType(t *testing.T) {
	cases := []struct {
		data     string
		expected valueType
	}{
		{"key: foo", valueTypeStr},
		{"key: 123", valueTypeInt},
		{"key: [foo, bar]", valueTypeStrList},
		{"key: [1, 2]", valueTypeIntList},
	}
	for _, c := range cases {
		m := yaml.MapSlice{}
		err := yaml.Unmarshal([]byte(c.data), &m)
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		dict := NewDict(m)
		actual := dict.GetValueType("key")
		if c.expected != actual {
			t.Errorf("data: %s, expected %v, got %v", c.data, c.expected, actual)
		}
	}
}
