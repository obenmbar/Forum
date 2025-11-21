package forumino

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	forumino "forumino/models"
)

func ReactHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		RenderError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/Reaction" {
		RenderError(w, "Page Not Found ", http.StatusNotFound)
		return
	}

	// if !ValidateCSRF(r) {
	// 	RenderError(w, "403 Forbidden: CSRF Token Invalid", http.StatusForbidden)
	// 	return
	// }

	ContentType := r.FormValue("type")
	ContentId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		RenderError(w, "Bad Request", http.StatusBadRequest)
		return
	}
	ContentReaction := r.FormValue("reaction")
	ReactionData := forumino.ReactionData{
		ContentType:     ContentType,
		ContentId:       ContentId,
		ContentReaction: ContentReaction,
	}
	err = IsVAlidDAta(ReactionData)
	if err != nil {
		RenderError(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user_id, err := SelectUserID(w, r, db)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		RenderError(w, "500 Error en commencer  la transition", http.StatusInternalServerError)
		return
	}

	switch ReactionData.ContentType {
	case "post":

		err = ExistPostId(ReactionData.ContentId, tx)
		if err != nil {
			RenderError(w, " 500 internalserver error le post n'exist pas", http.StatusInternalServerError)
			return
		}

		NameTable_PostDB := "post_votes"
		err = InsertRraction(user_id, ReactionData.ContentId, NameTable_PostDB, ReactionData.ContentReaction, tx)
		if err != nil {
			RenderError(w, " 500 internalserver error error tanque la modification de reaction post", http.StatusInternalServerError)
			return
		}

	case "comment":

		err = ExistCommentID(ReactionData.ContentId, tx)
		if err != nil {
			RenderError(w, "404 Not Found le comment n'exist pas", http.StatusInternalServerError)
			return
		}

		NameTable_CommentDB := "comment_votes"
		err = InsertRraction(user_id, ReactionData.ContentId, NameTable_CommentDB, ReactionData.ContentReaction, tx)
		if err != nil {
			RenderError(w, "404 Not Found error tanque la modification de Comment", http.StatusNotFound)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		RenderError(w, " 500 internalserver error commit", http.StatusInternalServerError)
		return
	}
}

func IsVAlidDAta(contentdata forumino.ReactionData) error {
	if contentdata.ContentType != "post" && contentdata.ContentType != "comment" {
		return errors.New("400 Bad Request: type invalid")
	}
	if contentdata.ContentId < 1 {
		return errors.New("400 Bad Request: id invalid")
	}
	if contentdata.ContentReaction != "like" && contentdata.ContentReaction != "dislike" {
		return errors.New("400 Bad Request: reaction invalid")
	}
	return nil
}

func InsertRraction(user_id, Id int, reaction string, NameTAbleCoorPO string, tx *sql.Tx) error {
	switch reaction {
	case "like":
		err := ModifierReaction(user_id, Id, true, NameTAbleCoorPO, tx)
		if err != nil {
			return err
		}
	case "dislike":
		err := ModifierReaction(user_id, Id, false, NameTAbleCoorPO, tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExistPostId(post_id int, tx *sql.Tx) error {
	var exists int
	err := tx.QueryRow("SELECT 1 FROM posts WHERE id = ?", post_id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("404 Not Found: post n'existe pas")
		}

		return err
	}
	return nil
}

func ModifierReaction(user_id, id int, NewVote bool, NameTableVote string, tx *sql.Tx) error {
	var CurrentVote bool
	var Table_idpoorcom string

	if NameTableVote == "post_votes" {
		Table_idpoorcom = "post_id"
	} else {
		Table_idpoorcom = "comment_id"
	}

	query := fmt.Sprintf("SELECT vote_type FROM %s WHERE user_id = ? AND %s = ?", NameTableVote, Table_idpoorcom)

	err := tx.QueryRow(query, user_id, id).Scan(&CurrentVote)
	if err != nil {
		if err == sql.ErrNoRows {
			query := fmt.Sprintf("INSERT INTO %s (user_id,%s,vote_type)VALUES(?,?,?)", NameTableVote, Table_idpoorcom)
			_, err = tx.Exec(query, user_id, id, NewVote)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	if CurrentVote == NewVote {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id = ? AND %s = ?", NameTableVote, Table_idpoorcom)
		_, err := tx.Exec(query, user_id, id)
		if err != nil {
			return err
		}

	} else {
		query := fmt.Sprintf("UPDATE %s SET vote_type = ? WHERE user_id = ? AND %s = ?", NameTableVote, Table_idpoorcom)
		_, err := tx.Exec(query, NewVote, user_id, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExistCommentID(comment_id int, tx *sql.Tx) error {
	var Exist int
	err := tx.QueryRow("SELECT 1 FROM comments WHERE id= ?", comment_id).Scan(&Exist)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Error : Invalid Comment_Id")
		}
		return err
	}

	return nil
}
