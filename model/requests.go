package model

import (
	"context"
	"os"
)

//DeleteRequestFromDataBase deletes request and post from database
func (p *Post) DeleteRequestFromDataBase(userid int) error {
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
	 DELETE FROM image_post WHERE image_post.post_id = $1;`, p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `
	 DELETE FROM cat_posts WHERE post_id = $1;`, p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `
	UPDATE posts SET valid = 0, valid_admin = 0, isvalid = 0 WHERE id = $1;`, p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
	INSERT OR REPLACE INTO posts_request (from_auth, to_auth, post_id) 
	VALUES (?,?,?)`, p.AuthID, userid, p.ID)
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

//ReadRequestFromDb gets specific request for a moderator
func (p *Post) ReadRequestFromDb() error {
	row := Db.QueryRow(`
	SELECT posts.id, posts.title, posts.body, posts.posted_on, user.id as userid, user.username, posts.valid FROM posts
		INNER JOIN user ON posts.user_id = user.id
        WHERE posts.id = ? AND posts.isvalid = 1
`, p.ID)

	var valid interface{}
	err := row.Scan(&p.ID, &p.Title, &p.FullText, &p.Date, &p.AuthID, &p.AuthName, &valid)
	if err != nil {
		return err
	}

	p.Valid = valid.(bool)

	if err := row.Err(); err != nil {
		return err
	}

	if err = p.ReadCategories(); err != nil {
		return err
	}

	return nil
}

//UpdatePostRequest ...
func (f *Forum) UpdatePostRequest() error {
	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT OR REPLACE INTO posts_request (from_auth, to_auth, post_id) 
	VALUES (?,?,?)`, f.User.ID, f.User.ID, f.Post.ID)
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
