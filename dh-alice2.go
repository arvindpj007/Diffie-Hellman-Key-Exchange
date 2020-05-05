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
	fmt.Println("./dh-alice2 <filename of message from Bob> <filename to read secret key>")

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

	// fmt.Println(binaryText)
	return text
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

	var fileAlice string
	// var fileBobText string
	var secertKeyAlice string
	// var fileAliceText string

	fileAlice, secertKeyAlice = setupCLI()

	// fmt.Println(fileAlice, secertKeyAlice)

	P, G, Gb := getParameters(fileAlice)
	P, G, A := getParameters(secertKeyAlice)

	G.Add(G, G)
	Gab := calculateSharedSecret(Gb, A, P)

	fmt.Println(Gab)
}
