package numGen

var (
	start = 1
	step  = 1
)

// todo заменить на postgres sequence или придумать что еще
// todo также за раз получать диапазон из 1-2 тысяч чисел и использовать его и когда числа кончатся, запросить новый диапазон
func Generate() int {
	currStart := start
	start += step

	return currStart
}
