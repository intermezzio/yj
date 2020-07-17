package yaml

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"reflect"

	goyaml "gopkg.in/yaml.v3"
)

type JSON struct {
	EscapeHTML bool
}

func (j *JSON) Marshal(v interface{}) ([]byte, error) {
	keyJSON := &bytes.Buffer{}
	encoder := json.NewEncoder(keyJSON)
	encoder.SetEscapeHTML(j.EscapeHTML)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return keyJSON.Bytes()[:keyJSON.Len()-1], nil
}

func (j *JSON) Unmarshal(src []byte, v interface{}) error {
	if !json.Valid(src) {
		var null interface{}
		err := json.Unmarshal(src, &null)
		if err == nil {
			err = errors.New("invalid JSON")
		}
		return err
	}
	decoder := &Decoder{
		DecodeYAML: func(r io.Reader) (*goyaml.Node, error) {
			var data goyaml.Node
			return &data, goyaml.NewDecoder(r).Decode(&data)
		},
		NaN:    (*float64)(nil),
		PosInf: math.MaxFloat64,
		NegInf: -math.MaxFloat64,
	}
	out, err := decoder.JSON(bytes.NewReader(src))
	if err != nil {
		return err
	}
	val := reflect.ValueOf(v).Elem()
	if !val.CanSet() {
		return errors.New("cannot set value")
	}
	val.Set(reflect.ValueOf(out))
	return nil
}
