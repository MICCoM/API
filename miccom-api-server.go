package main

import (
	"encoding/json"
	//"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wilke/webserver/CollectionJson"
	"github.com/wilke/webserver/MICCoM"
	//"github.com/wilke/webserver/MICCoM/Experiment"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Item CollectionJson.Item
type Collection CollectionJson.Collection

var myURL url.URL
var baseURL string
var miccom MICCoM.MICCoM

func init() {
	myURL.Host = "http://localhost:8000"
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
	r.HandleFunc("/experiment", ExperimentHandler)
	r.HandleFunc("/experiment/{id:[a-zA-Z]*}", ExperimentHandler).Name("experiment")
	r.HandleFunc("/search", SearchHandler)
	r.HandleFunc("/search/{path:.+}", SearchHandler)
	r.HandleFunc("/register", RegisterHandler)
	r.HandleFunc("/register/{path:[a-z+]+}", RegisterHandler)
	r.HandleFunc("/upload", UploadHandler)
	r.HandleFunc("/download", DownloadHandler)
	r.HandleFunc("/transfer", TransferHandler)
	r.HandleFunc("/transfer/{id}", SearchHandler)
	r.HandleFunc("/test", GetExperimentHandler)

	// Bind to a port and pass our router in
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Printf("%s\n", "Started Server at port 8000")
}

func GetExperimentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n", miccom)
	MICCoM.GetExperiment(w, r, miccom)
	//miccom.GetExperiment(w, r, miccom)
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("MICCoM!\n"))
	//fmt.Printf("%+v\n", c)
	fmt.Printf("%s\n", "Test")
}

func BaseHandler(w http.ResponseWriter, r *http.Request) {

	c := new(CollectionJson.CollectionJson)

	//q := CollectionJson.Query{}

	experiment_query := CollectionJson.Query{
		Href:   myURL.Host + "/experiment",
		Rel:    "experiment",
		Prompt: "Query definitions for experiment",
		Data:   nil,
	}

	c.Collection.Queries = []CollectionJson.Query{experiment_query}

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

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	path := vars["path"]

	//w.Write([]byte(r.UserAgent))
	//println(r.UserAgent())
	fmt.Printf("%+v\n", r.UserAgent())
	println(path)
	w.Write([]byte(r.UserAgent()))
	//c.version = 1
	//fmt.Printf("%+v\n", c)
	fmt.Printf("%s\n", "Test")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	path := vars["path"]

	href := myURL.Host + "/register" + "/" + path

	item := CollectionJson.Item{
		Href:  href,
		Data:  MICCoM.Experiment{}, //[]string{"Hallo", "Du"},
		Links: "no links",
		//Queries:  nil,
		//Template: nil,
	}

	// if path != "" {
	// 		alist := item.Data.([]string)
	// 		//alist.([]string)
	// 		//append(alist, path)
	// 		item.Data = append(alist, path)
	// 	}

	experiment := MICCoM.Experiment{}
	experiment.Data.ID = "A1001"

	if path != "" {
		alist := item.Data.([]MICCoM.Experiment)
		//alist.([]string)
		//append(alist, path)
		item.Data = append(alist, experiment)
		//item.AddData(experiment)
	}

	var list []MICCoM.Experiment
	col := CollectionJson.Collection{
		Version: "1",
		Href:    baseURL + r.URL.Path,
		Items:   list,
	}

	nr := experiment.AddToItems(&col)

	//list = col.Items.([]Frame.Item)
	//col.Items = append(list, item)

	fmt.Printf("%+v\n", col.Items)
	fmt.Printf("%+v\n", "Items: "+strconv.Itoa(nr))

	register := CollectionJson.CollectionJson{
		Collection: col,
	}

	println("Collection", register.Collection.Version)

	jb, err := json.Marshal(register)
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

func ExperimentHandler(w http.ResponseWriter, r *http.Request) {

	// Retrieve path and query parameters for experiment resource
	vars := mux.Vars(r)
	id := vars["id"]

	// initialize r.Form for query parameters
	r.ParseForm()

	method := r.Method
	switch {
	default:
		fmt.Printf("Found unsupported Method %s\n", r.Method)

	case method == "GET":
		fmt.Printf("Found Get\n")
		MICCoM.GetExperiment(w, r, miccom)
	case method == "POST":
		fmt.Printf("Found POST\n")
		c, err := miccom.CreateExperiment(r)
		if err != nil {
			miccom.SendError(w, err, 500)
		} else {
			fmt.Printf("Collection: %+v\n", c)
			miccom.SendCollection(w, c)
		}
	}

	return

	// initialize experiment
	experiment, err := MICCoM.NewExperiment(id)

	if err != nil {
		CollectionJson.SendError(w, err)
		return
		panic(err)
	}

	var list []MICCoM.Experiment
	col := CollectionJson.Collection{
		Version: "1",
		Href:    baseURL + r.URL.String(),
		Items:   list,
	}

	nr := experiment.AddToItems(&col)

	//list = col.Items.([]Frame.Item)
	//col.Items = append(list, item)

	fmt.Printf("%+v\n", col.Items)
	fmt.Printf("%+v\n", "Items: "+strconv.Itoa(nr))

	frame := CollectionJson.CollectionJson{
		Collection: col,
	}

	jb, err := json.Marshal(frame)
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

}

func OldExperimentHandler(w http.ResponseWriter, r *http.Request) {

	// Retrieve path and query parameters for experiment resource
	vars := mux.Vars(r)
	id := vars["id"]

	// initialize r.Form for query parameters
	r.ParseForm()

	// fmt.Printf("%s\n", "Length: "+strconv.Itoa(len(r.Form)))
	//
	// 	if id != "" {
	// 		w.Write([]byte(id))
	// 	}
	//
	// 	for key, value := range r.Form {
	// 		fmt.Printf("%s -> %s\n", key, value)
	// 		w.Write([]byte(key + "\n"))
	//
	// 	}
	//
	// 	for key, value := range vars {
	// 		fmt.Printf("%s :: %s\n", key, value)
	// 		w.Write([]byte(key + "\n"))
	//
	// 	}

	// initialize experiment
	experiment, err := MICCoM.NewExperiment(id)

	if err != nil {
		CollectionJson.SendError(w, err)
		return
		panic(err)
	}

	var list []MICCoM.Experiment
	col := CollectionJson.Collection{
		Version: "1",
		Href:    baseURL + r.URL.String(),
		Items:   list,
	}

	nr := experiment.AddToItems(&col)

	//list = col.Items.([]Frame.Item)
	//col.Items = append(list, item)

	fmt.Printf("%+v\n", col.Items)
	fmt.Printf("%+v\n", "Items: "+strconv.Itoa(nr))

	frame := CollectionJson.CollectionJson{
		Collection: col,
	}

	jb, err := json.Marshal(frame)
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

func UploadHandler(w http.ResponseWriter, r *http.Request)   {}
func DownloadHandler(w http.ResponseWriter, r *http.Request) {}
func TransferHandler(w http.ResponseWriter, r *http.Request) {}
