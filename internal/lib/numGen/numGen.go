package numGen

import (
	"math/rand"
	"time"
)

// Generate временная функция генерации чисел
// todo заменить на postgres sequence или придумать что еще после миграции на БД
// todo также за раз получать диапазон из 1-2 тысяч чисел и использовать его и когда числа кончатся, запросить новый диапазон
func Generate() int {
	rand.Seed(time.Now().UnixNano())

	return rand.Int()
}
