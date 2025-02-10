package numbertheoretic_methods_in_cryptography

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"slices"
)

func ExtendedEuclideanAlgorithm(a int, b int) int {
	remainders := []int{a, b}
	xs := []int{1, 0}
	ys := []int{0, 1}
	counterIterations := 1
	for {
		currentQuotines := remainders[counterIterations-1] / remainders[counterIterations]
		nextRemainder := remainders[counterIterations-1] % remainders[counterIterations]
		if nextRemainder == 0 {
			return remainders[counterIterations]
		} else {
			remainders = append(remainders, nextRemainder)
			xs = append(xs, xs[counterIterations-1]-currentQuotines*xs[counterIterations])
			ys = append(ys, ys[counterIterations-1]-currentQuotines*ys[counterIterations])
			counterIterations += 1
		}
	}
}

func PrimaryCheck(number int) bool {
	for currentNumber := 2; currentNumber <= int(math.Sqrt(float64(number))); currentNumber++ {
		if number%currentNumber == 0 {
			return false
		}
	}
	return true
}

func GetEulersFunction(number int) (result int) {
	result = 1
	numberStart := number
	primaryNumbers := []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107,
		109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229,
		233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293,
	}
	multipers := make(map[int]int, 1)
	for {
		if slices.Index(primaryNumbers, number) != -1 {
			log.Printf("%d - primary", number)
			multipers[number] += 1
			logMainStrings := fmt.Sprintf("\nf(%d) = ", numberStart)
			var logDecompositionString, logMultipulicationString string
			for currentKey, currentValue := range multipers {
				logDecompositionString = fmt.Sprintf("%s * f(%d^%d)", logDecompositionString, currentKey, currentValue)
				currentResult := int(math.Pow(float64(currentKey), float64(currentValue)) - math.Pow(float64(currentKey), float64(currentValue-1)))
				logMultipulicationString = fmt.Sprintf("%s * %d", logMultipulicationString, currentResult)
				result *= currentResult
			}
			fmt.Printf("%s%s = %s = %d\n", logMainStrings, logDecompositionString[2:], logMultipulicationString[2:], result)
			return
		}
		for _, currentPrimaryNumber := range primaryNumbers {
			if number%currentPrimaryNumber == 0 {
				currentMultiper := number / currentPrimaryNumber
				multipers[currentPrimaryNumber] += 1
				number = currentMultiper
				break
			}
		}
	}
}

func ModuloReduction(numberInt int, power int, module int) (remainder int) {
	numberBig := big.NewInt(int64(numberInt))
	numberInPower := numberBig.Exp(numberBig, big.NewInt(int64(power)), nil)
	_, resultBig := new(big.Int).DivMod(numberInPower, big.NewInt(int64(module)), new(big.Int))
	return int(resultBig.Int64())
}

func GetInverseElement(number int, power int) (result int, err error) {
	log.Printf("%d^(-1)mod(%d)\n", number, power)
	if ExtendedEuclideanAlgorithm(number, power) != 1 {
		err = fmt.Errorf("there is no inverse element")
		return
	}
	log.Printf("%d^(f(%d) - 1)mod(%d)\n", number, power, power)
	result = ModuloReduction(number, GetEulersFunction(power)-1, power)
	return
}
