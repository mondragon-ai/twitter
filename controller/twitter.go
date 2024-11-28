package controller

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/twitter/data/request"
	"github.com/twitter/data/response"
	"github.com/twitter/helper"
	"github.com/twitter/service"
)

type TwitterController struct {
	TwitterService service.TwitterService
}

func NewTwitterController(twitterService service.TwitterService) *TwitterController {
	return &TwitterController{TwitterService: twitterService}
}

func (c *TwitterController) PostTweet(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	tweetCreate := request.TweetCreateRequest{}
	helper.ReadRequestBody(requests, &tweetCreate)

	// Call the Tweet service
	resp, err := c.TwitterService.PostTweet(requests.Context(), tweetCreate)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to post tweet: %s", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Respond with the status
	writer.WriteHeader(resp.StatusCode)
	writer.Write([]byte(fmt.Sprintf("Tweet posted successfully with status: %s", resp.Status)))
}

func (c *TwitterController) FetchMentions(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {

	// c.TwitterService.FetchMentions(requests.Context())
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)

}

func (c *TwitterController) ReplyMention(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	// thread := params.ByName("mentionId")

	// c.TwitterService.ReplyMention(requests.Context(), thread)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

func (controller *TwitterController) ReplyDM(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	// thread := params.ByName("thread")

	// c.TwitterService.ReplyDM(requests.Context(), thread)
	// if err != nil {
    //     http.Error(writer, fmt.Sprintf("Error: %s", err.Error()), http.StatusNotFound)
    //     return
    // }

    webResponse := response.WebResponse{
        Code:   http.StatusOK,
        Status: "Ok",
        Data:   nil,
    }

	helper.WriteResponseBody(writer, webResponse)
}