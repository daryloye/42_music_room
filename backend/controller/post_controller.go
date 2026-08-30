package controller

import (
	"net/http"
	"server/data/request"
	"server/data/response"
	"server/helper"
	"server/service"

	"github.com/julienschmidt/httprouter"
)

type PostController struct {
	PostService service.PostService
}

func NewPostController(postService service.PostService) *PostController {
	return &PostController{PostService: postService}
}

func (controller *PostController) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	postCreateRequest := request.PostCreateRequest{}
	helper.ReadRequestBody(r, &postCreateRequest)

	controller.PostService.Create(r.Context(), postCreateRequest)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}
	helper.WriteResponseBody(w, webResponse)
}

func (controller *PostController) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	postUpdateRequest := request.PostUpdateRequest{}
	helper.ReadRequestBody(r, &postUpdateRequest)

	postId := params.ByName("postId")
	postUpdateRequest.Id = postId

	controller.PostService.Update(r.Context(), postUpdateRequest)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}
	helper.WriteResponseBody(w, webResponse)
}

func (controller *PostController) FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result := controller.PostService.FindAll(r.Context())
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   result,
	}
	helper.WriteResponseBody(w, webResponse)
}

func (controller *PostController) FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	postId := params.ByName("postId")

	result := controller.PostService.FindbyId(r.Context(), postId)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   result,
	}
	helper.WriteResponseBody(w, webResponse)
}

func (controller *PostController) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	postId := params.ByName("postId")

	controller.PostService.Delete(r.Context(), postId)
	webResponse := response.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}
	helper.WriteResponseBody(w, webResponse)
}
