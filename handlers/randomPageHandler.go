package handlers

import (
	"forum/fake"
	"math/rand"
	"net/http"
	"strconv"
)

func RandomPageHandler(w http.ResponseWriter, r *http.Request) {
	postList := fake.GetAllPosts()

	if len(postList) == 0 {
		http.NotFound(w, r)
		return
	}

	randomPost := postList[rand.Intn(len(postList))]

	http.Redirect(w, r, "/posts/"+strconv.Itoa(randomPost.ID), http.StatusSeeOther)
}