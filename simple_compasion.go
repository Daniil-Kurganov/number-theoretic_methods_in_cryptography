package numbertheoretic_methods_in_cryptography

import (
	"fmt"
	"log"
	"math"
	"math/big"
)

func ExpansionOfSuitableFraction(numerator uint, divider uint) (result []uint) {
	for {
		if numerator%divider == 0 {
			result = append(result, numerator/divider)
			return
		}
		result = append(result, numerator/divider)
		numerator, divider = divider, numerator-result[len(result)-1]*divider
	}
}

func CalculationOfSuitableFactorsByTable(continuedFraction []uint) (result [][]uint) {
	numerators, dividers := []uint{1, continuedFraction[0]}, []uint{0, 1}
	for currentIndex := 1; currentIndex < len(continuedFraction); currentIndex++ {
		numerators = append(numerators, continuedFraction[currentIndex]*numerators[currentIndex]+numerators[currentIndex-1])
		dividers = append(dividers, continuedFraction[currentIndex]*dividers[currentIndex]+dividers[currentIndex-1])
		result = append(result, []uint{numerators[len(numerators)-1], dividers[len(dividers)-1]})
	}
	return
}

func SimpleCompasionSolutionByContinuedFraction(a uint, b uint, module uint) (result []uint, err error) {
	log.Printf("\n%d * x ≡ %d mod(%d)", a, b, module)
	DCM := ExtendedEuclideanAlgorithm(int(a), int(module))
	log.Printf("DCM(%d, %d) = %d", a, module, DCM)
	if b%uint(DCM) != 1 {
		a /= uint(DCM)
		b /= uint(DCM)
		module /= uint(DCM)
		log.Printf("DCM > 1 -> compasion has %d solutions\n\n%d * x ≡ %d mod(%d)", DCM, a, b, module)
		if ExtendedEuclideanAlgorithm(int(a), int(module)) != 1 || b%uint(ExtendedEuclideanAlgorithm(int(a), int(module))) != 0 {
			err = fmt.Errorf("current compasion hasn't got solution")
			return
		}
	}
	continuedFraction := ExpansionOfSuitableFraction(a, module)
	log.Printf("Continued fraction: %v", continuedFraction)
	suitableFactors := CalculationOfSuitableFactorsByTable(continuedFraction)
	log.Printf("Suitable factors: %v", suitableFactors)
	if a != suitableFactors[len(suitableFactors)-1][0] || module != suitableFactors[len(suitableFactors)-1][1] {
		err = fmt.Errorf("error on table calculating method")
		return
	}
	k, P, Q := len(continuedFraction)-1, suitableFactors[len(suitableFactors)-2][0], suitableFactors[len(suitableFactors)-2][1]
	log.Printf("k = %d, Pₖ₋₁ = %d, Qₖ₋₁ = %d", k, P, Q)
	zeroX := ModuloReduction(int(math.Pow(-1, float64(k-1)))*int(b)*int(Q), 1, int(module))
	result = append(result, uint(zeroX))
	log.Printf("x₀ = %d", zeroX)
	for counter := 1; counter < DCM; counter++ {
		currentX := zeroX + counter*int(module)
		result = append(result, uint(currentX))
		log.Printf("Current solution: %d", currentX)
	}
	return
}

func SimpleCompasionSolution(a int, b int, m int) (solutions []int, err error) {
	var LCD []int
	aHistory := []int{a}
	bHistory := []int{b}
	mHistory := []int{m}
	for {
		fmt.Printf("\nExpression: %d * x ≡ %d mod(%d)\n", aHistory[len(aHistory)-1], bHistory[len(bHistory)-1], mHistory[len(mHistory)-1])
		currentLCD := ExtendedEuclideanAlgorithm(aHistory[len(aHistory)-1], mHistory[len(mHistory)-1])
		LCD = append(LCD, currentLCD)
		fmt.Printf("LCD(%d, %d) = %d\n", aHistory[len(aHistory)-1], mHistory[len(mHistory)-1], currentLCD)
		if b%currentLCD == 0 {
			if len(LCD) == 1 {
				fmt.Printf("%d | %d -> the comparison is solvable with %d solutions\n", currentLCD, bHistory[len(bHistory)-1], currentLCD)
			} else {
				fmt.Printf("%d | %d -> use Euler's theorem\n", currentLCD, bHistory[len(bHistory)-1])
			}
		} else {
			err = fmt.Errorf("compasion hasn't got solution")
			return
		}
		if currentLCD == 1 {
			break
		}
		aHistory = append(aHistory, a/currentLCD)
		bHistory = append(bHistory, b/currentLCD)
		mHistory = append(mHistory, m/currentLCD)
	}
	fmt.Printf("x0 = (%d ^ (f(%d) - 1)) * %d mod(%d)\n", aHistory[len(aHistory)-1], mHistory[len(mHistory)-1], bHistory[len(bHistory)-1], mHistory[len(mHistory)-1])
	power := GetEulersFunction(mHistory[len(mHistory)-1]) - 1
	fmt.Printf("x0 = (%d ^ %d) * %d mod(%d)\n", aHistory[len(aHistory)-1], power, bHistory[len(bHistory)-1], mHistory[len(mHistory)-1])
	var expression, powerInt = big.NewInt(int64(aHistory[len(aHistory)-1])), big.NewInt(int64(power))
	result := expression.Exp(expression, powerInt, nil)
	result.Mul(result, big.NewInt(int64(bHistory[len(bHistory)-1])))
	_, x0 := new(big.Int).DivMod(result, big.NewInt(int64(mHistory[len(mHistory)-1])), new(big.Int))
	fmt.Printf("x0 = %d mod(%d) = %d\n", result, mHistory[len(mHistory)-1], x0)
	solutions = append(solutions, int(x0.Int64()))
	for counterIterations := 1; counterIterations < LCD[0]; counterIterations++ {
		currentRoot := x0.Int64() + int64(counterIterations)*int64(mHistory[1])
		solutions = append(solutions, int(currentRoot))
		fmt.Printf("Root № %d: %d + %d * %d = %d\n", counterIterations, x0.Int64(), int64(counterIterations), int64(mHistory[1]), currentRoot)
	}
	return
}
