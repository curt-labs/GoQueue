package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-nsq"
	"net/http"
	"strings"
)

const (
	AdminPartIndexUrl     = "http://iapi.curtmfg.com/index/part/"
	AdminCategoryIndexUrl = "http://iapi.curtmfg.com/index/category/"
)

type AdminHandler struct {
	ModificationType string `json:"modification_type"`
	ChangeType       string `json:"change_type"`
	Identifier       string `json:"id"`
}

func (a *AdminHandler) HandleMessage(message *nsq.Message) error {

	err := json.Unmarshal(message.Body, &a)
	if err != nil {
		return err
	}

	switch strings.ToLower(a.ModificationType) {
	case "part":
		err = a.index(AdminPartIndexUrl)
	case "category":
		err = a.index(AdminCategoryIndexUrl)
	}

	if err != nil {
		return err
	}

	message.Finish()

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

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("index called failed with %d", res.StatusCode)
	}

	return nil
}
