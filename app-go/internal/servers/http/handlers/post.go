package handlers

import (
	postmodule "app-go/internal/modules/post"
	"app-go/internal/modules/post/models"
	"app-go/internal/servers/http/middlewares"
	"log"
	"net/http"
)

type post struct {
	postProvider postmodule.IPost
}

func NewPost(postProvider postmodule.IPost) *post {
	return &post{
		postProvider: postProvider,
	}
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
	group.Handle("/post/", middlewares.Auth(http.StripPrefix("/post", postmux))) // cover the entrypoint with middleware

	return group
}

// ListPost handle http request to get a paginated list of posts
func (p *post) ListPost(w http.ResponseWriter, r *http.Request) {
	respondJson(w, http.StatusOK, r.URL.Query())
}

// Post handle http request to create post entry
func (p *post) Post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// set size limit of incoming request
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20)

	// parse incoming request
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		log.Printf("error on parsing request into multipart/form. err=%v", err)
		respondJson(w, http.StatusBadRequest, map[string]string{
			"error": "error on parsing request",
		})
		return
	}

	// attempt to read file from request
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("error on getting file. err=%v", err)
		respondJson(w, http.StatusBadRequest, map[string]string{
			"error": "error on reading file on request",
		})
		return
	}
	if file != nil {
		defer file.Close()
	}

	// read non-file payload from request
	text := r.FormValue("text")

	// call service
	err = p.postProvider.CreatePost(r.Context(), models.CreatePostArg{
		Text: text,
		File: file,
	})
	if err != nil {
		log.Printf("error on submitting post. err=%v", err)
		respondJson(w, http.StatusBadRequest, map[string]string{
			"error": "error processing post submission",
		})
		return
	}

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
