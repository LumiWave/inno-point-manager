package util

import "strconv"

// string -> int64
func ParseInt(data string) int64 {
	value, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

// string 숫자로 변환해서 크기 비교
func CompareString(data1, data2 string) int64 {
	value1 := ParseInt(data1)
	value2 := ParseInt(data2)

	if value1 < value2 {
		return -1 // data1이 작으면
	} else if value1 > value2 {
		return 1 // data1이 더 크면
	}
	return 0 //같으면 0 리턴
}

// string 숫자로 변환해서 더하기
func SumString(data1, data2 string) string {
	return strconv.FormatInt(ParseInt(data1)+ParseInt(data2), 10)
}

// string 숫자로 변환해서 빼기
func SubString(data1, data2 string) string {
	return strconv.FormatInt(ParseInt(data1)-ParseInt(data2), 10)
}
