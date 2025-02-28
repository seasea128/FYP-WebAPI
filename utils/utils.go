package utils

import "encoding/binary"

func ToByteArray(i int32) []byte {
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr[0:4], uint32(i))
	return arr
}
