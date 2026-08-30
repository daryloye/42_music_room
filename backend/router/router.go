package router

import (
	"net/http"
	"server/controller"
	"server/data/response"
	"server/helper"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(postController *controller.PostController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		webResponse := response.WebResponse{
			Code:   200,
			Status: "Hello",
			Data:   nil,
		}
		helper.WriteResponseBody(w, webResponse)
	})

	router.GET("/api/post", postController.FindAll)
	router.GET("/api/post/:postId", postController.FindById)
	router.POST("/api/post", postController.Create)
	router.PATCH("/api/post/:postId", postController.Update)
	router.DELETE("/api/post/:postId", postController.Delete)

	return router
}
