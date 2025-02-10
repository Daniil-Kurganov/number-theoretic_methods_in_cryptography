package numbertheoretic_methods_in_cryptography

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func systemPrint(system [][]int) {
	log.Println("\nSystem:")
	for _, currentEqution := range system {
		log.Printf(" %d * x ≡ %d mod(%d)\n", currentEqution[0], currentEqution[1], currentEqution[2])
	}
	log.Println()
}

func SolveCompasionSystem1Degree(powerSystem int) (solution int, err error) {
	log.Println("\nEquation has view like: c * x ≡ a mod(m)\nEnter \"c\", \"a\" and \"m\", separate by space")
	scanner := bufio.NewScanner(os.Stdin)
	var system, bigCoefficients [][]int
	for counterIteration := 0; counterIteration < powerSystem; counterIteration++ {
		fmt.Printf("Enter %d equation parameters: ", counterIteration+1)
		scanner.Scan()
		currentParametersString := strings.Split(scanner.Text(), " ")
		var currentParameters []int
		for _, currentParameterString := range currentParametersString {
			var currentParameter int
			if currentParameter, err = strconv.Atoi(currentParameterString); err != nil {
				log.Fatalf("Error on parsing current equation's parameters: %s", err)
			}
			currentParameters = append(currentParameters, currentParameter)
		}
		if currentParameters[0] > 1 {
			bigCoefficients = append(bigCoefficients, []int{counterIteration, currentParameters[0]})
		}
		system = append(system, currentParameters)
	}
	if len(bigCoefficients) != 0 {
		for _, currentCoefficients := range bigCoefficients {
			if system[currentCoefficients[0]][1]%currentCoefficients[1] == 0 && system[currentCoefficients[0]][2]%currentCoefficients[1] == 0 {
				system[currentCoefficients[0]][0] = 1
				system[currentCoefficients[0]][1] /= currentCoefficients[1]
				system[currentCoefficients[0]][2] /= currentCoefficients[1]
			} else {
				var currentSolutions []int
				if currentSolutions, err = SimpleCompasionSolution(system[currentCoefficients[0]][0], system[currentCoefficients[0]][1], system[currentCoefficients[0]][2]); err != nil {
					log.Fatalf("Error on getting solution of %d-equation: %s", currentCoefficients[0], err)
				}
				system[currentCoefficients[0]][1] = currentSolutions[0]
			}
			system[currentCoefficients[0]][0] = 1
			log.Printf("%d-equation became to: %d * x ≡ %d mod(%d)",
				currentCoefficients[0], system[currentCoefficients[0]][0], system[currentCoefficients[0]][1], system[currentCoefficients[0]][2])
		}
	}
	systemPrint(system)
	modulesMultiple := 1
	for currentIndex, currentEquation := range system {
		var secondIndex int
		if currentIndex == len(system)-1 {
			secondIndex = 0
		} else {
			secondIndex = currentIndex + 1
		}
		if ExtendedEuclideanAlgorithm(system[currentIndex][2], system[secondIndex][2]) != 1 {
			log.Fatal("The modules of the system are not pairwise coprime")
		}
		modulesMultiple *= currentEquation[2]
	}
	var solutionsCoefficients [][]int
	for currentIndex, currentEquation := range system {
		var currentSolutionsCoefficient []int
		currentM := modulesMultiple / currentEquation[2]
		currentSolutionsCoefficient = append(currentSolutionsCoefficient, currentM)
		log.Printf("N-%d = %d^(-1)mod(%d)\n", currentIndex+1, currentM, currentEquation[2])
		if ExtendedEuclideanAlgorithm(currentM, currentEquation[2]) != 1 {
			log.Fatalf("There is no inverse element for %d-equation", currentIndex+1)
		}
		log.Printf(" N-%d = %d^(f(%d) - 1)mod(%d)\n", currentIndex+1, currentM, currentEquation[2], currentEquation[2])
		currentPower := GetEulersFunction(currentEquation[2]) - 1
		log.Printf("Power: %d", currentPower)
		var currentExpression, currentPowerInt = big.NewInt(int64(currentM)), big.NewInt(int64(currentPower))
		currentMathResult := currentExpression.Exp(currentExpression, currentPowerInt, nil)
		_, currrentN := new(big.Int).DivMod(currentMathResult, big.NewInt(int64(currentEquation[2])), new(big.Int))
		log.Printf(" N-%d = %d mod(%d) = %d\n", currentIndex+1, currentMathResult, currentEquation[2], currrentN)
		currentSolutionsCoefficient = append(currentSolutionsCoefficient, int(currrentN.Int64()))
		currentSolutionsCoefficient = append(currentSolutionsCoefficient, currentEquation[1])

		solutionsCoefficients = append(solutionsCoefficients, currentSolutionsCoefficient)
	}
	log.Print("\nSolutions coefficients:")
	for currentIndex, currentCoefficients := range solutionsCoefficients {
		log.Printf("M-%d = %d; N-%d = %d", currentIndex+1, currentCoefficients[0], currentIndex+1, currentCoefficients[1])
		solution += currentCoefficients[0] * currentCoefficients[1] * currentCoefficients[2]
	}
	solution = solution % modulesMultiple
	log.Printf("\nResult: %d\n", solution)
	return
}
