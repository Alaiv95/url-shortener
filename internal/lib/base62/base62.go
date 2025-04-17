package base62

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func ConvertNum(num int) string {
	if num == 0 {
		return ""
	}

	result := ""

	for num > 0 {
		result += string(base62Chars[num%62])
		num /= 62
	}

	return result
}
