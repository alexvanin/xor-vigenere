package main

func groupText(text []rune, period int) []map[rune]int {
	res := make([]map[rune]int, period)
	for i := range res {
		res[i] = make(map[rune]int, alphabetLen)
	}

	for i, r := range text {
		index := i % period
		res[index][r] += 1
	}

	return res
}

func popularRune(m map[rune]int) rune {
	var (
		res rune
		max = 0
	)
	for l, counter := range m {
		if counter > max {
			res = l
			max = counter
		}
	}

	return res
}

func cartesian(runes [][]rune, processCombination func([]rune)) {
	c := 1
	for _, a := range runes {
		c *= len(a)
	}
	if c == 0 {
		return
	}

	b := make([]rune, c*len(runes))
	n := make([]int, len(runes))
	s := 0

	for i := 0; i < c; i++ {
		e := s + len(runes)
		pi := b[s:e]
		s = e
		for j, n := range n {
			pi[j] = runes[j][n]
		}
		for j := len(n) - 1; j >= 0; j-- {
			n[j]++
			if n[j] < len(runes[j]) {
				break
			}
			n[j] = 0
		}
		processCombination(pi)
	}
}
