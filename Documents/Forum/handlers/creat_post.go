package forumino

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"unicode"

	forumino "forumino/models"
)

func Creat_PostGlobale(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Path != "/CreatePost" {
		RenderError(w, "404 not found path CreatPost", http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		Creat_PostPage(w, r, forumino.PageData{})
		return
	case http.MethodPost:
		Creat_PostHandler(w, r, db)
		return
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func Creat_PostPage(w http.ResponseWriter, r *http.Request, PageData forumino.PageData) {
	if r.Method != http.MethodGet{
		RenderError(w, "moethod not allowed ", http.StatusMethodNotAllowed)
		return 
	}
	temp, err := template.ParseFiles("template/creat_post.html")
	if err != nil {
		RenderError(w, "template/creat_post.html error Parsing", http.StatusInternalServerError)
		return
	}
	token := GetCSRFToken(r)

	if token == "" {
		token = SetCSRFToken(w)
	}

	PageData.CSRFToken = token

	// var buff bytes.Buffer

	if err = temp.Execute(w, PageData); err != nil {
		RenderError(w, "Template execute error in creat_postPage", http.StatusInternalServerError)
		return
	}
	// if _, err := buff.WriteTo(w); err != nil {
	// 	RenderError(w, "Template execute error in creat_postPage", http.StatusInternalServerError)
	// 	return
	// }
}

func Creat_PostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		RenderError(w, "moethod not allowed ", http.StatusMethodNotAllowed)
		return 
	}
	err := r.ParseForm()
	if err != nil {
		RenderError(w, "500 internal server error , en parsform dans creat_post", http.StatusInternalServerError)
		return
	}

	if !ValidateCSRF(r) {
		RenderError(w, "403 Forbidden: CSRF Token Invalid", http.StatusForbidden)
		return
	}

	Title := r.FormValue("Title")
	Contenue := r.FormValue("Contenue")
	category := r.Form["Category"]

	DataPost := forumino.MY_Post{
		Title:    Title,
		Contenue: Contenue,
		Category: category,
	}
	_, err = validet_post(Title, Contenue, category)
	if err != nil {
		PageData := forumino.PageData{
			ErrorMessege: err,
			Post:         DataPost,
			CSRFToken:    GetCSRFToken(r),
		}
		Creat_PostPage(w, r, PageData)
		return
	}

	err = InsertPostToDB(w, r, db, &DataPost)
	if err != nil {
	if err.Error() == "le Coockie n'exist pas ou user id n'exist pas" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}	
		RenderError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func validet_post(Title string, Contenue string, category []string) (int, error) {
	title := strings.TrimSpace(Title)
	contenue := strings.TrimSpace(Contenue)
	if title == "" {
		return http.StatusBadRequest, errors.New("400 bad request title is empty")
	}
	if contenue == "" {
		return http.StatusBadRequest, errors.New("400 bad request contenue is empty")
	}
	if len(title) > 150 {
		return http.StatusBadRequest, errors.New("400 bad request le nombre maximale des caracter dans  le titre  est 150")
	}
	if len(contenue) > 50000 {
		return http.StatusBadRequest, errors.New("400 bad request le nombre maximale des caracter dans le contenue est 50000")
	}
	if len(category) == 0 {
		return http.StatusBadRequest, errors.New("400 bad request il n y a aucun categorie")
	}
	for _, char := range title {
		if !unicode.IsPrint(char) {
			return http.StatusBadRequest, fmt.Errorf("400 bad request : %s is imprintabl caracter", string(char))
		}
	}
	for _, char := range contenue {
		if !unicode.IsPrint(char) {
			return http.StatusBadRequest, fmt.Errorf("400 bad request : %s is imprintabl caracter", string(char))
		}
	}

	for _, catecategoryName := range category {
		if _, exist := forumino.CategorID.Categoryid[catecategoryName]; !exist {
			return http.StatusBadRequest, errors.New("400 bad request: ce category n'exist pas")
		}
	}
	return 200, nil
}

func InsertPostToDB(w http.ResponseWriter, r *http.Request, db *sql.DB, data *forumino.MY_Post) error {
	UserId, err := SelectUserID(w, r, db)
	if err != nil {
		return errors.New("le Coockie n'exist pas ou user id n'exist pas")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO posts (user_id,title,content) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(UserId, data.Title, data.Contenue)
	if err != nil {
		return err
	}

	PostID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	stmtt, err := tx.Prepare("INSERT OR IGNORE INTO post_categories (post_id,category_id) VALUES(?,?)")
	if err != nil {
		return err
	}

	defer stmtt.Close()

	for _, categoryNAme := range data.Category {
		_, err = stmtt.Exec(int(PostID), forumino.CategorID.Categoryid[categoryNAme])
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func SelectUserID(w http.ResponseWriter, r *http.Request, db *sql.DB) (int, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, err
	}

	TockenValue := cookie.Value
	var UserId int

	err = db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", TockenValue).Scan(&UserId)
	if err != nil {
		return 0, err
	}

	if UserId < 1 {
		return 0, errors.New("invalid user id")
	}

	return UserId, nil
}
