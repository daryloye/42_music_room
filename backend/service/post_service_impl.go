package service

import (
	"context"
	"server/data/request"
	"server/data/response"
	"server/helper"
	"server/model"
	"server/repository"
)

type PostServiceImpl struct {
	PostRepository repository.PostRepository
}

func NewPostService(postRepository repository.PostRepository) PostService {
	return &PostServiceImpl{PostRepository: postRepository}
}

func (p *PostServiceImpl) Create(ctx context.Context, request request.PostCreateRequest) {
	postData := model.Post{
		Title:       request.Title,
		Published:   request.Published,
		Description: request.Description,
	}
	p.PostRepository.Save(ctx, postData)
}

func (p *PostServiceImpl) Update(ctx context.Context, request request.PostUpdateRequest) {
	postData := model.Post{
		Id:          request.Id,
		Title:       request.Title,
		Published:   request.Published,
		Description: request.Description,
	}
	p.PostRepository.Update(ctx, postData)
}

func (p *PostServiceImpl) Delete(ctx context.Context, postId string) {
	_, err := p.PostRepository.FindById(ctx, postId)
	helper.ErrorPanic(err)
	p.PostRepository.Delete(ctx, postId)
}

func (p *PostServiceImpl) FindbyId(ctx context.Context, postId string) response.PostResponse {
	post, err := p.PostRepository.FindById(ctx, postId)
	helper.ErrorPanic(err)

	postResponse := response.PostResponse{
		Id:          post.Id,
		Title:       post.Title,
		Published:   post.Published,
		Description: post.Description,
	}

	return postResponse
}

func (p *PostServiceImpl) FindAll(ctx context.Context) []response.PostResponse {
	posts := p.PostRepository.FindAll(ctx)

	var postResp []response.PostResponse

	for _, value := range posts {
		post := response.PostResponse{
			Id:          value.Id,
			Title:       value.Title,
			Published:   value.Published,
			Description: value.Description,
		}
		postResp = append(postResp, post)
	}

	return postResp
}
