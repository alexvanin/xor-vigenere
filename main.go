package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	outputSampleSize   = 6
	maxSecretKeySize   = 12
	expectedMinTextLen = 2000
)

func main() {
	input := flag.String("i", "", "cypher text")
	debug := flag.Bool("debug", false, "print debug info")
	secretKeyLen := flag.Int("l", 0, "length of secret key if known")
	permLettersLen := flag.Int("p", 4, "number of most frequest russian letters for cartesian permutation")
	workers := flag.Int("w", 0, "number of parallel workers (default number of CPU)")
	resultLen := flag.Int("r", 50, "number of top results to process")
	sampleSize := flag.Int("s", 2000, "size of cypher text sample to analyze, 0 to analyze whole text")
	flag.Parse()

	start := time.Now()

	// Read cypher text.
	data, err := os.ReadFile(*input)
	if err != nil {
		fmt.Printf("reading cypher text: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Cypher text file: %s\n", *input)
	cypherText := []rune(string(data))

	if len(cypherText) < expectedMinTextLen {
		fmt.Printf("! Text length %d < %d, frequency analysis may be affected", len(cypherText), expectedMinTextLen)
	}

	// Set number of parallel workers.
	if *workers == 0 {
		cpuN := runtime.NumCPU()
		workers = &cpuN
	}
	fmt.Printf("Parallel workers: %d (override with -w flag)\n", *workers)

	// Finding secret key length.
	if *secretKeyLen == 0 {
		fmt.Print("Trying to predict secret key length ...")
		prediction := predictSecretKeyLen(cypherText, maxSecretKeySize)
		fmt.Printf(" %d (override with -l flag)\n", prediction)
		secretKeyLen = &prediction
	} else {
		fmt.Printf("Secret key length: %d (override with -w flag)\n", *secretKeyLen)
	}

	// Set sample size to analyze crypto text.
	if *sampleSize == 0 {
		size := len(cypherText)
		sampleSize = &size
	}
	fmt.Printf("Cypher text sample size to analyze: %d (override with -s flag)\n", *sampleSize)

	// Print other parameters.
	fmt.Printf("Number of letters in permutations: %d [%s] (override with -p flag)\n", *permLettersLen, string(mostFrequest[:*permLettersLen]))

	//
	fmt.Printf("Number of results to show: %d (override with -r flag)\n", *resultLen)

	// Start decoding.
	coder := NewCoder()
	groups := groupText(cypherText, *secretKeyLen)

	permutations := make([][]rune, *secretKeyLen)
	for i := range permutations {
		permutations[i] = make([]rune, *permLettersLen)
		mostUsedRune := popularRune(groups[i])
		for j, r := range mostFrequest[:*permLettersLen] {
			permutations[i][j] = rune(coder.ToRune[coder.ToCode[mostUsedRune]^coder.ToCode[r]])
		}
	}

	possibleKeysLen := 1
	for _, p := range permutations {
		possibleKeysLen *= len(p)
	}

	result := NewTop(*resultLen)
	pb := newProgressBar(possibleKeysLen)

	jobs := make(chan []rune)
	wg := new(sync.WaitGroup)
	wg.Add(possibleKeysLen)

	for i := 0; i < *workers; i++ {
		go func(jobs <-chan []rune) {
			for job := range jobs {
				decodedFreq := coder.RuneFrequency(cypherText[:*sampleSize], job)
				result.Add(string(job), coder.AlphabetDivergence(decodedFreq, float64(*sampleSize)))
				pb.Add(1)
				wg.Done()
			}
		}(jobs)
	}

	cartesian(permutations, func(r []rune) {
		jobs <- r
	})
	wg.Wait()
	fmt.Println()

	// Print results.
	for _, s := range result.List() {
		fmt.Printf("Key=%s Sample:%s Div:%f\n", s, string(coder.Code(cypherText[:outputSampleSize**secretKeyLen], []rune(s))), result.Value(s))
	}

	fmt.Printf("Execution duration: %s\n", time.Since(start))
	if *debug {
		PrintMemUsage()
	}
}

func newProgressBar(cap int) *progressbar.ProgressBar {
	return progressbar.NewOptions(cap,
		progressbar.OptionSetDescription("Analyzing"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("keys"),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "#",
			SaucerHead:    "#",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

func predictSecretKeyLen(text []rune, limit int) int {
	lim := len(text) / 2
	if lim > limit {
		lim = limit
	}

	ln := len(text)

	var (
		predictedLen     int
		predictedOverlap int
	)
	for i := 1; i <= lim; i++ {
		overlap := 0
		for j := 0; j < ln; j++ {
			shiftIndex := (j + i) % ln
			if text[j] == text[shiftIndex] {
				overlap++
			}
		}
		if overlap > predictedOverlap {
			predictedOverlap = overlap
			predictedLen = i
		}
	}

	return predictedLen
}
