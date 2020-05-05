package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
)

var numberOfBitsQ = 1023

//Funciton to throw an error when the input CLI has missing/wrong parameters
func missingParametersError() {

	fmt.Println("ERROR: Parameters missing!")
	fmt.Println("HELP:")
	fmt.Println("./dh-alice1 <filename for message to Bob> <filename to store secret key>")

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

// Funciton that generates a random 1023 bit prime number Q
func generateQ() *big.Int {

	key, err := rand.Prime(rand.Reader, numberOfBitsQ)
	if err != nil {
		fmt.Println("Error")
		panic(err)
	}

	return key
}

// Funciton that discovers a random 1024 bit prime number P, such that P = 2Q + 1 and return P and Q
func getPandQ() (*big.Int, *big.Int) {

	one := new(big.Int).SetInt64(1)
	P := new(big.Int).SetInt64(0)
	Q := new(big.Int).SetInt64(0)
	isPrime := false

	for !isPrime {

		Q = generateQ()
		P.Add(Q, Q)
		P.Add(P, one)

		isPrime = P.ProbablyPrime(128)
	}

	return P, Q

}

// Function to find the primitive root or the generator element
func getGeneratorG(P, Q *big.Int) *big.Int {

	var err error

	G := new(big.Int).SetInt64(0)
	zero := new(big.Int).SetInt64(0)
	one := new(big.Int).SetInt64(1)
	two := new(big.Int).SetInt64(2)
	one024 := new(big.Int).SetInt64(1024)
	PminusOne := new(big.Int).SetInt64(0)
	powerOne := new(big.Int).SetInt64(0)
	powerTwo := new(big.Int).SetInt64(0)
	testOne := new(big.Int).SetInt64(0)
	testTwo := new(big.Int).SetInt64(0)
	maxInt := new(big.Int).SetInt64(0)

	checkGenereator := true

	// Factors of P - 1 are 2 and Q
	PminusOne.Add(Q, Q)

	// Test for primitive root requires raising the test element to the powers of (P-1)/2 and (P-1)/Q and if both return not equal to one it is a genreator element
	// (P-1)/2
	powerOne.Div(PminusOne, two)
	// (P-1)/Q
	powerTwo.Div(PminusOne, Q)
	// MAX integer
	maxInt.Exp(two, one024, nil)

	for checkGenereator {

		G, err = rand.Int(rand.Reader, maxInt)

		if err != nil {
			fmt.Println("error:", err)
			os.Exit(0)
		}

		if G.Cmp(zero) == 0 {
			continue
		}
		G.Mod(G, P)

		if G.Cmp(zero) == 0 {
			continue
		}

		testOne.Exp(G, powerOne, P)
		testTwo.Exp(G, powerTwo, P)

		// fmt.Println("Generator: ", G)
		// fmt.Println("One: ", testOne)
		// fmt.Println("Two: ", testTwo)

		if testOne.Cmp(one) == 0 || testTwo.Cmp(one) == 0 {
			checkGenereator = true
		} else {
			checkGenereator = false
		}

		// fmt.Println("Check: ", checkGenereator)
		// fmt.Println("Next:")

	}

	return G
}

// Function that generates A and also G power A mod P
func getAandGa(P, G *big.Int) (*big.Int, *big.Int) {

	var err error

	A := new(big.Int).SetInt64(0)
	Ga := new(big.Int).SetInt64(0)
	zero := new(big.Int).SetInt64(0)
	two := new(big.Int).SetInt64(2)
	one024 := new(big.Int).SetInt64(1024)
	maxInt := new(big.Int).SetInt64(0)

	checkA := true

	// MAX integer
	maxInt.Exp(two, one024, nil)

	A, err = rand.Int(rand.Reader, maxInt)

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(0)
	}

	for checkA {

		if A.Cmp(zero) == 0 {
			continue
		}

		A.Mod(A, P)

		if A.Cmp(zero) == 0 {
			continue
		}

		if A.Cmp(G) == 0 {
			continue
		}

		checkA = false

	}

	Ga.Exp(G, A, P)

	return A, Ga
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

//Main Function
func main() {

	var fileBob string
	var fileBobText string
	var secretKeyAlice string
	var secretKeyAliceText string

	fileBob, secretKeyAlice = setupCLI()

	// fmt.Println(fileBob, secretKeyAlice)

	P, Q := getPandQ()
	G := getGeneratorG(P, Q)
	A, Ga := getAandGa(P, G)

	fileBobText = getTextFromInts(P, G, Ga)
	secretKeyAliceText = getTextFromInts(P, G, A)

	// fmt.Println("fileBobText: ", fileBobText)
	// fmt.Println("secretKeyAlice: ", secretKeyAliceText)

	setOutputText(fileBobText, fileBob)
	setOutputText(secretKeyAliceText, secretKeyAlice)
}
