package helpers

import "strconv"

func ConvertPointerString(x *string) *string {
	if *x == "" {
		x = nil
	}
	return x
}

func ConvertNilPointerString(x *string) string {
	if x == nil {
		return ""
	}
	return *x
}

func ConvertStringtoInt(value string) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return result
}