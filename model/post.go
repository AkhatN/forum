package model

import (
	"context"
	"errors"
	"os"
)

//Post struct
type Post struct {
	ID               int
	Title            string
	FullText         string
	Date             string
	AuthName         string
	AuthID           int
	Cats             []Category
	Comments         []Comment
	Com              Comment
	Sel              SelectLike
	AmountOfLikes    int
	AmountOfDislikes int
	Image            string
	IsImage          bool
	NotificUsername  string
	NotificDate      string
	Valid            bool
	ValidAdmin       bool
	IsValid          bool
	AcceptedName     string
	AcceptedAuthRole int
	Report           string
	Ncats            int
}

//ReadPostByID gets the post from database by Post_ID
func (p *Post) ReadPostByID() error {
	row := Db.QueryRow(`SELECT posts.id, posts.title, posts.body, posts.posted_on, user.username, posts.amount_likes, posts.amount_dislikes 
	FROM posts
	INNER JOIN user ON posts.user_id = user.id
    WHERE posts.id = ? and posts.valid = 1 and posts.isvalid = 1`, p.ID)
	err := row.Scan(&p.ID, &p.Title, &p.FullText, &p.Date, &p.AuthName, &p.AmountOfLikes, &p.AmountOfDislikes)
	if err != nil {
		return err
	}

	if err := row.Err(); err != nil {
		return err
	}

	if err = p.ReadCategories(); err != nil {
		return err
	}

	return nil
}

//WritePosttoDataBase Adding post to Sqlite DataBase, returning error if it couldn't insert ...
func (f *Forum) WritePosttoDataBase() error {
	valid := 0
	if f.User.Role == 2 || f.User.Role == 1 {
		valid = 1
	}
	row := Db.QueryRow(`
	INSERT INTO posts (title, body, posted_on, user_id, amount_likes, amount_dislikes, valid, valid_admin, isvalid) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		f.Post.Title, f.Post.FullText, f.Post.Date, f.Post.AuthID, f.Post.AmountOfLikes, f.Post.AmountOfDislikes, valid, 0, 1)
	err := row.Scan(&f.Post.ID)
	if err != nil {
		return err
	}

	if valid == 1 {
		if err := f.UpdatePostRequest(); err != nil {
			return err
		}
	}

	return nil
}

//ReadCategories gets categories for the post
func (p *Post) ReadCategories() error {
	row, err := Db.Query(`SELECT categories.name 
	FROM cat_posts
	INNER JOIN categories ON categories.id = cat_posts.cat_id WHERE cat_posts.post_id = ?`, p.ID)
	if err != nil {
		return err
	}

	for row.Next() {
		cat := Category{}
		err := row.Scan(&cat.Name)
		if err != nil {
			return err
		}
		p.Cats = append(p.Cats, cat)
	}
	p.Ncats = len(p.Cats)

	return nil
}

//ReadPostForUserFromDB gets specific post for the user
func (p *Post) ReadPostForUserFromDB(ID int) error {
	row := Db.QueryRow(`
	SELECT posts.id, posts.title, posts.body, posts.posted_on, user.id as userid, user.username, amount_likes, amount_dislikes, liked.liked, liked.disliked FROM posts
		INNER JOIN user ON posts.user_id = user.id
		LEFT JOIN liked ON CASE
		WHEN liked.post_id = posts.id AND liked.auth_id = ? THEN 1
		ELSE 0
		END
        WHERE posts.id = ? and posts.valid = 1 and posts.isvalid = 1
`, ID, p.ID)

	var a, b interface{}
	err := row.Scan(&p.ID, &p.Title, &p.FullText, &p.Date, &p.AuthID, &p.AuthName, &p.AmountOfLikes, &p.AmountOfDislikes, &a, &b)
	if err != nil {
		return err
	}

	if err := row.Err(); err != nil {
		return err
	}

	if a == nil && b == nil {
		p.Sel.IsLike = false
		p.Sel.IsDislike = false
	} else {
		p.Sel.IsLike = a.(bool)
		p.Sel.IsDislike = b.(bool)

	}

	if err := row.Err(); err != nil {
		return err
	}

	if err = p.ReadCategories(); err != nil {
		return err
	}

	return nil
}

//DeletePostFromDataBase deletes post from database for registered user
func (p *Post) DeletePostFromDataBase() error {
	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	toDelete := true
	image := ""
	image, err = p.ReadImagefromDb()
	if err != nil {
		toDelete = false
	}

	if toDelete {
		err := os.Remove("." + image)
		if err != nil {
			return err
		}
	}

	_, err = tx.ExecContext(ctx, `
	 DELETE FROM image_post WHERE image_post.post_id = $1 AND
	 (SELECT posts.user_id FROM posts WHERE posts.id = $2) = $3;
	  `, p.ID, p.ID, p.AuthID)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	DELETE FROM posts WHERE id = $1 and user_id = $2;`, p.ID, p.AuthID)
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

//ReadPostForEdit gets specific post for the user
func (p *Post) ReadPostForEdit(UserID int) (int, error) {
	row := Db.QueryRow(`SELECT posts.title, posts.body, user_id FROM posts WHERE posts.id = ? and user_id`, p.ID, UserID)

	userid := 0
	err := row.Scan(&p.Title, &p.FullText, &userid)
	if err != nil {
		return userid, err
	}

	if err := row.Err(); err != nil {
		return userid, err
	}

	if err = p.ReadCategories(); err != nil {
		return userid, err
	}

	return userid, nil
}

//CheckPost ...
func (f *Forum) CheckPost(UserID int) error {
	if f.User.ID != UserID {
		return errors.New("Different UserID")
	}
	return nil
}

//EditPostInDataBase updates post in Sqlite DataBase, returning error if it couldn't update ...
func (f *Forum) EditPostInDataBase() error {
	valid := 0
	if f.User.Role == 2 || f.User.Role == 1 {
		valid = 1
	}
	_, err := Db.Exec(`
	UPDATE posts
	SET title = ?,
	body = ?, 
	valid = ? , 
	valid_admin = ?, 
	isvalid = ?
	WHERE id = ?`, f.Post.Title, f.Post.FullText, valid, 0, 1, f.Post.ID)
	if err != nil {
		return err
	}

	if valid == 1 {
		if err := f.UpdatePostRequest(); err != nil {
			return err
		}
	}
	return nil
}

//ReadAuthorID ...
func (p *Post) ReadAuthorID() error {
	row := Db.QueryRow(`SELECT user.id
	FROM posts
	INNER JOIN user ON posts.user_id = user.id
    WHERE posts.id = ?`, p.ID)
	err := row.Scan(&p.AuthID)
	if err != nil {
		return err
	}

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}
