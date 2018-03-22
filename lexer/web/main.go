package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bucketd/go-graphqlparser/lexer"
	"github.com/bucketd/go-graphqlparser/token"
	"github.com/go-chi/chi"
)

type Request struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables"`
}

func main() {
	r := chi.NewRouter()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var req Request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "oops", http.StatusInternalServerError)
			return
		}

		lxr := lexer.New(strings.NewReader(req.Query))

		for {
			tok := lxr.Scan()
			if tok.Type == token.EOF {
				break
			}

			fmt.Fprintf(w, "%s(%s)\n", tok.Type.String(), tok.Literal)
			//fmt.Printf("%s(%s)\n", tok.Type.String(), tok.Literal)
		}
	})

	http.ListenAndServe(":3000", r)
}
