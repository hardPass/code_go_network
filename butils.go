package binary

import (
  "io"
)

// All here are in the order of bigEndian

// Convert uint32 to []byte
func Uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)

	return b
}

// Convert []byte to uint32
func BytesToUint32(b []byte) uint32 {
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

// Write data[]byte with a prefix head which represents the length of data
func Send(w io.Writer, data []byte) (int, error) {
	n := len(data)
	b := Uint32ToBytes(uint32(n))
	_, err := w.Write(b)
	if err != nil {
		return 0, err
	}

	_, err = w.Write(data)
	if err != nil {
		return 0, err
	}

	return (n + 4), nil
}

// Read data which has a head representing size of data
func Receive(r io.Reader) ([]byte, error) {
	b4 := make([]byte, 4)
	_, err := io.ReadFull(r, b4)
	if err != nil {
		return nil, err
	}
	len := BytesToUint32(b4)

	buf := make([]byte, len)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
