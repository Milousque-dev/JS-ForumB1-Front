package handlers

import (
	"forum/database"
	"forum/fake"
	"net/http"
	"strconv"
)

func CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	toggleCommentReaction(w, r, 1)
}

func CommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
	toggleCommentReaction(w, r, -1)
}

func toggleCommentReaction(w http.ResponseWriter, r *http.Request, value int) {
	user, isLogged := fake.GetCurrentUserFull(r)
	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentIDStr := r.PathValue("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	comment, found := fake.GetCommentByID(commentID)
	if !found {
		http.NotFound(w, r)
		return
	}

	database.ToggleCommentLike(user.ID, commentID, value)

	http.Redirect(w, r, "/posts/"+strconv.Itoa(comment.PostID), http.StatusSeeOther)
}
