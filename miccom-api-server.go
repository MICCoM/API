package main

import (
	"encoding/json"
	//"errors"
	"flag"
	"fmt"
	"github.com/MICCoM/API/MICCoM"
	"github.com/MICCoM/API/MICCoM/Experiment"
	"github.com/gorilla/mux"
	"github.com/wilke/RESTframe/CollectionJSON"
	"github.com/wilke/RESTframe/ShockClient"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	//"strconv"
)

type Item CollectionJSON.Item
type Collection CollectionJSON.Collection
type Client ShockClient.Client

var myURL url.URL
var baseURL string
var miccom MICCoM.MICCoM
var tmpDir string

var shock_ip = flag.String("shock", "http://localhost:7745", "URL for Shock host")
var mongo_ip = flag.String("mongo", "localhost:21071", "IP for MongoDB host")
var my_port = flag.String("port", "8001", "port for this service")
var my_url = flag.String("url", "http://localhost", "Display URL for this service")

func init() {

	flag.Parse()

	myURL.Host = "http://localhost:8001"
	baseURL = myURL.Host
	miccom.New(MICCoM.Parameter{ShockHost: *shock_ip,
		MongoHost: *mongo_ip,
		API:       *my_url})
	tmpDir = "/tmp"
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
	r.HandleFunc("/experiment", GetExperimentHandler).Methods("GET")
	r.HandleFunc("/experiment", CreateExperimentHandler).Methods("POST")
	r.HandleFunc(`/experiment/{ID:[a-zA-Z0-9\-]*}`, GetExperimentHandler).Name("experiment")
	// 	r.HandleFunc("/search", SearchHandler)
	// 	r.HandleFunc("/search/{path:.+}", SearchHandler)
	// 	r.HandleFunc("/register", RegisterHandler)
	// 	r.HandleFunc("/register/{path:[a-z+]+}", RegisterHandler)
	// 	r.HandleFunc("/upload", UploadHandler)
	// 	r.HandleFunc("/download", DownloadHandler)
	// 	r.HandleFunc("/transfer", TransferHandler)
	// 	r.HandleFunc("/transfer/{id}", SearchHandler)
	r.HandleFunc("/test", TestResourceHandler).Methods("POST")

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

	// Get URL Query Parameter
	error := r.ParseForm()
	if error != nil {
		fmt.Printf("ERROR: %+v", error)
	}

	// Defined Path Parameter
	vars := mux.Vars(r)
	id, ok := vars["ID"]

	// Build option set
	options := r.Form

	if ok {
		options["ID"] = []string{id}
	}

	c := miccom.GetExperiment(options)

	jb, err := json.Marshal(c)
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

// Create Experiments , POST request
func CreateExperimentHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Creating Experiment from POST")

	//Create collection object for response
	collection := Collection{Items: []Experiment.Experiment{}}

	// Get URL Query Parameter
	error := r.ParseForm()
	if error != nil {
		log.Println("Error, can't parse Form", error)
		fmt.Printf("ERROR: %+v", error)
	}

	// Defined Path Parameter
	vars := mux.Vars(r)
	id, ok := vars["ID"]
	if ok {
		fmt.Printf(id)
	}

	// Build option set
	options := r.Form
	options["success"] = []string{"0"}

	//Read multiparts from body

	log.Println("Reading Body")

	read_form, err := r.MultipartReader()
	if err != nil {
		log.Println("Error parsing multipart:", err)
		e := CollectionJSON.Error{Title: "Invalid form submitted",
			Code:    400,
			Message: err.Error()}
		collection.Error = &e

	} else {
		for {
			log.Println("Part")
			part, err_part := read_form.NextPart()

			// Error handling
			if err_part != nil {
				// reached end of stream
				if err_part == io.EOF {
					break
				} else {
					// Something went wrong
					// Return with error object
					log.Println("Problem", err_part)
				}
			}
			log.Println("Read Part")

			//Check for expected parts
			if part.FormName() == "file" {
				//Handle file upload
				if part.FileName() == "" {
					//Error, file part but no file upload
					log.Println("Creating Experiment but have 'file' part without filename ")
				} else {
					//save file to disk
					log.Println("Creating Experiment but have file - need item part")

					// Create FilePath , store in a tmp dir
					var tmpPath string = fmt.Sprintf("%s/%s", tmpDir, part.FileName())
					if tmpFile, err := os.Create(tmpPath); err == nil {
						defer tmpFile.Close()
						if _, err = io.Copy(tmpFile, part); err != nil {
							log.Println("Error", err)
						}
					}
				}
			} else if part.FormName() == "attributes" {
				//arbitrary json structure
				log.Println("Attributes ", part)
			} else if part.FormName() == "item" {
				//Collection Item
				log.Println("Creating Experiment from item part")

				// expect json; read content into string
				//decoder := json.NewDecoder(part)
				var d Experiment.Data

				slurp, err := ioutil.ReadAll(part)
				if err != nil {
					log.Fatal(err)
				}
				err = json.Unmarshal(slurp, &d)

				if err != nil {
				}

				//err = decoder.Decode(&d)

				if err != nil {
					log.Println(err)
					miccom.SendError(w, err, 400)
					//panic(err)
					fmt.Println("Yeah")
					return
				}

				log.Println("Got experiment: %s ", d.Name)

				e := Experiment.Experiment{Type: "Experiment", Data: d}

				collection.Items = append(collection.Items.([]Experiment.Experiment), e)

				// create experiment from json
				// create shock node and use as experiment id
				// update shock node

			} else {
				//Unexpected part
				fmt.Printf("Didn't see this part comming: %+v", part)
			}
			// END loop block
		}
		// END reading multipart
	}

	//var o map[string][]string

	//var c string = "{\"success\" : 0 }"
	//c := miccom.GetExperiment(o)

	jb, err := json.Marshal(collection)
	if err != nil {
		println(jb)
		w.Write([]byte(err.Error()))
		http.Error(w, err.Error(), 500)
	} else {
		fmt.Printf("Json: %s\n", jb)
		fmt.Printf("Last Error:  %+v\n", err)
		w.Header().Set("Content-Type", "application/json")
		if collection.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(jb))

	}
	//miccom.GetExperiment(w, r, miccom)
}

// Change name after test is successfull
func TestResourceHandler(w http.ResponseWriter, r *http.Request) {
	// File upload

	collection := Collection{}
	//file := Experiment.File{}

	fmt.Printf("Request: %+v\n", r)

	// Get URL Query Parameter
	error := r.ParseForm()
	if error != nil {
		fmt.Printf("ERROR: %+v", error)
	}

	// Defined Path Parameter
	vars := mux.Vars(r)
	id, ok := vars["ID"]
	if ok {
		fmt.Printf(id)
	}

	// Build option set
	options := r.Form
	options["success"] = []string{"0"}

	//Read multiparts from body
	read_form, err := r.MultipartReader()
	if err != nil {

	} else {
		for {
			part, err_part := read_form.NextPart()

			// Error handling
			if err_part != nil {
				// reached end of stream
				if err_part == io.EOF {
					break
				} else {
					// Something went wrong
					// Return with error object
				}
			}

			//Check for expected parts
			if part.FormName() == "file" {
				//Handle file upload
				if part.FileName() == "" {
					//Error, file part but no file upload
				} else {
					//save file to disk
					fmt.Printf("PART: %+v\n", part)
					var tmpPath string = fmt.Sprintf("/Users/Andi/Development/tmp/%s", part.FileName())
					if tmpFile, err := os.Create(tmpPath); err == nil {
						defer tmpFile.Close()
						if _, err = io.Copy(tmpFile, part); err != nil {
							log.Println("Error", err)
						}
					}
				}
			} else if part.FormName() == "attributes" {
				//arbitrary json structure
				fmt.Printf("PART: %+v\n", part)
				log.Println("Attributes ", part)
			} else if part.FormName() == "item" {
				//Collection Item
			} else {
				//Unexpected part
				fmt.Printf("Didn't see this part comming: %+v", part)
			}
			// END Loop block , reading part by part
		}
		// END Multipart reading block
	}

	// Store data
	// 1. Create Shock node for files
	// 2. Add file url/object to metadata
	// 3. Create experiment node

	// For every file in file list create file object and send to storage
	//if files > 0 {
	//}

	//var o map[string][]string

	//var c string = "{\"success\" : 0 }"
	//c := miccom.GetExperiment(o)

	jb, err := json.Marshal(collection)
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
