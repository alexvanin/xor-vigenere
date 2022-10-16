package main

import "math"

var (
	alphabet    = []rune("АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ")
	alphabetLen = len(alphabet)

	alphabetFreq = []float64{ // took from https://dpva.ru/Guide/GuideUnitsAlphabets/Alphabets/FrequencyRuLetters/ (without Ë)
		0.080159, // А
		0.015942, // Б
		0.045400, // ...
		0.016957,
		0.029801,
		0.084523,
		0.009398,
		0.016492,
		0.073559,
		0.012090,
		0.034952,
		0.044013,
		0.032080,
		0.066997,
		0.109714,
		0.028117,
		0.047352,
		0.054698,
		0.062606,
		0.026225,
		0.002645,
		0.009710,
		0.004829,
		0.014453,
		0.007283,
		0.003608,
		0.000367,
		0.018999,
		0.017392,
		0.003188,
		0.006377,
		0.020074,
	}

	mostFrequest = []rune("ОЕАИНТСРВЛК")
)

type Coder struct {
	ToCode map[rune]int
	ToRune map[int]rune
	RuFreq map[rune]float64
}

func NewCoder() Coder {
	coder := Coder{
		ToCode: make(map[rune]int, alphabetLen),
		ToRune: make(map[int]rune, alphabetLen),
		RuFreq: make(map[rune]float64, alphabetLen),
	}

	for i, r := range alphabet {
		coder.ToCode[r] = i
		coder.ToRune[i] = r
		coder.RuFreq[r] = alphabetFreq[i]
	}

	return coder
}

// Code with xor encoding.
func (c Coder) Code(text, key []rune) []rune {
	res := make([]rune, 0, len(text))
	for i, r := range text {
		keyIndex := i % len(key)
		res = append(res, rune(c.ToRune[c.ToCode[r]^c.ToCode[key[keyIndex]]]))
	}
	return res
}

func (c Coder) RuneFrequency(text, key []rune) map[rune]int {
	res := make(map[rune]int, alphabetLen)
	for i, r := range text {
		keyIndex := i % len(key)
		res[rune(c.ToRune[c.ToCode[r]^c.ToCode[key[keyIndex]]])] += 1
	}
	return res
}

func (c Coder) AlphabetDivergence(analysis map[rune]int, textLen float64) float64 {
	var result float64
	for r, freq := range analysis {
		result += math.Abs(float64(freq)/textLen - c.RuFreq[r])
	}
	return result
}
