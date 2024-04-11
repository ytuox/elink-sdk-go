package util

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/spf13/cast"
	"github.com/ytuox/elink-sdk-go/common"
)

func GetNRegister(data int) uint16 {
	return uint16((data + 15) / 16)
}

func GetNByte(data int) uint16 {
	return uint16((data + 7) / 8)
}

func GetDataTypeByteLen(_type string) uint16 {
	var byteLen uint16

	switch _type {
	case common.ValueTypeBool, common.ValueTypeInt8, common.ValueTypeUint8, common.ValueTypeInt16, common.ValueTypeUint16:
		byteLen = 2
	case common.ValueTypeInt32, common.ValueTypeUint32, common.ValueTypeFloat32:
		byteLen = 4
	case common.ValueTypeInt64, common.ValueTypeUint64, common.ValueTypeFloat64:
		byteLen = 8
	}

	return byteLen
}

func BytesToX(inData []byte, outType string) interface{} {

	var res interface{}

	switch outType {
	case common.ValueTypeBool:
		res = binary.BigEndian.Uint16(inData)
	case common.ValueTypeUint8:
		res = binary.BigEndian.Uint16(inData)
	case common.ValueTypeInt8:
		res = int8(binary.BigEndian.Uint16(inData))
	case common.ValueTypeUint16:
		res = binary.BigEndian.Uint16(inData)
	case common.ValueTypeInt16:
		res = int16(binary.BigEndian.Uint16(inData))
	case common.ValueTypeUint32:
		res = binary.BigEndian.Uint32(inData)
	case common.ValueTypeUint64:
		res = binary.BigEndian.Uint64(inData)
	case common.ValueTypeInt64:
		res = int64(binary.BigEndian.Uint64(inData))
	case common.ValueTypeFloat32:
		// CDAB
		raw := binary.BigEndian.Uint32(inData)
		res = math.Float32frombits(raw)
	case common.ValueTypeFloat64:
		raw := binary.BigEndian.Uint64(inData)
		res = math.Float64frombits(raw)
	}

	return res
}

func XToBytes(inData interface{}, outType string) ([]byte, error) {

	byteLen := GetDataTypeByteLen(outType)
	bytes := make([]byte, byteLen)

	switch outType {
	case common.ValueTypeBool:
		binary.BigEndian.PutUint16(bytes, cast.ToUint16(inData))
	case common.ValueTypeInt8:
		binary.BigEndian.PutUint16(bytes, cast.ToUint16(inData))
	case common.ValueTypeUint8:
		binary.BigEndian.PutUint16(bytes, cast.ToUint16(inData))
	case common.ValueTypeInt:
		binary.BigEndian.PutUint16(bytes, cast.ToUint16(inData))
	case common.ValueTypeInt16:
		binary.BigEndian.PutUint16(bytes, cast.ToUint16(inData))
	case common.ValueTypeUint16:
		binary.BigEndian.PutUint16(bytes, cast.ToUint16(inData))
	case common.ValueTypeFloat32:
		bits := math.Float32bits(cast.ToFloat32(inData))
		binary.BigEndian.PutUint32(bytes, bits)
	case common.ValueTypeFloat64:
		bits := math.Float64bits(cast.ToFloat64(inData))
		binary.BigEndian.PutUint64(bytes, bits)
	default:
		return nil, fmt.Errorf("unsupported output type: %s", outType)
	}

	return bytes, nil
}
