package account

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type CardGenerator struct {
	rng *rand.Rand
}

func NewCardGenerator() *CardGenerator {
	source := rand.NewSource(time.Now().UnixNano())
	return &CardGenerator{
		rng: rand.New(source),
	}
}
func (cg *CardGenerator) GenerateCVC() string {
	cvc := cg.rng.Intn(1000)
	for cg.isWeakCVC(cvc) {
		cvc = cg.rng.Intn(1000)
	}
	return fmt.Sprintf("%03d", cvc)
}

func (cg *CardGenerator) isWeakCVC(cvc int) bool {
	if cvc < 10 {
		return true
	}
	if cvc < 100 && cvc%11 == 0 {
		return true
	}
	if cvc > 99 {
		str := strconv.Itoa(cvc)
		if str[0] == str[1] && str[1] == str[2] {
			return true
		}
		if str[0]+1 == str[1] && str[1]+1 == str[2] {
			return true
		}
		if str[0]-1 == str[1] && str[1]-1 == str[2] {
			return true
		}
	}
	return false
}

func (cg *CardGenerator) GenerateCardNumber() string {
	digits := []int{4, 4, 0, 0, 4, 3, 0, 2}
	for i := 0; i < 7; i++ {
		digits = append(digits, cg.rng.Intn(10))
	}

	checkDigit := cg.calculateLuhnCheckDigit(digits)
	digits = append(digits, checkDigit)

	return cg.formatCardNumber(digits)
}

// total := 0
// 	for i, num := range arr {
// 		if i%2 == 0 {
// 			num *= 2
// 			if num > 9 {
// 				num -= 9
// 			}
// 		}
// 		total += num
// 	}
// 	last_numb := (10 - (total % 10)) % 10
// 	arr = append(arr, last_numb)

func (cg *CardGenerator) calculateLuhnCheckDigit(digits []int) int {
	sum := 0
	for i, digit := range digits {
		// Позиция справа: (len(digits) - i)
		if i%2 == 0 { // Четные позиции справа
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return (10 - (sum % 10)) % 10
}

func (cg *CardGenerator) formatCardNumber(digits []int) string {
	var result string
	for i, digit := range digits {
		if i > 0 && i%4 == 0 {
			result += " "
		}
		result += strconv.Itoa(digit)
	}
	return result
}
