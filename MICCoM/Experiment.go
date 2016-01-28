package MICCoM

import (
	//"fmt"
	"errors"
	"fmt"
	"github.com/MICCoM/API/CollectionJson"
	//"log"
	//"net/http"
	"time"
)

type Data struct {
	Name     string   `bson:"name"`
	ID       string   `bson:"ID"`
	Version  string   `bson:"version"`
	Date     string   `bson:"date"`
	Duration string   `bson:"duration"`
	Files    []string `bson:"files"`
	Samples  []string `bson:"samples"`
	Analysis interface{}
	Workflow interface{}
}

type Experiment struct {
	CollectionJson.Item
	Data Data `json:"data"`
}

//var template CollectionJson.Template{}

var experimentTemplate = CollectionJson.Template{
	{Name: "Name", Value: "String", Prompt: "Experiment name"},
	{Name: "Date", Value: "yyyy-mm-dd", Prompt: "Start date of experiment"},
	{Name: "Duration", Value: "integer", Prompt: "Duration of the experiment in seconds"}}

func NewExperiment(id string) (Experiment, error) {

	var e Experiment

	if id != "" {

		t := time.Now()
		fmt.Print(t)

		e = Experiment{}
		e.Data.ID = id
		e.Data.Version = "1"
		e.Data.Date = time.Now().Format(time.ANSIC)
	} else {
		return e, errors.New("Can't initialize experiment, no ID given")
	}

	return e, nil
}

func (e Experiment) GetTemplate() (CollectionJson.Template, error) {
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
	t := c.(CollectionJson.Collection)
	alist := t.Items.([]Experiment)
	t.Items = append(alist, e)

}

// Add experiment to items list in collection
func (e Experiment) AddToItems(c *CollectionJson.Collection) int {
	//var alist []ExperimentStruct

	if c.Items == nil {
		c.Items = []Experiment{e}
	} else {
		alist := c.Items.([]Experiment)
		c.Items = append(alist, e)
	}

	return len(c.Items.([]Experiment))
}

func (e Experiment) ToData() ([]CollectionJson.DataItem , error) {
	
	var dl  []CollectionJson.DataItem
	var err error
	
	
	return dl , err
}
