package Experiment

import (
	//"fmt"
	"errors"
	"fmt"
	"github.com/wilke/RESTframe/CollectionJSON"
	//"log"
	//"net/http"
	"time"
)

const (
	DataType = "Experiment"
)

type File struct {
	Filename string   `json:"filename"`
	Path     string   `json:"path"`
	MD5      string   `json:"md5"`
	Size     string   `json:"size"`
	Type     string   `json:"type"`
	Format   string   `json:"format"`
	Tags     []string `json:"tags"`
	LRI      string   `json:"lri"`
	URI      string   `json:"uri"`
	ID       string   `json:"id"`
}

type Data struct {
	Type     string      `json:"type"`
	Name     string      `bson:"name" json:"name"`
	Codes    []string    `bson:"codes" json:"codes"`
	ID       string      `bson:"id"`
	Version  string      `bson:"version"`
	Date     string      `bson:"date"`
	Duration string      `bson:"duration"`
	Files    []File      `bson:"files"`
	Samples  []string    `bson:"samples"`
	Analysis interface{} `json:"analysis"`
	Workflow interface{} `json:"workflow"`

	// self.type = "Experiment"
	// self.name = None
	// self.version = None
	// self.codes = []
	// self.files = []
	// self.workflow = [] # list of steps
	// self.author = None
	// self.date = None
	// self.analysis = None
	// self.samples = {}

}

type Experiment struct {
	CollectionJSON.Item
	Type string `json:"type"`
	Data Data   `json:"data"`
}

//var template CollectionJSON.Template{}

var experimentTemplate = CollectionJSON.Template{
	{Name: "name", Value: "string", Prompt: "Experiment name"},
	{Name: "date", Value: "yyyy-mm-dd", Prompt: "Start date of experiment"},
	{Name: "duration", Value: "integer", Prompt: "Duration of the experiment in seconds"}}

func NewExperiment(id string) (Experiment, error) {

	var e Experiment

	if id != "" {

		t := time.Now()
		fmt.Print(t)

		e = Experiment{}
		e.Type = DataType
		e.Data.Type = DataType
		e.Data.ID = id
		e.Data.Version = "1"
		e.Data.Date = time.Now().Format(time.ANSIC)
	} else {
		return e, errors.New("Can't initialize experiment, no ID given")
	}

	return e, nil
}

func (e Experiment) GetTemplate() (CollectionJSON.Template, error) {
	template := experimentTemplate
	return template, nil
}

func (e Experiment) GetItem() Experiment {

	var err error

	e, err = NewExperiment(e.Data.ID)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}

	return e
}

func (e Experiment) AddToData(c interface{}) {
	//var alist []ExperimentStruct
	t := c.(CollectionJSON.Collection)
	alist := t.Items.([]Experiment)
	t.Items = append(alist, e)

}

// Add experiment to items list in collection
func (e Experiment) AddToItems(c *CollectionJSON.Collection) int {
	//var alist []ExperimentStruct

	if c.Items == nil {
		c.Items = []Experiment{e}
	} else {
		alist := c.Items.([]Experiment)
		c.Items = append(alist, e)
	}

	return len(c.Items.([]Experiment))
}

func (e Experiment) ToData() ([]CollectionJSON.DataItem, error) {

	var dl []CollectionJSON.DataItem
	var err error

	return dl, err
}
