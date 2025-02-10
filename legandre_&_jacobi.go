package numbertheoretic_methods_in_cryptography

import (
	"log"
	"math"
	"math/big"
)

func GetLegandreSymbol(first int, second int) (result int) {
	firstBig, power := big.NewInt(int64(first)), big.NewInt(int64((second-1)/2))
	expression := firstBig.Exp(firstBig, power, nil)
	_, resultBig := new(big.Int).DivMod(expression, big.NewInt(int64(second)), new(big.Int))
	if int(resultBig.Int64()) == second-1 {
		return -1
	}
	return int(resultBig.Int64())
}

func factorization(number int) map[int]int {
	primeFactors := make(map[int]int)
	for currentDivider := 2; currentDivider < int(math.Sqrt(float64(number))); currentDivider++ {
		for number%currentDivider == 0 {
			primeFactors[currentDivider] += 1
			number /= currentDivider
		}
	}
	if number != 1 {
		primeFactors[number] += 1
	}
	return primeFactors
}

func GetLegendreOrJacobiSymbol(first int, second int) (result int) {
	if PrimaryCheck(second) {
		log.Printf("%d is primary -> starting calculation the L(%d, %d)", second, first, second)
		result = GetLegandreSymbol(first, second)
		log.Printf(" L(%d, %d) = %d\n", first, second, result)
	} else {
		log.Printf("%d isn't primary -> starting calculation the J(%d, %d)", second, first, second)
		factorials := factorization(second)
		result = 1
		log.Printf("Factorization of %d: %v\n", second, factorials)
		for currentNumber, currentPower := range factorials {
			currentResult := GetLegandreSymbol(first, currentNumber)
			log.Printf(" L(%d, %d) = %d", first, currentNumber, currentResult)
			result *= int(math.Pow(float64(currentResult), float64(currentPower)))
		}
	}
	return
}
