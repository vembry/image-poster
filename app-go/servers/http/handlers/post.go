package handlers

import (
	"app-go/server/http/middlewares"
	"net/http"
)

type post struct {
}

func NewPost() *post {
	return &post{}
}

func (p *post) GetRoutes() *http.ServeMux {
	postmux := http.NewServeMux()

	postmux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		respondJson(w, http.StatusOK, map[string]string{
			"message": "ok",
		})
	})

	postmux.HandleFunc("GET /list", p.ListPost)
	postmux.HandleFunc("POST /", p.Post)
	postmux.HandleFunc("POST /{postId}/comment", p.PostComment)
	postmux.HandleFunc("DELETE /{postId}/comment/{commentId}", p.DeleteComment)

	// group endpoints with 1 prefix
	group := http.NewServeMux()
	group.Handle("/post/", middlewares.Auth(http.StripPrefix("/post", postmux)))

	return group
}

// ListPost handle http request to get a paginated list of posts
func (p *post) ListPost(w http.ResponseWriter, r *http.Request) {
	respondJson(w, http.StatusOK, r.URL.Query())
}

// Post handle http request to create post entry
func (p *post) Post(w http.ResponseWriter, r *http.Request) {
	respondJson(w, http.StatusOK, nil)
}

// PostComment handle http request to post a comment on a post
func (p *post) PostComment(w http.ResponseWriter, r *http.Request) {
	respondJson(w, http.StatusOK, map[string]string{
		"postId": r.PathValue("postId"),
	})
}

// DeleteComment handle http request to delete a comment from a post.
func (p *post) DeleteComment(w http.ResponseWriter, r *http.Request) {
	respondJson(w, http.StatusOK, map[string]string{
		"postId":    r.PathValue("postId"),
		"commentId": r.PathValue("commentId"),
	})
}
