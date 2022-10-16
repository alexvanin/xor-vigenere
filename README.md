# xor-vigenere

xor-encoded vigenere cypher decoder for ru-ru alphabet.
Decodes uppercase text without white spaces and Ë letter.

## Build

Requires [go1.19](https://go.dev/dl/)

```
$ go build .
```

## Usage

Specify input file and wait for result.

```
$ ./xor-vigenere -i example/cipher.txt -r 2
Cypher text file: example/cipher.txt
Parallel workers: 12 (override with -w flag)
Trying to predict secret key length ... 8 (override with -l flag)
Cypher text sample size to analyze: 2000 (override with -s flag)
Number of letters in permutations: 4 [ОЕАИ] (override with -p flag)
Number of results to show: 2 (override with -r flag)
Analyzing 100% [##################################################] (65536/65536, 36315 keys/s)        

Key=АБЫРВАЛГ Sample:ГДЕЖЕВЕСЬМИРВДЕНЬМОЕГОРОЖДЕНИЯГДЕЭЛЕКТРИЧЕСКИЕФО Div:0.145333
Key=АЙЫРВАЛГ Sample:ГМЕЖЕВЕСЬДИРВДЕНЬДОЕГОРОЖМЕНИЯГДЕХЛЕКТРИЧНСКИЕФО Div:0.166576
Execution duration: 1.805705397s
```

`АБЫРВАЛГ` is a valid secret key on top.

To receive better (but slower) results use:
- bigger sample size for analyzing cypher text (`-s 0` for analyzing whole text)
- look through more possible keys (`-r` flag)
- use more frequently used letters to create more permutations (slows **really** hard, `-p` flag)
