package numGen

var (
	start = 1
	step  = 1
)

func Generate() int {
	currStart := start
	start += step

	return currStart
}
