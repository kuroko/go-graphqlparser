package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bucketd/go-graphqlparser/lexer"
	"github.com/bucketd/go-graphqlparser/token"
	"github.com/go-chi/chi"
)

type Request struct {
	Query         []byte                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables"`
}

func main() {
	r := chi.NewRouter()
	r.Post("/graphql", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		var lxr *lexer.Lexer

		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/graphql":
			bs, _ := ioutil.ReadAll(r.Body)

			lxr = lexer.New(bs)
		case "application/json":
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "oops", http.StatusInternalServerError)
				return
			}

			lxr = lexer.New(req.Query)
		default:
			// TODO(seeruk): Invalid content type.
			return
		}

		for {
			tok, _ := lxr.Scan()
			if tok.Type == token.EOF {
				break
			}

			fmt.Fprintf(w, "%s(%s)\n", tok.Type.String(), tok.Literal)
			//fmt.Printf("%s(%s)\n", tok.Type.String(), tok.Literal)
		}
	})

	http.ListenAndServe(":3000", r)
}
