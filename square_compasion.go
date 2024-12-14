package numbertheoreticmethodsincryptography

import (
	"fmt"
	"log"
	"math"
)

func SolvingSquareCompasion(a int, power int) (result int, err error) {
	if !PrimaryCheck(power) {
		err = fmt.Errorf("power isn't primary")
		return
	}
	log.Printf("x² ≡ %d mod(%d)", a, power)
	entranceCheck := GetLegandreSymbol(a, power)
	if entranceCheck == -1 {
		log.Fatalf("\nL(%d, %d) = -1 -> The comparison is undecidable", a, power)
	}
	log.Printf("\nL(%d, %d) = 1 -> The comcomparison is solvable and has 2 solutions", a, power)
	if ModuloReduction(power, 1, 4) == 3 {
		factor := (power - 3) / 4
		result = ModuloReduction(a, factor+1, power)
		log.Printf("%d ≡ 3 mod(4) -> use the first simple case solution:\n x = ±%d^(%d + 1) mod(%d) ≡ ±%d mod(%d)",
			power, a, factor, power, result, power)
		return
	}
	if ModuloReduction(power, 1, 8) == 5 {
		factor := (power - 5) / 8
		selectWayResult := ModuloReduction(a, 2*factor+1, power)
		log.Printf("%d ≡ 5 mod(8) -> use the second simple case solution", power)
		if selectWayResult == 1 {
			result = ModuloReduction(a, factor+1, power)
			log.Printf(" first subcase: x = ±%d^(%d + 1) mod(%d) = ±%d mod(%d)", a, factor, power, result, power)
			return
		}
		if selectWayResult == power-1 {
			result = ModuloReduction(ModuloReduction(a, factor+1, power)*ModuloReduction(2, 2*factor+1, power), 1, power)
			log.Printf(" second subcase: x = %d^(%d + 1) * 2^(2 * %d + 1) mod(%d) = ±%d mod(%d)",
				a, factor, factor, power, result, power)
			return
		}
	}
	n, nCheck := 2, entranceCheck*(-1)
	for {
		currentCheck := GetLegandreSymbol(n, power)
		if currentCheck == power-1 {
			currentCheck = -1
		}
		if currentCheck == nCheck {
			break
		}
		n += 1
	}
	log.Printf("\nUse general algoritm solution\nL(%d, %d) = %d -> L(N, %d) = %d; N = %d", a, power, entranceCheck, power, nCheck, n)
	var h int
	k := 1
	for {
		if (power-1)%int(math.Pow(2, float64(k))) == 0 && ((power-1)/int(math.Pow(2, float64(k))))%2 == 1 {
			h = (power - 1) / int(math.Pow(2, float64(k)))
			break
		}
		k += 1
	}
	log.Printf("%d = 1 + 2^%d * %d -> k = %d, h = %d", power, k, h, k, h)
	a1 := ModuloReduction(a, (h+1)/2, power)
	var a2 int
	if a2, err = GetInverseElement(a, power); err != nil {
		err = fmt.Errorf("error on calculating a2: %s", err)
		return
	}
	n1 := ModuloReduction(n, h, power)
	n2 := 1
	log.Printf("\nStart values: a1 = %d, a2 = %d, N1 = %d, N2 = %d", a1, a2, n1, n2)
	var table [][]int // [[i, b, c, d, j, N2]]
	for currentIteration := 0; currentIteration < k-1; currentIteration++ {
		currentTableRow := []int{currentIteration}
		b := ModuloReduction(a1*n2, 1, power)
		currentTableRow = append(currentTableRow, b)
		c := ModuloReduction(a2*ModuloReduction(b, 2, power), 1, power)
		currentTableRow = append(currentTableRow, c)
		d := ModuloReduction(c, int(math.Pow(2, float64(k)-2-float64(currentIteration))), power)
		var j int
		if d == -1 || d == power-1 {
			d, j = -1, 1
		} else {
			j = 0
		}
		currentTableRow = append(currentTableRow, d, j)
		n2 = ModuloReduction(n2*ModuloReduction(n1, int(math.Pow(2, float64(currentIteration)))*j, power), 1, power)
		currentTableRow = append(currentTableRow, n2)
		table = append(table, currentTableRow)
	}
	log.Print("\nTable\n i b c d j N2\n")
	for _, currentRow := range table {
		log.Print(currentRow)
	}
	result = ModuloReduction(a1*n2, 1, power)
	log.Printf("x = ±%d mod(p)", result)
	return
}
