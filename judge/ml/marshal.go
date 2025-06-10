package ml

import (
	"bytes"
	"encoding/gob"
)

func zero[T any]() (res T) { return }

func init() {
	gob.Register(zero[Title]())
	gob.Register(zero[List]())
	gob.Register(zero[Paragraph]())
	gob.Register(zero[CodeBlock]())
	gob.Register(zero[Example]())
	gob.Register(zero[Quote]())
	gob.Register(zero[Image]())
}

func (d *Document) MarshalBinary() ([]byte, error) {
	return d.AppendBinary(nil)
}

func (d *Document) AppendBinary(dst []byte) ([]byte, error) {
	buf := bytes.NewBuffer(dst)
	if err := gob.NewEncoder(buf).Encode(d); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *Document) UnmarshalBinary(src []byte) error {
	return gob.NewDecoder(bytes.NewReader(src)).Decode(d)
}
