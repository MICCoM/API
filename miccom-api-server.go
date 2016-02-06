package main

import (
	"encoding/json"
	//"errors"
	"fmt"
	"github.com/MICCoM/API/MICCoM"
	"github.com/gorilla/mux"
	"github.com/wilke/RESTframe/CollectionJSON"
	"github.com/wilke/RESTframe/ShockClient"
	"log"
	"net/http"
	"net/url"
	//"strconv"
)

type Item CollectionJSON.Item
type Collection CollectionJSON.Collection
type Client ShockClient.Client

var myURL url.URL
var baseURL string
var miccom MICCoM.MICCoM

func init() {
	myURL.Host = "http://localhost:8001"
	baseURL = myURL.Host
	miccom.New("", "", "", "", "")
	// i = new(Item)
	// 	c = new(Frame.Collection)
	// 	fmt.Printf("%+v\n", c)
	// 	fmt.Printf("%s\n", "Test")
}

func main() {

	fmt.Printf("%s\n", "Starting Server")
	fmt.Printf("Miccom:\n%+v\n", miccom)
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.

	r.HandleFunc("/", BaseHandler)
	//r.HandleFunc("/experiment", ExperimentHandler)
	// r.HandleFunc("/experiment/{id:[a-zA-Z]*}", ExperimentHandler).Name("experiment")
	// 	r.HandleFunc("/search", SearchHandler)
	// 	r.HandleFunc("/search/{path:.+}", SearchHandler)
	// 	r.HandleFunc("/register", RegisterHandler)
	// 	r.HandleFunc("/register/{path:[a-z+]+}", RegisterHandler)
	// 	r.HandleFunc("/upload", UploadHandler)
	// 	r.HandleFunc("/download", DownloadHandler)
	// 	r.HandleFunc("/transfer", TransferHandler)
	// 	r.HandleFunc("/transfer/{id}", SearchHandler)
	r.HandleFunc("/test", GetExperimentHandler)

	// Bind to a port and pass our router in
	err := http.ListenAndServe(":8001", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Printf("%s\n", "Started Server at port 8000")
}

func BaseHandler(w http.ResponseWriter, r *http.Request) {

	c := new(CollectionJSON.CollectionJSON)

	//q := CollectionJSON.Query{}

	experiment_query := CollectionJSON.Query{
		Href:   myURL.Host + "/experiment",
		Rel:    "experiment",
		Prompt: "Query definitions for experiment",
		Data:   nil,
	}

	c.Collection.Queries = []CollectionJSON.Query{experiment_query}

	// Create json from collection
	jb, err := c.ToJson()

	// Send json
	if err != nil {
		println(jb)
		w.Write([]byte(err.Error()))
		http.Error(w, err.Error(), 500)
	} else {
		fmt.Printf("%s\n", jb)
		fmt.Printf("%+v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jb))

	}

}

func GetExperimentHandler(w http.ResponseWriter, r *http.Request) {

	var o interface{}
	o = miccom.Get(o)

	jb, err := json.Marshal(o)
	if err != nil {
		println(jb)
		w.Write([]byte(err.Error()))
		http.Error(w, err.Error(), 500)
	} else {
		fmt.Printf("Json: %s\n", jb)
		fmt.Printf("Last Error:  %+v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jb))

	}
	//miccom.GetExperiment(w, r, miccom)
}
