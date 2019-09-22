package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
)

var one = big.NewInt(1)

type PublicKey struct {
	N uint64
	G uint64
}

var pubKey = getPubKey()
var maxVotes = 7
var currentVotes = 0
var encryptedVotes []int

func main() {
	http.HandleFunc("/vote", receiveVote)

	log.Fatal(http.ListenAndServe(":8888", nil))
}

func receiveVote(w http.ResponseWriter, r *http.Request) {
	if currentVotes < maxVotes {
		currentVotes++
		receivedVote, err := strconv.Atoi(r.FormValue("c"))
		if err != nil {
			panic(err)
		}
		fmt.Println("received encrypted vote: " + strconv.Itoa(receivedVote))
		encryptedVotes = append(encryptedVotes, receivedVote)
		if currentVotes == maxVotes {
			fmt.Println(encryptedVotes)
			sumVotes()
			currentVotes = 0
			encryptedVotes = nil
		}
	}
}

func sumVotes() {
	currentMult := big.NewInt(0)

	for i := 0; i < maxVotes-1; i++ {
		if currentMult.Cmp(big.NewInt(0)) == 0 {
			num1 := big.NewInt(int64(encryptedVotes[i]))
			num2 := big.NewInt(int64(encryptedVotes[i+1]))

			currentMult = new(big.Int).Mul(num1, num2)
		} else {
			num := big.NewInt(int64(encryptedVotes[i+1]))

			currentMult = new(big.Int).Mul(currentMult, num)
		}
	}

	nSquared := new(big.Int).Exp(big.NewInt(int64(pubKey.N)), big.NewInt(2), nil)

	cSum := new(big.Int).Mod(currentMult, nSquared)

	fmt.Println(cSum.String())

	resp, err := http.PostForm("http://localhost:8080/decrypt",
		url.Values{"c": {cSum.String()}})
	if err != nil {
		// panic(err)
		fmt.Println(err)
	}
	defer resp.Body.Close()
}

func getPubKey() PublicKey {
	pubKey := PublicKey{}

	res, err := http.Get("http://localhost:8080/pubkey")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(body, &pubKey)
	if err != nil {
		fmt.Println(err)
	}

	return pubKey
}
