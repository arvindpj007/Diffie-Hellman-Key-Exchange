package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
)

//Funciton to throw an error when the input CLI has missing/wrong parameters
func missingParametersError() {

	fmt.Println("ERROR: Parameters missing!")
	fmt.Println("HELP:")
	fmt.Println("./dh-bob <filename of message from Alice> <filename of message back to Alice>")

}

//Funciton to setup the CLI
func setupCLI() (string, string) {

	if len(os.Args) < 3 {

		missingParametersError()
		os.Exit(1)
	}

	input1 := os.Args[1]
	input2 := os.Args[2]
	return input1, input2

}

//Function to get the binary value from the given input file and returns value
func getInputText(inputText string) string {

	file, err := os.Open(inputText)
	if err != nil {
		log.Fatal(err)
	}

	dataBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	text := string(dataBytes)

	return text
}

// Function to write the text content into the output file
func setOutputText(text, output string) {

	var _, err = os.Stat(output)

	// Delete file if exists
	if os.IsExist(err) {

		err = os.Remove(output)
		if err != nil {
			log.Fatal(err)
			fmt.Println("ERROR: cannot open: ", err)
		}

	}

	// Create file
	file, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
		fmt.Println("ERROR: cannot open: ", err)
	}

	// Open file in append mode
	file, err = os.OpenFile(output, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
		fmt.Println("ERROR: cannot open: ", err)
	}

	// Write content in file
	l, err := file.WriteString(text)
	if err != nil {
		fmt.Println("ERROR: cannot write", err)
		file.Close()
		return
	}
	if l < 0 {

	}
	// fmt.Println(l, "bits written successfully to the file", output)
	file.Sync()
	file.Close()
}

// Function that generates A and also G power A mod P
func getBandGb(P, G *big.Int) (*big.Int, *big.Int) {

	var err error

	B := new(big.Int).SetInt64(0)
	Gb := new(big.Int).SetInt64(0)
	zero := new(big.Int).SetInt64(0)
	two := new(big.Int).SetInt64(2)
	one024 := new(big.Int).SetInt64(1024)
	maxInt := new(big.Int).SetInt64(0)

	checkB := true

	// MAX integer
	maxInt.Exp(two, one024, nil)

	B, err = rand.Int(rand.Reader, maxInt)

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(0)
	}

	for checkB {

		if B.Cmp(zero) == 0 {
			continue
		}

		B.Mod(B, P)

		if B.Cmp(zero) == 0 {
			continue
		}

		if B.Cmp(G) == 0 {
			continue
		}

		checkB = false

	}

	Gb.Exp(G, B, P)

	return B, Gb
}

// Function that outputs formatted text from the provided integer inputs
func getTextFromInts(P, G, A *big.Int) string {

	var text string

	text = "( "
	text = text + fmt.Sprintf("%s", P)
	text = text + ","
	text = text + fmt.Sprintf("%s", G)
	text = text + ","
	text = text + fmt.Sprintf("%s", A)
	text = text + " )"

	return text
}

//
func calculateSharedSecret(Ga, B, P *big.Int) *big.Int {

	Gab := new(big.Int)
	Gab.Exp(Ga, B, P)

	return Gab
}

func getParameters(fileBob string) (*big.Int, *big.Int, *big.Int) {

	P := new(big.Int)
	G := new(big.Int)
	Ga := new(big.Int)

	fileBobText := getInputText(fileBob)
	fileBobText = fileBobText[2 : len(fileBobText)-2]
	fileBobTextSplit := strings.Split(fileBobText, ",")
	PText := fileBobTextSplit[0]
	GText := fileBobTextSplit[1]
	GaText := fileBobTextSplit[2]

	P.SetString(PText, 10)
	G.SetString(GText, 10)
	Ga.SetString(GaText, 10)

	return P, G, Ga
}

//Main Function
func main() {

	var fileBob string
	// var fileBobText string
	var fileAlice string
	// var fileAliceText string

	fileBob, fileAlice = setupCLI()

	// fmt.Println(fileBob, fileAlice)

	P, G, Ga := getParameters(fileBob)

	B, Gb := getBandGb(P, G)
	// fmt.Println("Q:", Q)
	// fmt.Println("P:", P)
	// fmt.Println("G:", G)
	// fmt.Println("A:", A)
	// fmt.Println("Ga:", Ga)

	fileAliceText := getTextFromInts(P, G, Gb)

	// fmt.Println("secretKeyAlice: ", fileAliceText)

	setOutputText(fileAliceText, fileAlice)

	Gab := calculateSharedSecret(Ga, B, P)

	fmt.Println(Gab)
}
