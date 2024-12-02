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

type MentionsController struct {
	MentionService service.MentionService
}

func NewMentionsController(mentionService service.MentionService) *MentionsController {
	return &MentionsController{MentionService: mentionService}
}

func (controller *MentionsController) Create(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	mentionCreateRequest := request.MentionCreateRequest{}
	helper.ReadRequestBody(requests, &mentionCreateRequest)

	controller.MentionService.Create(requests.Context(), mentionCreateRequest)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

func (controller *MentionsController) Delete(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	mentionId := params.ByName("mentionId")

	controller.MentionService.Delete(requests.Context(), mentionId)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

func (controller *MentionsController) FindAll(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	result := controller.MentionService.FindAll(requests.Context())
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   result,
	}

	helper.WriteResponseBody(writer, webResponse)
}

func (controller *MentionsController) FindById(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	mentionId := params.ByName("mentionId")

	result, err := controller.MentionService.FindById(requests.Context(), mentionId)
	if err != nil {
        http.Error(writer, fmt.Sprintf("Error: %s", err.Error()), http.StatusNotFound)
        return
    }

    webResponse := response.WebResponse{
        Code:   http.StatusOK,
        Status: "Ok",
        Data:   result,
    }

	helper.WriteResponseBody(writer, webResponse)
}





// Save Tweet Idea
func (controller *MentionsController) CreateTweetIdea(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	tweetIdeaRequest := request.TweetIdea{}
	helper.ReadRequestBody(requests, &tweetIdeaRequest)

	controller.MentionService.CreateTweetIdea(requests.Context(), tweetIdeaRequest)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

// Delete Tweet Idea
func (controller *MentionsController) DeleteTweetIdea(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	mentionId := params.ByName("mentionId")

	controller.MentionService.Delete(requests.Context(), mentionId)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

// FindAll Tweet Idea
func (controller *MentionsController) FindAllTweetIdeas(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	result := controller.MentionService.FindAll(requests.Context())
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   result,
	}

	helper.WriteResponseBody(writer, webResponse)
}