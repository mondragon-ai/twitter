package router

import (
	"fmt"

	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/twitter/controller"
)

func MentionsRouter(bookController *controller.MentionsController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprint(w, "Welcome Home")
	})

	router.GET("/api/mentions", bookController.FindAll)
	router.GET("/api/mentions/:mentionId", bookController.FindById)
	router.POST("/api/mentions", bookController.Create)

	router.DELETE("/api/mentions/:mentionId", bookController.Delete)

	return router
}