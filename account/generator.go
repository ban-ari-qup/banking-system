package account

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// структура генератора карт
type CardGenerator struct {
	rng *rand.Rand
}

// создание нового генератора карт
func NewCardGenerator() *CardGenerator {
	source := rand.NewSource(time.Now().UnixNano())
	return &CardGenerator{
		rng: rand.New(source),
	}
}

// генерация номера карты
func (cg *CardGenerator) GenerateCardNumber() string {
	digits := []int{4, 4, 0, 0, 4, 3, 0, 2}

	for i := 0; i < 7; i++ {
		digits = append(digits, cg.rng.Intn(10))
	}

	checkDigit := cg.calculateLuhnCheckDigit(digits)
	digits = append(digits, checkDigit)

	return cg.formatCardNumber(digits)
}

// генерация CVC кода
func (cg *CardGenerator) GenerateCVC() string {
	cvc := cg.rng.Intn(1000)
	for cg.isWeakCVC(cvc) {
		cvc = cg.rng.Intn(1000)
	}
	return fmt.Sprintf("%03d", cvc)
}

// вычисление контрольной цифры по алгоритму Луна
func (cg *CardGenerator) calculateLuhnCheckDigit(digits []int) int {
	sum := 0
	for i, digit := range digits {
		if i%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return (10 - (sum % 10)) % 10
}

// проверка на слабый CVC
func (cg *CardGenerator) isWeakCVC(cvc int) bool {
	if cvc < 10 {
		return true
	}

	str := strconv.Itoa(cvc)

	if len(str) == 3 && str[0] == str[1] && str[1] == str[2] {
		return true
	}

	if len(str) == 3 {
		if str[0]+1 == str[1] && str[1]+1 == str[2] {
			return true
		}
		if str[0]-1 == str[1] && str[1]-1 == str[2] {
			return true
		}
	}

	return false
}

func (cg *CardGenerator) formatCardNumber(digits []int) string { //функция форматирования номера карты
	var result string
	for _, digit := range digits {
		result += strconv.Itoa(digit)
	}
	return result
}
