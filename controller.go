package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type memoData struct {
	Memo
	Title string
	ID    string
}

type indexData struct {
	Title  string
	Filter string
	Memos  []Memo
}

func memoGET(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("id")
	if ids == "" || !bson.IsObjectIdHex(ids) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id := bson.ObjectIdHex(ids)

	// Get memo
	memo, err := getMemoByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Load templates
	tmpl, _ := template.ParseFiles(
		"view/memo.html",
		"view/head.html",
	)

	tmpl.ExecuteTemplate(w, "memo", memoData{
		Memo:  memo,
		Title: "Memo • " + memo.Name,
		ID:    id.Hex(),
	})
}

func memoPOST(w http.ResponseWriter, r *http.Request) {
	id := bson.NewObjectId()
	ids := r.URL.Query().Get("id")
	if ids != "" && bson.IsObjectIdHex(ids) {
		id = bson.ObjectIdHex(ids)
	}

	r.ParseMultipartForm(32 << 20)

	var memo Memo
	memo.ID = id
	memo.Name = r.FormValue("name")
	memo.Category = r.FormValue("category")
	memo.Content = template.HTML(r.FormValue("content"))
	memo.CreatedAt = time.Now()

	updateMemo(memo)

	// Delete memo picture
	if r.FormValue("delete-avatar") == "on" {
		os.Remove("data/" + id.Hex() + ".png")
	}

	// Memo picture
	file, _, err := r.FormFile("avatar")
	if err == nil {
		f, _ := os.OpenFile("data/"+id.Hex()+".png", os.O_WRONLY|os.O_CREATE, 0666)
		io.Copy(f, file)

		f.Close()
		file.Close()
	}

	http.Redirect(w, r, "/memo?id="+memo.ID.Hex(), http.StatusSeeOther)
}

func editMemoGET(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("id")
	if ids == "" || !bson.IsObjectIdHex(ids) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id := bson.ObjectIdHex(ids)

	// Get memo
	memo, err := getMemoByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Load templates
	tmpl, _ := template.ParseFiles(
		"view/memo_edit.html",
		"view/head.html",
	)

	tmpl.ExecuteTemplate(w, "memo", memoData{
		Memo:  memo,
		Title: "Title • " + memo.Name + " (edition)",
		ID:    id.Hex(),
	})
}

func newMemoGET(w http.ResponseWriter, r *http.Request) {
	// Load templates
	tmpl, _ := template.ParseFiles(
		"view/memo_edit.html",
		"view/head.html",
	)

	tmpl.ExecuteTemplate(w, "memo", memoData{
		Title: "Memo • Nouveau",
	})
}

func deleteMemoPOST(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("id")
	if ids == "" || !bson.IsObjectIdHex(ids) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id := bson.ObjectIdHex(ids)

	if err := deleteMemoByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func indexGET(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")

	// Load templates
	tmpl, _ := template.ParseFiles(
		"view/home.html",
		"view/head.html",
	)

	// Get memos
	memos, _ := getAllMemos(filter)

	tmpl.ExecuteTemplate(w, "home", indexData{
		Title:  "Memo",
		Filter: filter,
		Memos:  memos,
	})
}
