package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type TestCase struct {
	name     string
	address   map[string]string
	codes      []string
	count     int
	opts map[string]interface{}
}

var testCase = TestCase{
	name: "test",
	address: map[string]string{
		"street": "Liconl",
		"city":   "London",
	},
	codes: []string{
		"123",
		"456",
	},
	count: 2,
	opts: map[string]interface{}{
		"address2": map[string]interface{}{
			"street": "742 Evergreen Terrace",
			"city":   "LA",
		},
		"codes": []string{
			"985",
			"144",
		},
		"notes": []string{
			"note1",
			"note2",
		},
	},
}

func generate(data map[string]interface{}) map[string]interface{}  {

		for key, v := range data {
			if assert, ok := v.(map[string]interface{}); ok {
				generate(assert)
			} else {
				if av, ok := v.(string); ok {
					data[key] = imitateString(av)
				}
				if av, ok := v.(float64); ok {
					data[key] = imitateNumber(av)
				}
				if assert, ok := v.([]interface{}); ok {
					rv := reflect.ValueOf(assert)
					ret := make([]interface{}, rv.Len())
					for i, v := range assert {
						if av, ok := v.(string); ok {
							ret[i] = imitateString(av)
						}
						if av, ok := v.(float64); ok {
							ret[i] = imitateNumber(av)
						}
						if assert, ok := v.(map[string]interface{}); ok {
							ret[i] = generate(assert)
						}
					}
					data[key] = ret
				}
			}
	}

	return data
}

func imitateString(s string) string {
	digitStringReg := regexp.MustCompile(`^[0-9]+$`)
	specials := []string{"]", "^", "\\\\", "[", ".", "(", ")", "-", "*", "+", "?", "|", "{", " "}
	specialsReg := regexp.MustCompile("[" + strings.Join(specials, "") + "]+")

	vovels := "aeiou"
	consonants := "bcdfghjklmnpqrstvwxyz"
	vovelsReg := regexp.MustCompile("[" + vovels + "]+")

	b := make([]string, len(s))
	for i := 0; i < len(s); i++ {
		fromIndex := s[i:i+1]
		if specialsReg.MatchString(fromIndex) {
			b[i] = fromIndex
		} else {
			if digitStringReg.MatchString(fromIndex) {
				b[i] = strconv.Itoa(rand.Intn(10 - 1) + 1)
			}else{
				if vovelsReg.MatchString(fromIndex) {
					b[i] = string(vovels[rand.Intn(len(vovels)-1)])
				} else {
					b[i] = string(consonants[rand.Intn(len(consonants)-1)])
				}
				if strings.ToUpper(fromIndex) == fromIndex {
					b[i] = strings.ToUpper(b[i])
				}
			}
		}
	}
	generated := strings.Join(b, "")
	return generated
}


func imitateNumber(n float64) float64 {
	simple := uint64(n)

	count := 0
	for simple != 0 {
		simple /= 10
		count += 1
	}

	s := ""

	for i := 0; i < count; i++ {
		s += strconv.Itoa(rand.Intn(10 - 1) + 1)
	}

	floatNum, _ := strconv.ParseFloat(s, 64)

	return floatNum
}


func CreateRequest() *http.Request {
	var pes = map[string]interface{}{
		"address2": map[string]interface{}{
			"street": "742 Evergreen Terrace",
			"city":   "LA",
			"numver": 2,
		},
		"codes": []interface{}{
			"985",
			"144",
		},
		"notes": []string{
			"note1",
			"note2",
			"ka[usta",
		},
		"string": "string",
		"inte": 665,
	}
	b, _ := json.Marshal(pes)
	reader := bytes.NewReader(b)

	req := httptest.NewRequest(http.MethodGet, "/imitate", reader)
	req.Header.Set("Content-Type", "application/json")

	var body map[string]interface{}

	json.NewDecoder(req.Body).Decode(&body)


	generate(body)

	return req
}

func TestUpperCaseHandler(t *testing.T) {
	CreateRequest()
}
