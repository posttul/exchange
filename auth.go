package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/posttul/exchange/storage"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// AuthServer handle the authetication of the users
type AuthServer struct {
	Store  storage.Storage
	tokens []Token
}

// Token holds base token information.
type Token struct {
	Token string `json:"token"`
	RPM   int    `json:"rps"` // Request Per Minute
}

func (t *Token) String() string {
	return fmt.Sprintf("token -> %s requestPerSecond -> %d ", t.Token, t.RPM)
}

// GetToken register an get token
func (a *AuthServer) GetToken() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t := Token{
			Token: generateToken(40),
			RPM:   20,
		}
		a.tokens = append(a.tokens, t)
		bts, err := json.Marshal(a.tokens)
		if err != nil {
			fail(w, err.Error())
		}
		a.Store.Write(bts)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"token":"%s"}`, t.Token)))
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
