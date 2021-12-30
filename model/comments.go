package model

import (
	"context"
	"errors"
)

//Comment struct
type Comment struct {
	ID               int
	Comm             string
	AuthID           int
	AuthName         string
	PostName         string
	PostID           int
	Date             string
	Sel              SelectComLike
	AmountOfLikes    int
	AmountOfDislikes int
	ToAuth           int
	NotificDate      string
	NotificUsername  string
}

//ReadCommensByID gets comments for the post from database Post_ID
func (p *Post) ReadCommensByID() error {
	row, err := Db.Query(`SELECT comments.id, comments.comment, comments.auth_id, user.username, comments.post_id, comments.amount_likes, comments.amount_dislikes, comments.commented_on 
	FROM comments 
	INNER JOIN user on user.id = comments.auth_id 
	WHERE post_id = ?
	`, p.ID)
	if err != nil {
		return err
	}

	for row.Next() {
		comm := Comment{}
		err := row.Scan(&comm.ID, &comm.Comm, &comm.AuthID, &comm.AuthName, &comm.PostID, &comm.AmountOfLikes, &comm.AmountOfDislikes, &comm.Date)
		if err != nil {
			return err
		}
		p.Comments = append(p.Comments, comm)
	}
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}

//ReadComLikesForUser gets comments for specific post for the user
func (p *Post) ReadComLikesForUser(ID int) error {
	row, err := Db.Query(`SELECT comments.id, comments.comment,  comments.auth_id, user.username, comments.post_id, comments.amount_likes, comments.amount_dislikes, comlikes.liked, comlikes.disliked, comments.commented_on 
	FROM comments
			INNER JOIN user ON comments.auth_id = user.id
			LEFT JOIN comlikes ON CASE
			WHEN comlikes.com_id = comments.id AND comlikes.auth_id = ? THEN 1
			ELSE 0
			END
			WHERE comments.post_id = ?
	`, ID, p.ID)
	if err != nil {
		return err
	}

	for row.Next() {
		comm := Comment{}
		var a, b interface{}
		err := row.Scan(&comm.ID, &comm.Comm, &comm.AuthID, &comm.AuthName, &comm.PostID, &comm.AmountOfLikes, &comm.AmountOfDislikes, &a, &b, &comm.Date)
		if err != nil {
			return err
		}

		if a == nil && b == nil {
			comm.Sel.IsLike = false
			comm.Sel.IsDislike = false
		} else {
			comm.Sel.IsLike = a.(bool)
			comm.Sel.IsDislike = b.(bool)

		}
		p.Comments = append(p.Comments, comm)
	}

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

//ReadPostAuthIDByPostID gets id of author post
func (c *Comment) ReadPostAuthIDByPostID() error {
	row := Db.QueryRow("SELECT user_id FROM posts WHERE id = ?", c.PostID)

	err := row.Scan(&c.ToAuth)
	if err != nil {
		return err
	}

	return nil
}

//WriteCommentToDataBase writes comments into database
func (c *Comment) WriteCommentToDataBase() error {
	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	row := tx.QueryRow("INSERT INTO comments (comment, auth_id, post_id, commented_on) VALUES (?, ?, ?,?) RETURNING id", c.Comm, c.AuthID, c.PostID, c.Date)

	err = row.Scan(&c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO notific_comment (valid, notific_date, post_id, com_id, from_auth, to_auth) 
	VALUES (?, ?, ?, ?, ?, ?)`, 1, c.Date, c.PostID, c.ID, c.AuthID, c.ToAuth)
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

//DeleteCommentFromDataBase deletes comments from database
func (c *Comment) DeleteCommentFromDataBase(userid int) error {
	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `	
 DELETE FROM comments WHERE comments.id = $1 and comments.auth_id = $2;
	`, c.ID, c.AuthID)
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

//ReadCommentByID gets comment for the post from database Post_ID
func (f *Forum) ReadCommentByID() error {
	row := Db.QueryRow(`SELECT comments.id, comments.comment, comments.auth_id, user.username, comments.post_id
	FROM comments 
	INNER JOIN user on user.id = comments.auth_id 
	WHERE comments.id = ? and user.id = ?
	`, f.Post.Com.ID, f.User.ID)

	err := row.Scan(&f.Post.Com.ID, &f.Post.Com.Comm, &f.Post.Com.AuthID, &f.Post.Com.AuthName, &f.Post.ID)
	if err != nil {
		return err
	}

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

//CheckComment ...
func (f *Forum) CheckComment(PostID int) error {
	if f.Post.ID != PostID {
		return errors.New("Different PostID")
	}
	return nil
}

//EditCommentInDataBase edites comments in database
func (c *Comment) EditCommentInDataBase() error {
	_, err := Db.Exec(`
	UPDATE comments
	SET comment = ?
	WHERE id = ? and auth_id = ?`, c.Comm, c.ID, c.AuthID)
	if err != nil {
		return err
	}

	return nil
}
