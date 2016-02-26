package MICCoM

import (
	"fmt"
	"github.com/MICCoM/API/MICCoM/Experiment"
	//"github.com/wilke/webserver/Frame"
	//"errors"
	"github.com/wilke/RESTframe/CollectionJSON"
	"github.com/wilke/RESTframe/ShockClient"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"log"
	//"sync"
	"net/http"
	//"net/url"
	//"code.google.com/p/go-uuid/uuid"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	MongoDBHosts  = "localhost:27017"
	ShockHost     = "http://localhost:7445"
	AuthDatabase  = "miccom"
	AuthUserName  = "" //"miccom"
	AuthPassword  = "" //"miccom"
	MongoDatabase = "miccom"
	ApiUrl        = "http://localhost:8001"
)

type Parameter struct {
	MongoHost string
	MongoDB   string
	API       string
	User      string
	Password  string
	ShockHost string
}

type MICCoM struct {
	Api           string //ApiUrl
	MongoHost     string //MongoDBHosts
	MongoDB       string //MongoDatabase
	MongoUser     string //"miccom"
	MongoPassword string //"miccom"
	Mongo         *mgo.Session
	Shock         *ShockClient.Client
	ShockHost     string //Shock URL
}

type Sample struct{}
type Condition struct{}

var err error
var miccom MICCoM

// func init() {
//
// 	err = errors.New("")
//
// 	if err != nil {
// 		panic(err)
// 	}
// }

// intialize miccom struct with values
func (m *MICCoM) New(p Parameter) { //api string, mongoHost string, mongoDB string, user string, password string) {

	if p.API != "" {
		m.Api = p.API
	} else {
		m.Api = ApiUrl
	}

	if p.MongoHost != "" {
		m.MongoHost = p.MongoHost
	} else {
		m.MongoHost = MongoDBHosts
	}

	if p.MongoDB != "" {
		m.MongoDB = p.MongoDB
	} else {
		m.MongoDB = MongoDatabase
	}

	if p.User != "" {
		m.MongoUser = p.User
	} else {
		m.MongoUser = AuthUserName
	}
	if p.Password != "" {
		m.MongoPassword = p.Password
	} else {
		m.MongoPassword = AuthPassword
	}

	if m.MongoHost != "" {

		mongoDBDialInfo := &mgo.DialInfo{
			Addrs:    []string{m.MongoHost},
			Timeout:  60 * time.Second,
			Database: m.MongoDB,
			Username: m.MongoUser,
			Password: m.MongoPassword,
		}

		// Create a session which maintains a pool of socket connections
		// to our MongoDB.
		mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
		if err != nil {
			log.Fatalf("CreateSession: %s\n", err)
		} else {
			mongoSession.SetMode(mgo.Monotonic, true)
			m.Mongo = mongoSession
		}
	}

	if m.ShockHost == "" {
		var client ShockClient.Client
		h := ShockHost
		if p.ShockHost != "" {
			h = p.ShockHost
		}
		client.URL = h
		fmt.Printf("ShockClient %+v\n", client)
		m.ShockHost = h
		m.Shock = &ShockClient.Client{URL: h}
	} else {
		fmt.Printf("Something wrong\n")
	}

	fmt.Printf("Init miccom:\n%+v\n", m)
}

func (m MICCoM) GetExperiment(options map[string][]string) CollectionJSON.Collection {

	// Get options to query active store
	// var options map[string]string
	//var data Experiment.Data

	fmt.Printf("Options: %+v\n", options)

	collection := m.Get(options)

	c := CollectionJSON.Collection{}
	e := Experiment.Experiment{}
	c.Template, err = e.GetTemplate()

	fmt.Printf("Debug: %+v\n", c)
	var l []Experiment.Experiment
	//nodes := collection.Data.([]ShockClient.Node)
	if collection.Data != nil {
		for i, n := range collection.Data.([]ShockClient.Node) {

			fmt.Printf("Node %v\n", i)
			var attr map[string]interface{}
			if n.Attributes != nil {

				attr = n.Attributes.(map[string]interface{})
				var data Experiment.Data
				if attr["ID"] != nil {
					data.ID = attr["ID"].(string)
				}
				if attr["Type"] != nil {
					data.Type = attr["Type"].(string)
				} else {
					data.Type = "Experiment"
				}
				if attr["codes"] != nil {
					for _, c := range attr["codes"].([]interface{}) {
						data.Codes = append(data.Codes, c.(string))
					}
				}
				if attr["version"] != nil {
					data.Version = attr["version"].(string)
				}
				if attr["files"] != nil {
					for _, f := range attr["files"].([]interface{}) {
						data.Files = append(data.Files, map2file(f.(map[string]interface{})))
					}
				}

				e := Experiment.Experiment{Data: data}
				e.Href = m.Api + "/experiment/" + n.ID
				l = append(l, e)
			} else {
				e := Experiment.Experiment{}
				e.Href = m.Api + "/experiment/" + n.ID
				l = append(l, e)
			}
		}
	}
	c.Items = l
	c.Count = len(l)

	return c

	//jb, err := c.ToJson()
	//w.Write([]byte(jb))
	//w.Write([]byte("Got it"))
}

func (m MICCoM) CreateExperiment(r *http.Request) (*CollectionJSON.CollectionJSON, error) { //r *http.Request) {

	// Declare local variables
	// var c CollectionJSON.CollectionJSON // empty collection object

	fmt.Printf("Request %+v\n", r)

	// Init regexp for appliction/json
	re := regexp.MustCompile("application/json")

	// Read body if content type is application/json
	if re.Find([]byte(r.Header.Get("Content-Type"))) != nil {

		// Get length for read buffer
		length, err := strconv.Atoi(r.Header.Get("Content-Length"))

		// Define read buffer
		var p []byte
		p = make([]byte, length)

		a, err := r.Body.Read(p)
		jsonString := string(p[:])
		fmt.Printf("Body (JSON): %i , %+v\n", a, err)
		fmt.Printf("%s\n", jsonString)

		// Decode json string and create experiment item

		var t CollectionJSON.Template

		err = json.Unmarshal(p, &t)

		if err != nil {
			fmt.Println("Can not decode json string into template object:", err)
			return nil, err
		} else {
			fmt.Printf("%+v\n", t)
			expData := Experiment.Data{}
			for _, d := range t {
				fmt.Printf("Name: %+s\n", d.Name)
				if d.Name != "" {
					switch {
					default:
						fmt.Printf("Data %+v\n", expData)
						err := errors.New("Not supported value for attribute 'name' ( \"name\":" + d.Name + "\")\n")
						return nil, err
					case d.Name == "name":
						expData.Name = d.Value
					case d.Name == "ID":
						expData.ID = d.Value
					case d.Name == "Version":
						expData.Version = d.Value
					case d.Name == "Date":
						expData.Date = d.Value
					case d.Name == "Duration":
						expData.Duration = d.Value
					case d.Name == "Files":
						file := Experiment.File{ID: d.Value}
						expData.Files = []Experiment.File{file}
					case d.Name == "Samples":
						expData.Samples = []string{d.Value}
					}

				} else {
					fmt.Printf("Data %+v\n", expData)
					err := errors.New("Attribute 'name' has no value:\"name\":\"\" ")
					return nil, err
				}

			}

			e, err := Experiment.NewExperiment((time.Now()).String())
			e.Data = expData
			c := CollectionJSON.Collection{}
			nrItems := e.AddToItems(&c)

			// Set number of created experiments
			c.Count = nrItems
			c.Template, err = e.GetTemplate()

			// Mongo Check and Create
			sessionCopy := m.Mongo.Copy()
			defer sessionCopy.Close()

			// Create handle to Experiments collection
			collection := sessionCopy.DB(m.MongoDB).C("Experiments")

			//var experiments []Experiment
			query := collection.Find(&e)

			count, err := query.Count()
			if count == 0 {
				err := collection.Insert(&e)
				if err != nil {
					fmt.Printf("Error inserting document into mongo collection (%+v)\n", e)
					return nil, err
				}
				newExp := Experiment.Experiment{}
				err = collection.Find(&e).One(&newExp)

				if err != nil {
					fmt.Printf("Can't retrieve newly created document\n")
					return nil, err
				}

				nrItems := e.AddToItems(&c)
				fmt.Printf("CollectionJSON contains %d items (should 2)\n", nrItems)
			} else {
				fmt.Printf("Document already exists (%d) , please use update (%+v)\n", count, e)
			}

			// Set Outer collection frame and return
			cj := CollectionJSON.CollectionJSON{Collection: c}
			fmt.Printf("Returning new CollectionJSON: %+v\n", cj)
			return &cj, err
		}

	} else {
		fmt.Printf("Content-Type : %s\n", r.Header.Get("Content-Type"))
	}

	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	sessionCopy := m.Mongo.Copy()
	defer sessionCopy.Close()

	// // Get a mongo collection to execute the query against.
	// 	collection := sessionCopy.DB(m.MongoDB).C("Experiments")
	//
	// 	// Retrieve the list of experiments.
	// 	var experiments []Data
	// 	err := collection.Find(nil).All(&experiments)
	// 	if err != nil {
	// 		log.Printf("RunQuery : ERROR : %s\n", err)
	// 		return c, err
	// 	}
	//
	// 	c = CollectionJSON.CollectionJSON{}
	// 	fmt.Printf("Debug: %+v\n", c)
	// 	var l []Experiment
	// 	for _, d := range experiments {
	// 		e := Experiment{Data: d}
	// 		e.Href = m.Api + "/experiment/" + d.ID
	// 		l = append(l, e)
	// 	}
	// 	c.Collection.Items = l
	// 	c.Collection.Count = len(l)

	return nil, nil

}

// *********************
// http return functions

func (m MICCoM) SendCollection(w http.ResponseWriter, c *CollectionJSON.CollectionJSON) {
	jb, err := c.ToJson()
	if err != nil {
		code := 500
		m.SendError(w, err, code)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jb))
}

func (m MICCoM) SendError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

// **********************

func UpdateExperiment() {}

func (m MICCoM) mongo(*mgo.Session, error) {}

func (m MICCoM) Get(o map[string][]string) ShockClient.Collection {

	if m.Shock != nil {
		fmt.Printf("Test Shock\n")

		if m.ShockHost == "" {
			fmt.Printf("Error: No Host\n")
			var col ShockClient.Collection
			col.Status = 404
			col.Error = "Missing host in path"
			return col
		}

		uri := strings.Join([]string{m.ShockHost, "node"}, "/")

		id, ok := o["ID"]
		if ok {

			if len(id) > 1 {
				uri = strings.Join([]string{uri, "querynode"}, "?")
				for _, i := range id {
					uri = strings.Join([]string{uri, strings.Join([]string{"id", i}, "=")}, "&")
				}
			} else {
				uri = strings.Join([]string{uri, id[0]}, "/")
			}
		}

		collection, err_code, err := m.Shock.Get(uri)

		if err != nil {
			fmt.Printf("Error: %v (%v) ", err_code, err.Error)
			return collection
		}

		return collection

	} else if m.Mongo != nil {
		fmt.Printf("Test Mongo\n")

		// Request a socket connection from the session to process our query.
		// Close the session when the goroutine exits and put the connection back
		// into the pool.
		sessionCopy := m.Mongo.Copy()
		defer sessionCopy.Close()

		// Get a mongo collection to execute the query against.
		collection := sessionCopy.DB(m.MongoDB).C("Experiments")

		// Retrieve the list of experiments.
		var experiments []Experiment.Experiment
		err := collection.Find(nil).All(&experiments)
		if err != nil {
			log.Printf("RunQuery : ERROR : %s\n", err)
			return ShockClient.Collection{}
		}
	}
	return ShockClient.Collection{}
}

func (m MICCoM) testget(o interface{}) interface{} {

	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	sessionCopy := m.Mongo.Copy()
	defer sessionCopy.Close()

	// Get a mongo collection to execute the query against.
	collection := sessionCopy.DB(m.MongoDB).C("Experiments")

	// Retrieve the list of experiments.
	var experiments []Experiment.Experiment
	err := collection.Find(nil).All(&experiments)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
		return nil
	}

	return o
}

func map2file(a map[string]interface{}) Experiment.File {
	var f Experiment.File

	b, err := json.Marshal(a)
	err = json.Unmarshal(b, &f)

	if err != nil {
		log.Printf("map2file : ERROR : %s\n", err)
		return Experiment.File{}
	}

	//var list_of_keys []string
	// list_of_keys := []string{"filename", "path", "md5", "size", "Type", "format", "tags", "lri", "uri", "ID"}
	// 	for _, k := range list_of_keys {
	// 		f.(k = a[k].(string)
	// 	}
	return f
}
