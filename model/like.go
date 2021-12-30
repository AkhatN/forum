package model

import (
	"context"
	"time"
)

//SelectLike struct of post like
type SelectLike struct {
	IsLike    bool
	IsDislike bool
	PostID    int
	AuthID    int
	ToAuth    int
}

//SelectComLike struct of comment like
type SelectComLike struct {
	IsLike    bool
	IsDislike bool
	ComID     int
	AuthID    int
}

//WriteLikeToDataBase writes like of a post to database
func (s *SelectLike) WriteLikeToDataBase() error {
	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	INSERT OR REPLACE INTO liked 
	(id, liked, disliked, post_id, auth_id) VALUES(
		(SELECT Id FROM liked WHERE post_id = ? AND auth_id = ?), ?, ?, ?, ?)`,
		s.PostID, s.AuthID, s.IsLike, s.IsDislike, s.PostID, s.AuthID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `
	UPDATE posts SET amount_likes = (SELECT count(liked) FROM liked WHERE liked = 1 AND post_id = ?), 
	amount_dislikes = (SELECT count(disliked) FROM liked WHERE disliked = 1 AND post_id = ?) WHERE id = ?
	`, s.PostID, s.PostID, s.PostID)
	if err != nil {
		tx.Rollback()
		return err
	}

	//checks for valid notification
	valid := 1
	if !s.IsLike && !s.IsDislike {
		valid = 0
	}
	_, err = tx.ExecContext(ctx, `INSERT OR REPLACE INTO notific_liked (id, valid, notific_date, post_id, from_auth, to_auth) 
	VALUES ((SELECT id FROM notific_liked WHERE post_id = ? AND from_auth = ?), ?, ?, ?, ?, ?)`, s.PostID, s.AuthID, valid, time.Now().Format("2006.01.02 15:04:05"), s.PostID, s.AuthID, s.ToAuth)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//WriteComLikeToDataBase writes like of a comment like to database
func (s *SelectComLike) WriteComLikeToDataBase() error {
	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	INSERT OR REPLACE INTO comlikes 
	(id, liked, disliked, auth_id, com_id) values(
		(SELECT Id FROM comlikes WHERE com_id = ? AND auth_id = ?), ?, ?, ?, ?)`,
		s.ComID, s.AuthID, s.IsLike, s.IsDislike, s.AuthID, s.ComID)
	if err != nil {
		tx.Rollback()
		return err
	}
	var amountOfLikes, amountOfDislikes int
	row := tx.QueryRow(`
	SELECT count(liked) AS count_liked, (SELECT count(disliked) FROM comlikes WHERE disliked = 1 AND com_id = ?) FROM comlikes WHERE liked = 1 AND com_id = ?
	`, s.ComID, s.ComID)
	err = row.Scan(&amountOfLikes, &amountOfDislikes)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE comments SET amount_likes = ?, amount_dislikes = ? WHERE id = ?", amountOfLikes, amountOfDislikes, s.ComID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

//ReadPostAuthIDByPostID ...
func (s *SelectLike) ReadPostAuthIDByPostID() error {
	row := Db.QueryRow("SELECT user_id FROM posts WHERE id = ?", s.PostID)

	err := row.Scan(&s.ToAuth)
	if err != nil {
		return err
	}

	return nil
}
