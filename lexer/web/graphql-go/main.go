package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/graphql-go/graphql/language/lexer"
	"github.com/graphql-go/graphql/language/source"
)

type Request struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables"`
}

func main() {
	r := chi.NewRouter()
	r.Post("/graphql", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		var body *source.Source

		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/graphql":
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "oopsie", http.StatusInternalServerError)
				return
			}

			body = source.NewSource(&source.Source{
				Body: bs,
			})
		case "application/json":
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "oops", http.StatusInternalServerError)
				return
			}

			body = source.NewSource(&source.Source{
				Body: []byte(req.Query),
			})
		default:
			// TODO(seeruk): Invalid content type.
			return
		}

		lxr := lexer.Lex(body)

		for {
			tok, err := lxr(0)
			if err != nil {
				http.Error(w, "oopsie", http.StatusInternalServerError)
				return
			}

			if tok.Kind == lexer.EOF {
				break
			}

			fmt.Fprintf(w, "%+v\n", tok)
			//fmt.Printf("%s(%s)\n", tok.Type.String(), tok.Literal)
		}
	})

	http.ListenAndServe(":3000", r)
}
