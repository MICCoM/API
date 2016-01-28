package CollectionJson

import (
	"encoding/json"
	"fmt"
	//"github.com/gorilla/mux"
	"net/http"
)

type Itemer interface {
	AddData(d interface{})
	AddItem(c *Collection) int
	AddToItems(c *Collection) int
	GetItem(i interface{}) interface{}
	ToData() []DataItem
}

type Templater interface {
	GetTemplate() Template
	ToTemplate() Template
}

type URL string

type Queries string
type Template []DataItem
type Error string

type DataItem struct {
	Name   string `json:"name" bson:"name"`
	Value  string `json:"value" bson:"value"` // should be string
	Prompt string `json:"prompt" bson:"prompt"`
}

//Query structure for elements in query list
type Query struct {
	Href   string
	Rel    string
	Prompt string
	Data   []DataItem
}

// Link structure
type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Prompt string `json:"prompt"`
	Name   string `json:"name"`
	Render string `json:"render"`
}

// Element in items list in a collection
type Item struct {
	Href  string      `json:"href"`
	Data  interface{} `json:"data"`
	Links string      `json:"links"`
	//Queries  string      `json:"queries"`
	//Template string      `json:"template"`
	//Error    Error       `json:"error"`
}

// Standard collection
type Collection struct {
	Version  string      `json:"version"`
	Href     string      `json:"href"`
	Links    []Link      `json:"href"`
	Items    interface{} `json:"items"`
	Queries  []Query     `json:"queries"`
	Template Template    `json:"template"`

	ID       int    `json:"ID"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Count    int    `json:"count"`
	Limit    int    `json:"limit"`
}

// Top level data structure, everything is a collection
type CollectionJson struct {
	Collection Collection `json:"collection"`
}

// Functions

func (c Collection) AddItem(i Item) {
	ilist := c.Items.([]Item)
	c.Items = append(ilist, i)
}

func (c CollectionJson) AddItem(i Item) {
	ilist := c.Collection.Items.([]Item)
	c.Collection.Items = append(ilist, i)
}

func (i Item) AddData(d Itemer) {

	alist := i.Data.([]interface{})
	i.Data = append(alist, d)

}

func SendError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}

func (c Collection) ToJson() ([]byte, error) {

	jb, err := json.Marshal(c)
	if err != nil {

	} else {
		fmt.Printf("%s\n", jb)
		fmt.Printf("%+v\n", err)

	}

	return jb, err
}

func (c CollectionJson) Send(w http.ResponseWriter) error {

	jb, err := c.ToJson()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jb))

	return nil
}
func (c CollectionJson) ToJson() ([]byte, error) {

	jb, err := json.Marshal(c)
	if err != nil {

	} else {
		fmt.Printf("%s\n", jb)
		fmt.Printf("%+v\n", err)

	}

	return jb, err
}
