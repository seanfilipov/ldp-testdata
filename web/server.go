package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cbroglie/mustache"
	"github.com/folio-org/ldp-testdata/testdata"
)

type Obj map[string]interface{}
type Array []map[string]interface{}

var fileDefs Array

func viewHandler(w http.ResponseWriter, r *http.Request) {
	wrapper :=
		Obj{
			"defs": fileDefs,
		}
	str, _ := mustache.RenderFile("web/index.html", wrapper)
	fmt.Fprintf(w, str)
}

// Run the web server
func Run(openInBrowser bool, fDefs []testdata.FileDef) {
	jsonObj, marshalErr := json.Marshal(fDefs)
	if marshalErr != nil {
		panic(marshalErr)
	}
	if err := json.Unmarshal(jsonObj, &fileDefs); err != nil {
		panic(err)
	}

	http.HandleFunc("/", viewHandler)
	if openInBrowser {
		go open("http://localhost:8080/")
	}
	panic(http.ListenAndServe(":8080", nil))
}
