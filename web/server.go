package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cbroglie/mustache"
	"github.com/folio-org/ldp-testdata/testdata"
)

// Obj is a JSON object; it can hold a string
type Obj map[string]interface{}

// Array is a slice of JSON objects
type Array []map[string]interface{}

var fileDefs Array // cache the filedefs on server startup

func viewHandler(w http.ResponseWriter, r *http.Request) {
	wrapper :=
		Obj{
			"defs": fileDefs, // mustache requires a field ('defs') to access the list
		}
	str, _ := mustache.RenderFile("web/index.html", wrapper)
	fmt.Fprintf(w, str)
}
func fakeHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[len("/fake/"):]
	filepath := "output/default/" + filename
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	wrapper :=
		Obj{
			"data": string(b),
			"file": filename,
		}
	jsonObj, marshalErr := json.Marshal(wrapper)
	if marshalErr != nil {
		panic(marshalErr)
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonObj))
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
	http.HandleFunc("/fake/", fakeHandler)
	if openInBrowser {
		go open("http://localhost:8080/")
	}
	panic(http.ListenAndServe(":8080", nil))
}
