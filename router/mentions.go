package router

import (
	"fmt"

	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/twitter/controller"
)

func MentionsRouter(mentionController *controller.MentionsController, twitterController *controller.TwitterController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprint(w, "Welcome Home")
	})

	router.GET("/api/mentions", mentionController.FindAll)
	router.GET("/api/mentions/:mentionId", mentionController.FindById)
	router.POST("/api/mentions", mentionController.Create)
	router.DELETE("/api/mentions/:mentionId", mentionController.Delete)


	router.POST("/api/twitter/tweet", twitterController.PostTweet)
	router.GET("/api/twitter/mentions", twitterController.FetchMentions)
	router.POST("/api/twitter/:mendionId", twitterController.ReplyMention)
	router.POST("/api/twitter/:id", twitterController.ReplyDM)

	return router
}