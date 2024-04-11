package util

var sortTab = map[string]int{
	"A": 0,
	"B": 1,
	"C": 2,
	"D": 3,
	"E": 4,
	"F": 5,
	"G": 6,
	"H": 7,
}

// list 元素排序
func ListEleSort(data []byte, format string) []byte {
	result := make([]byte, len(data))
	for i, v := range data {
		result[sortTab[string(format[i])]] = v
	}
	return result
}

func GetBitVal(data, index int) int {
	/*
		得到某个字节中某一位（Bit）的值
		:params data: 待取值
		:params index: 待读取位的序号，从右向左0开始，0-7为一个完整字节的8个位
		:returns: 返回读取该位的值，0或1
	*/
	if (data & (1 << index)) > 0 {
		return 1
	} else {
		return 0
	}
}

func GetBitsVal(data, index, len int) int {
	/*
		得到某个字节中某一位（Bit）的值
		:params data: 待取值的字节值, 10进制
		:params startBit: 待读取位的起始序号，从右向左0开始，0-7为一个完整字节的8个位
		:params bitLen: 待读取位长度，1bit，2bit
		:returns: 返回10进制值
	*/

	res := 0

	if len == 1 {
		res = GetBitVal(data, index)
		return res
	}

	dataBlock := index + len
	for i := index; i < dataBlock; i++ {
		if (data & (1 << i)) > 0 {
			res = SetBitVal(res, i-index, 1)
		}
	}
	return res
}

func SetBitVal(data, index, val int) int {
	/*
		data：准备更改的原值。
		index：待更改位的序号，从右向左，以0开始。0-7 表示一个完整字节的8个位。
		val：目标位预更改的值，为 0 或 1。
	*/
	if val == 1 {
		data |= (1 << index)
	} else {
		data &^= (1 << index)
	}
	return data
}
func SetBitsVal(data, index, len, val int) int {
	/*
		data：准备更改的原值。
		index：待更改位的序号，从右向左，以0开始。0-7 表示一个完整字节的8个位。
		len: 从 index 开始更改多少位。
		setVal：目标位预更改的值。

		data = 99
		index = 5
		len = 2
		val = 0

		result = 3
	*/

	res := data

	for i := 0; i < len; i++ {
		bitVal := (val >> i) & 1
		res = SetBitVal(res, i+index, bitVal)
	}
	return res
}
