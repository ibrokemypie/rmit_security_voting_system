package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/therecipe/qt/widgets"
)

type PublicKey struct {
	N uint64
	G uint64
}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(250, 200)
	window.SetWindowTitle("Voting Booth")

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	input := widgets.NewQLineEdit(nil)
	input.SetPlaceholderText("Voter Number")

	widget.Layout().AddWidget(input)

	buttonGroup := widgets.NewQButtonGroup(nil)
	candidateOne := widgets.NewQRadioButton(nil)
	candidateOne.SetText("Candidate One")
	buttonGroup.AddButton(candidateOne, 1)
	candidateTwo := widgets.NewQRadioButton(nil)
	candidateTwo.SetText("Candidate Two")
	buttonGroup.AddButton(candidateTwo, 2)

	groupbox := widgets.NewQGroupBox(nil)
	groupbox.SetLayout(widgets.NewQVBoxLayout())
	groupbox.Layout().AddWidget(candidateOne)
	groupbox.Layout().AddWidget(candidateTwo)

	widget.Layout().AddWidget(groupbox)

	button := widgets.NewQPushButton2("Submit Vote", nil)
	button.ConnectClicked(func(bool) {
		var candidate = buttonGroup.CheckedId()
		var r, err = strconv.Atoi(input.Text())

		if err == nil && candidate != -1 {
			submitVote(r, candidate)
		}

	})
	input.ConnectReturnPressed(button.Click)
	widget.Layout().AddWidget(button)

	window.Show()
	app.Exec()
}

func submitVote(r int, candidate int) {
	var m int
	if candidate == 1 {
		m = 8
	} else if candidate == 2 {
		m = 1
	}

	fmt.Println("m:" + strconv.Itoa(m) + "r:" + strconv.Itoa(r))

	encryptVote(m, r)
}

func encryptVote(message int, random int) {
	pubKey := getPubKey()

	g := big.NewInt(int64(pubKey.G))
	n := big.NewInt(int64(pubKey.N))
	m := big.NewInt(int64(message))
	r := big.NewInt(int64(random))

	temp1 := new(big.Int).Exp(g, m, nil)
	temp2 := new(big.Int).Exp(r, n, nil)
	temp3 := new(big.Int).Mul(temp1, temp2)

	c := new(big.Int).Mod(temp3, n.Exp(n, big.NewInt(2), nil))

	fmt.Println(c)

	resp, err := http.PostForm("http://localhost:8888/vote",
		url.Values{"c": {c.String()}})
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
