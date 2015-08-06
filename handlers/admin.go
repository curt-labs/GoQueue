package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-nsq"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	AdminPartIndexUrl     = "http://iapi.curtmfg.com/index/part/"
	AdminCategoryIndexUrl = "http://iapi.curtmfg.com/index/category/"
	PartIndexErrorUrl     = "http://iapi.curtmfg.com/index/part/error"
	PartIndexSuccessUrl   = "http://iapi.curtmfg.com/index/part/success"

	// AdminPartIndexUrl     = "http://localhost:8081/index/part/"
	// AdminCategoryIndexUrl = "http://localhost:8081/index/category/"
	// PartIndexErrorUrl     = "http://localhost:8081/index/part/error"
	// PartIndexSuccessUrl   = "http://localhost:8081/index/part/success"
)

type AdminHandler struct {
	ModificationType string        `json:"modification_type"`
	ChangeType       string        `json:"change_type"`
	Identifier       string        `json:"id"`
	TransitID        bson.ObjectId `json:"_id" bson:"_id"`
	Status           Status        `json:"status" bson:"status"`
	Error            string        `json:"error" bson:"error"`
}

type Status string

const (
	INTRANSIT Status = "In Transit"
	SUCCESS   Status = "Success"
	ERROR     Status = "Error"
)

func (a *AdminHandler) HandleMessage(message *nsq.Message) error {
	err := json.Unmarshal(message.Body, &a)
	if err != nil {
		return err
	}
	defer message.Finish()

	switch strings.ToLower(a.ModificationType) {
	case "part":
		a.Error = "" //Hmm...should already be empty....
		partIndexErr := a.index(AdminPartIndexUrl)
		err = a.updatePartIndexRecords(partIndexErr)
	case "category":
		a.Error = ""
		err = a.index(AdminCategoryIndexUrl)
	}

	if err != nil {
		return err
	}

	return nil
}

func (a *AdminHandler) updatePartIndexRecords(partIndexError error) error {
	url := PartIndexSuccessUrl
	if partIndexError != nil && partIndexError.Error() != "" {
		a.Error = partIndexError.Error()
		url = PartIndexErrorUrl
	}

	j, err := json.Marshal(&a)
	if err != nil {
		return err
	}
	res, err := http.Post(fmt.Sprintf("%s", url), "application/json", bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	body, err := (ioutil.ReadAll(res.Body))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("updatePartIndexRecords called failed with code: %d. %v", res.StatusCode, string(body))
	}
	return nil
}

func (a *AdminHandler) index(endpoint string) error {

	res, err := http.Get(fmt.Sprintf("%s%s", endpoint, a.Identifier))
	if err != nil {
		return err
	}
	if res == nil {
		return fmt.Errorf("%s", "request failed")
	}
	body, err := (ioutil.ReadAll(res.Body))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("index called failed with code: %d. %v", res.StatusCode, string(body))
	}

	return nil
}
