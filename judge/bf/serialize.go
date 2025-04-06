package bf

import (
	"encoding/binary"
	"errors"
)

func (bc ByteCode) MarshalBinary() ([]byte, error) { return bc.AppendBinary(nil) }

func (bc ByteCode) AppendBinary(b []byte) ([]byte, error) {
	b = binary.AppendUvarint(b, uint64(len(bc.ops)))

	s := []Op(bc.String())
	// each command takes up 3 bits (8 states)
	// we pack 8 commands into 24 bits = 3 bytes
	// message is padded to a multiple of 3 bytes
	for i := 0; i < len(s); i += 8 {
		var buf [8]Op
		copy(buf[:], s[i:])
		for j := range buf {
			buf[j] = Op(buf[j].index())
		}
		b = append(b,
			byte(buf[2]<<6|buf[1]<<3|buf[0]),              // 22111000
			byte(buf[5]<<7|buf[4]<<4|buf[3]<<1|buf[2]>>2), // 54443332
			byte(buf[7]<<5|buf[6]<<2|buf[5]>>1),           // 77766655
		)
	}

	return b, nil
}

func (bc *ByteCode) UnmarshalBinary(data []byte) error {
	if bc == nil {
		panic("nil receiver")
	}

	size, n := binary.Uvarint(data)
	if n < 0 {
		return errors.New("bad length value")
	}

	data = data[n:]

	if len(data) != int(size+7)/8*3 {
		return errors.New("bad data length")
	}

	ops := make([]Op, 0, len(data)/3*8)

	idx := func(b byte) Op {
		return indexOp[b&7]
	}

	for range len(data) / 3 {
		b := data[:3]
		data = data[3:]

		ops = append(ops,
			idx(b[0]),            // xxxxx000
			idx(b[0]>>3),         // xx111xxx
			idx(b[0]>>6|b[1]<<2), // 22xxxxxx xxxxxxx2
			idx(b[1]>>1),         // xxxx333x
			idx(b[1]>>4),         // x444xxxx
			idx(b[1]>>7|b[2]<<1), // 5xxxxxxx xxxxxx55
			idx(b[2]>>2),         // xxx666xx
			idx(b[2]>>5),         // 777xxxxx
		)
	}

	ops = ops[:size]

	bc2, err := Compile(string(ops), len(ops))
	if err != nil {
		return err
	}

	bc.ops = bc2.ops
	return nil
}
