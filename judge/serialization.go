package judge

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)

type serializedGenerator struct {
	List *listGenerator
	BF   *bfGenerator
	Lua  *luaGenerator
}

func MarshalGenerator(g InputGenerator) ([]byte, error) {
	return AppendGenerator(g, nil)
}

func AppendGenerator(g InputGenerator, b []byte) ([]byte, error) {
	var val serializedGenerator
	switch gen := g.(type) {
	case bfGenerator:
		val = serializedGenerator{
			BF: &gen,
		}
	case listGenerator:
		val = serializedGenerator{
			List: &gen,
		}
	case luaGenerator:
		val = serializedGenerator{
			Lua: &gen,
		}
	case *smartGenerator:
		return AppendGenerator(gen.InputGenerator, b)
	default:
		return nil, fmt.Errorf("generator type %T is not known", g)
	}

	buf := bytes.NewBuffer(b)
	err := gob.NewEncoder(buf).Encode(val)
	return buf.Bytes(), err
}

func UnmarshalGenerator(b []byte) (InputGenerator, error) {
	var val serializedGenerator
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&val); err != nil {
		return nil, err
	}

	if val.Lua != nil {
		return &smartGenerator{*val.Lua}, nil
	}
	if val.List != nil {
		return *val.List, nil
	}
	if val.BF != nil {
		return &smartGenerator{*val.BF}, nil
	}
	return nil, errors.New("value did not contain any generator")
}
