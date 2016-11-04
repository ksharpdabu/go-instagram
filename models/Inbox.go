package models

import (
	"encoding/json"
	"fmt"

	"github.com/hieven/go-instagram/constants"
	"github.com/hieven/go-instagram/utils"
	"github.com/parnurzeal/gorequest"
)

// Inbox type
type Inbox struct {
	Threads []*Thread `json:"threads"`
	Request *gorequest.SuperAgent
}

// Parsing instagram response
type feedResponse struct {
	Status  string
	Message string
	Inbox   Inbox `json:inbox`
}

type approveAllThreadRequestSchema struct {
	UUID string `json:"_uuid"`
}

// GetFeed returns you inbox feed
func (inbox *Inbox) GetFeed() (threads []*Thread) {
	_, body, _ := utils.WrapRequest(inbox.Request.Get(constants.ROUTES.Inbox))

	var resp feedResponse
	json.Unmarshal([]byte(body), &resp)

	if resp.Status == "fail" {
		fmt.Println(resp.Message)
		return
	}
	// fmt.Println(body)
	inbox.Threads = resp.Inbox.Threads

	for _, thread := range inbox.Threads {
		thread.Request = inbox.Request

		for _, item := range thread.Items {
			// if item.Location == {} {
			// 	continue
			// }

			item.Location.Request = inbox.Request
		}
	}

	return inbox.Threads
}

// ApproveAllThreads will approve all pending message requests
func (inbox *Inbox) ApproveAllThreads() {
	payload := approveAllThreadRequestSchema{
		UUID: utils.GenerateUUID(),
	}

	jsonData, _ := json.Marshal(payload)

	_, body, _ := utils.WrapRequest(
		inbox.Request.Post(constants.ROUTES.ThreadsApproveAll).
			Type("multipart").
			Send(string(jsonData)))

	fmt.Println(body)
}
