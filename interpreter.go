package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var keywords map[rune]uint8 = map[rune]uint8{
	'>': 0, //incPtr
	'<': 1, //decPtr
	'+': 2, //incCell
	'-': 3, //decCell
	'.': 4, //out
	',': 5, //in
	'[': 6, //jmpFwd
	']': 7, //jmpBwd
}

type oBracket struct {
	LoopNb   int
	Previous *oBracket
}

var pathToSrc *string
var ptr uint16 = 0
var cells []uint8

func setFlag() {
	pathToSrc = flag.String("filepath", "./helloworld.bf", "spécifier le chemin d'accès au fichier source")
	flag.Parse()
}

func loadFile() *os.File {
	if pathToSrc == nil {
		log.Fatalln("Veuillez spécifier un chemin d'accès au code source. -filepath=./chemin/vers/fichier")
	}
	file, err := os.Open(*pathToSrc)

	if err != nil {
		log.Fatalln("Veuillez spécifier un chemin d'accès au code source. (" + err.Error() + ")")
	}

	return file
}

func compileSource() (comp []uint8, jmpList [][2]int) {
	file := loadFile()
	defer file.Close()

	//reader reading from the source file
	reader := bufio.NewReader(file)
	loopNb := 0

	var lastOpened *oBracket = &oBracket{
		LoopNb:   0,
		Previous: nil,
	}
	pos := 0

	for {
		r, _, err := reader.ReadRune()

		if err == io.EOF {
			break
		}
		//transfer the source code from file to slice
		if val, ok := keywords[r]; ok {
			if val == 6 {
				jmpList = append(jmpList, [2]int{pos, 0})

				lastOpened = &oBracket{
					LoopNb:   loopNb,
					Previous: lastOpened,
				}

				loopNb++
			}
			if val == 7 {
				jmpList[lastOpened.LoopNb][1] = pos
				lastOpened = lastOpened.Previous
			}

			comp = append(comp, val)
			pos++
		}
	}

	return comp, jmpList
}

func exec(comp []uint8, jmpList [][2]int) {

	cells = make([]uint8, 65536)
	loopNb := 0
	var loops map[int]bool = make(map[int]bool)

	var lastOpened *oBracket = &oBracket{
		LoopNb:   0,
		Previous: nil,
	}

	for i := 0; i < len(comp); i++ {
		switch comp[i] {
		case 0:
			ptr++
		case 1:
			ptr--
		case 2:
			cells[ptr]++
		case 3:
			cells[ptr]--
		case 4:
			fmt.Print(string(cells[ptr]))
		case 5:
			cells[ptr] = userInput()
		case 6:

			lastOpened = &oBracket{
				LoopNb:   loopNb,
				Previous: lastOpened,
			}

			if cells[ptr] == 0 {
				i = jmpList[lastOpened.LoopNb][1]
				lastOpened = lastOpened.Previous
			}

			if _, ok := loops[i]; !ok {
				loops[i] = true
				loopNb++
			}
		case 7:
			if cells[ptr] != 0 {
				i = jmpList[lastOpened.LoopNb][0]
			}
		}
	}
}

func userInput() uint8 {
	var input string
	n, _ := fmt.Scan(&input)

	//making sure the user input is not empty
	if n == 0 {
		return 0
	}
	//taking the first character
	s := rune(input[0])

	//checking for overflow
	if s < 256 {
		return uint8(s)
	}

	return 255
}

func main() {
	setFlag()
	exec(compileSource())
}
