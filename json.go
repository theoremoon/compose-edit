// copy from https://github.com/kayac/ecspresso/blob/d60a9c7d3a30218a2a4ecf8b4c38cc132f2b4965/json.go
// ecspresso has large dependencies

/*
MIT License

Copyright (c) 2017 KAYAC Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package composeedit

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/itchyny/gojq"
)

func MarshalJSONForAPI(v interface{}, queries ...string) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	walkMap(m, jsonKeyForAPI)
	if len(queries) > 0 {
		for _, q := range queries {
			if m, err = jqFilter(m, q); err != nil {
				return nil, err
			}
		}
	}
	bs, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, err
	}
	bs = append(bs, '\n')
	return bs, nil
}

func jsonKeyForAPI(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func walkMap(m map[string]interface{}, fn func(string) string) {
	for key, value := range m {
		delete(m, key)
		newKey := key
		if fn != nil {
			newKey = fn(key)
		}
		if value != nil {
			m[newKey] = value
		}
		switch value := value.(type) {
		case map[string]interface{}:
			switch strings.ToLower(key) {
			case "dockerlabels", "options":
				walkMap(value, nil) // do not rewrite keys for map[string]string
			default:
				walkMap(value, fn)
			}
		case []interface{}:
			if len(value) > 0 {
				walkArray(value, fn)
			} else {
				delete(m, newKey)
			}
		default:
		}
	}
}

func walkArray(a []interface{}, fn func(string) string) {
	for _, value := range a {
		switch value := value.(type) {
		case map[string]interface{}:
			walkMap(value, fn)
		case []interface{}:
			walkArray(value, fn)
		default:
		}
	}
}

func jqFilter(m map[string]interface{}, q string) (map[string]interface{}, error) {
	query, err := gojq.Parse(q)
	if err != nil {
		return nil, err
	}
	iter := query.Run(m)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}
		if m, ok = v.(map[string]interface{}); !ok {
			return nil, fmt.Errorf("query result is not map[string]interface{}: %v", v)
		}
	}
	return m, nil
}
