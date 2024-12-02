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

	// Tweet Clones
	router.GET("/api/clones", mentionController.FindAllTweetClones)
	router.POST("/api/clone", mentionController.CreateTweetClone)
	router.DELETE("/api/clone/:mentionId", mentionController.DeleteTweetClone)

	// Tweet Ideas
	router.GET("/api/ideas", mentionController.FindAllTweetIdeas)
	router.POST("/api/ideas", mentionController.CreateTweetIdea)
	router.DELETE("/api/ideas/:mentionId", mentionController.DeleteThreadIdea)

	// Tweet Thread Ideas
	router.GET("/api/threads", mentionController.FindAllThreadIdeas)
	router.POST("/api/threads", mentionController.CreateThreadIdea)
	router.DELETE("/api/threads/:mentionId", mentionController.DeleteThreadIdea)

	// Tweet Articles
	router.GET("/api/articles", mentionController.FindAllArticleUrls)
	router.POST("/api/articles", mentionController.CreateArticleUrl)
	router.DELETE("/api/articles/:mentionId", mentionController.DeleteArticleUrl)

	// Tweet Actions
	router.POST("/api/twitter/tweet", twitterController.PostTweet)
	router.GET("/api/twitter/mentions", twitterController.FetchMentions)
	router.POST("/api/twitter/mention/:mendionId", twitterController.ReplyMention)
	router.POST("/api/twitter/direct/:id", twitterController.ReplyDM)

	return router
}