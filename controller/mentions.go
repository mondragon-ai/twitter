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
	ideaID := params.ByName("ideaID")

	controller.MentionService.DeleteTweetIdea(requests.Context(), ideaID)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

// FindAll Tweet Idea
func (controller *MentionsController) FindAllTweetIdeas(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	result := controller.MentionService.FindAllTweetIdea(requests.Context())
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   result,
	}

	helper.WriteResponseBody(writer, webResponse)
}







// Save Tweet Thread Idea
func (controller *MentionsController) CreateThreadIdea(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	threadRequest := request.ThreadIdea{}
	helper.ReadRequestBody(requests, &threadRequest)

	controller.MentionService.CreateThreadIdea(requests.Context(), threadRequest)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

// Delete Tweet Thread Idea
func (controller *MentionsController) DeleteThreadIdea(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	ideaID := params.ByName("ideaID")

	controller.MentionService.DeleteThreadIdea(requests.Context(), ideaID)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

// FindAll Tweet Thread Ideas
func (controller *MentionsController) FindAllThreadIdeas(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	result := controller.MentionService.FindAllThreadIdea(requests.Context())
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   result,
	}

	helper.WriteResponseBody(writer, webResponse)
}





// Save Tweet Clone Idea
func (controller *MentionsController) CreateTweetClone(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	cloneRequest := request.TweetClone{}
	helper.ReadRequestBody(requests, &cloneRequest)

	controller.MentionService.CreateTweetClone(requests.Context(), cloneRequest)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

// Delete Tweet Clone Idea
func (controller *MentionsController) DeleteTweetClone(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	cloneID := params.ByName("cloneID")

	controller.MentionService.DeleteTweetClone(requests.Context(), cloneID)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

// FindAll Tweet Clone Ideas
func (controller *MentionsController) FindAllTweetClones(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	result := controller.MentionService.FindAllTweetClone(requests.Context())
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   result,
	}

	helper.WriteResponseBody(writer, webResponse)
}





// Save Article Url
func (controller *MentionsController) CreateArticleUrl(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	articleRequest := request.UrlCreate{}
	helper.ReadRequestBody(requests, &articleRequest)

	controller.MentionService.CreateArticleUrl(requests.Context(), articleRequest)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

// Delete Article Url
func (controller *MentionsController) DeleteArticleUrl(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	articleID := params.ByName("articleID")

	controller.MentionService.DeleteArticleUrl(requests.Context(), articleID)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	helper.WriteResponseBody(writer, webResponse)
}

// FindAll Article Urls
func (controller *MentionsController) FindAllArticleUrls(writer http.ResponseWriter, requests *http.Request, params httprouter.Params) {
	result := controller.MentionService.FindAllArticleUrl(requests.Context())
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   result,
	}

	helper.WriteResponseBody(writer, webResponse)
}