package lossless

import (
	"encoding/json"
	"testing"
	"time"
)

type Person struct {
	JSON `json:"-"`

	Name      string `json:"name"`
	Age       int    `json:"age"`
	Address   string
	CreatedAt time.Time
	Ignored   bool `json:"-"`
}

func (p *Person) UnmarshalJSON(data []byte) error {
	return p.JSON.UnmarshalJSON(p, data)
}

func (p Person) MarshalJSON() ([]byte, error) {
	return p.JSON.MarshalJSON(p)
}

type EqualityFunc func(a, b interface{}) bool

var jsondata = []byte(`
{"name": "Jack Wolfington",
 "age": 42,
 "address": "123 Fake St.",
 "CreatedAt": "2013-09-16T10:44:40.295451647-00:00",
 "Ignored": true,
 "Extra": {"foo": "bar"}}`)

func basicEquals(a, b interface{}) bool {
	return a == b
}

func timeEquals(a, b interface{}) bool {
	ta := a.(time.Time)
	tb := b.(time.Time)

	return ta.Equal(tb)
}

func assert(t *testing.T, equal EqualityFunc, actual, expected interface{}) {
	if !equal(actual, expected) {
		t.Fatalf("%#v (%T, actual) != %#v (%T, expected)", actual, actual, expected, expected)
	}
}

func TestDecode(t *testing.T) {
	var p Person
	err := json.Unmarshal(jsondata, &p)
	if err != nil {
		t.Fatal(err)
	}

	assert(t, basicEquals, p.Name, "Jack Wolfington")
	assert(t, basicEquals, p.Age, 42)
	assert(t, basicEquals, p.Address, "123 Fake St.")
	assert(t, timeEquals, p.CreatedAt, time.Date(2013, 9, 16, 10, 44, 40, 295451647, time.UTC))
	assert(t, basicEquals, p.Ignored, false)
}

func TestEncode(t *testing.T) {
	now := time.Now()

	p := Person{
		Name:      "Wolf Jackington",
		Age:       33,
		Address:   "742 Evergreen Terrace",
		CreatedAt: now,
		Ignored:   true,
	}

	p.Set("Pi", 3.14159)

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		t.Fatal(err)
	}

	v, ok := m["name"]
	assert(t, basicEquals, v, p.Name)

	v, ok = m["age"]
	assert(t, basicEquals, int(v.(float64)), p.Age)

	v, ok = m["Address"]
	assert(t, basicEquals, v, p.Address)

	v, ok = m["CreatedAt"]
	assert(t, basicEquals, v, p.CreatedAt.Format(time.RFC3339Nano))

	v, ok = m["Ignored"]
	assert(t, basicEquals, ok, false)

	v, ok = m["Pi"]
	assert(t, basicEquals, v, 3.14159)
}

func testDecodeEncode(t *testing.T) {
	var p Person
	err := json.Unmarshal(jsondata, &p)
	if err != nil {
		t.Fatal(err)
	}

	// Happy birthday, Jack
	p.Age++

	p.Ignored = true
	p.Set("age_printed", "forty-three")

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		t.Fatal(err)
	}

	v, ok := m["age"]
	assert(t, basicEquals, int(v.(float64)), p.Age)

	v, ok = m["Ignored"]
	assert(t, basicEquals, v, false)

	v, ok = m["Extra"]
	assert(t, basicEquals, ok, true)

	m2 := v.(map[string]interface{})
	v, ok = m2["foo"]
	assert(t, basicEquals, v, "bar")

	v, ok = m["age_printed"]
	assert(t, basicEquals, v, "forty-three")
}
