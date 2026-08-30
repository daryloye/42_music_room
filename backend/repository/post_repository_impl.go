package repository

import (
	"context"
	"errors"
	"fmt"
	"server/helper"
	"server/model"
	"server/prisma/db"
)

type PostRepositoryImpl struct {
	Db *db.PrismaClient
}

// NewPostRepository connects the DB client to the methods defined in the PostRepository interface
func NewPostRepository(Db *db.PrismaClient) PostRepository {
	return &PostRepositoryImpl{Db: Db}
}

func (p *PostRepositoryImpl) Save(ctx context.Context, post model.Post) {
	result, err := p.Db.Post.
		CreateOne(
			db.Post.Title.Set(post.Title),
			db.Post.Published.Set(post.Published),
			db.Post.Description.Set(post.Description),
		).
		Exec(ctx)
	helper.ErrorPanic(err)
	fmt.Println("Rows affected: ", result)
}

func (p *PostRepositoryImpl) Update(ctx context.Context, post model.Post) {
	result, err := p.Db.Post.
		FindUnique(db.Post.ID.Equals(post.Id)).
		Update(
			db.Post.Title.Set(post.Title),
			db.Post.Published.Set(post.Published),
			db.Post.Description.Set(post.Description),
		).
		Exec(ctx)
	helper.ErrorPanic(err)
	fmt.Println("Rows affected: ", result)
}

func (p *PostRepositoryImpl) Delete(ctx context.Context, postId string) {
	result, err := p.Db.Post.
		FindUnique(db.Post.ID.Equals(postId)).
		Delete().
		Exec(ctx)
	helper.ErrorPanic(err)
	fmt.Println("Rows affected: ", result)
}

func (p *PostRepositoryImpl) FindById(ctx context.Context, postId string) (model.Post, error) {
	post, err := p.Db.Post.
		FindUnique(db.Post.ID.Equals(postId)).
		Exec(ctx)
	helper.ErrorPanic(err)

	published, _ := post.Published()     // because published is Boolean? nullable
	description, _ := post.Description() // because description is String? nullable

	postData := model.Post{
		Id:          post.ID,
		Title:       post.Title,
		Published:   published,
		Description: description,
	}

	if post != nil {
		return postData, nil
	} else {
		return postData, errors.New("Post id not found")
	}
}

func (p *PostRepositoryImpl) FindAll(ctx context.Context) []model.Post {
	allPosts, err := p.Db.Post.
		FindMany().
		Exec(ctx)
	helper.ErrorPanic(err)

	var posts []model.Post

	for _, post := range allPosts {
		published, _ := post.Published()
		description, _ := post.Description()

		postData := model.Post{
			Id:          post.ID,
			Title:       post.Title,
			Published:   published,
			Description: description,
		}
		posts = append(posts, postData)
	}

	return posts
}
