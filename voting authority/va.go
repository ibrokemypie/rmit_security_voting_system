package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

var one = big.NewInt(1)

type PublicKey struct {
	N uint64
	G uint64
}

type PrivateKey struct {
	Lambda uint64
	Meu    uint64
}

func lcm(m *big.Int, n *big.Int) *big.Int {
	var z big.Int
	z.Mul(z.Div(m, z.GCD(nil, nil, m, n)), n)
	return &z
}

func l(u *big.Int, n *big.Int) *big.Int {
	var uMinusOne = new(big.Int).Sub(u, one)
	var divN = new(big.Int).Div(uMinusOne, n)
	return divN
}

var pubKey PublicKey
var privKey PrivateKey

func main() {
	var p = big.NewInt(89)
	var q = big.NewInt(53)
	var n = new(big.Int).Mul(p, q)
	var g = big.NewInt(8537)
	var nSquared = new(big.Int).Exp(n, big.NewInt(2), nil)

	pubKey = PublicKey{n.Uint64(), g.Uint64()}
	fmt.Println(pubKey)

	privKey = calculatePrivKey(p, q, n, g, nSquared)
	fmt.Println(privKey)

	http.HandleFunc("/pubkey", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pubKey)
	})

	http.HandleFunc("/decrypt", decrypt)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func decrypt(w http.ResponseWriter, r *http.Request) {
	c, err := strconv.Atoi(r.FormValue("c"))
	if err != nil {
		// panic(err)
		fmt.Println(err)
	}

	n := big.NewInt(int64(pubKey.N))
	cToLambda := new(big.Int).Exp(big.NewInt(int64(c)), big.NewInt(int64(privKey.Lambda)), nil)
	nSquared := new(big.Int).Exp(n, big.NewInt(2), nil)
	cLambdaModNSquared := new(big.Int).Mod(cToLambda, nSquared)
	lOut := l(cLambdaModNSquared, n)
	meuModN := new(big.Int).Mod(big.NewInt(int64(privKey.Meu)), n)
	temp := new(big.Int).Mul(lOut, meuModN)
	m := new(big.Int).Mod(temp, n)

	fmt.Println("Vote tally: " + strconv.FormatInt(m.Int64(), 2))
}

func calculatePrivKey(p *big.Int, q *big.Int, n *big.Int, g *big.Int, nSquared *big.Int) PrivateKey {
	var lambda = lcm(new(big.Int).Sub(p, one), new(big.Int).Sub(q, one))
	var gToLambda = new(big.Int).Exp(g, lambda, nSquared)
	var k = l(gToLambda, n)
	var meu = new(big.Int).Exp(k, big.NewInt(-1), n)

	return PrivateKey{lambda.Uint64(), meu.Uint64()}
}
