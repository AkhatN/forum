package model

import (
	"context"
	"net/http"
	"os"
)

//ReadRequestForModerator gets specific post from database for a moderator
func (f *Forum) ReadRequestForModerator(w http.ResponseWriter, r *http.Request) error {
	if err := f.Post.ReadRequestFromDb(); err != nil {
		return nil
	}

	return nil
}

//ReadRequestsForModeratorFromDB gets all requests for a moderator
func (f *Forum) ReadRequestsForModeratorFromDB() error {
	row, err := Db.Query(`
	SELECT posts.id, posts.title, posts.posted_on, coalesce(user.username,'is-null') as to_authname, coalesce(user.role_id, 3) as roleid, posts.valid, posts.valid_admin, posts.isvalid
	FROM posts
    LEFT JOIN posts_request on posts.id = posts_request.post_id 
	LEFT JOIN user ON posts_request.to_auth = user.id
	WHERE posts.isvalid = 1 `)
	if err != nil {
		return err
	}

	for row.Next() {
		post := Post{}
		var valid, admin, isvalid interface{}
		err := row.Scan(&post.ID, &post.Title, &post.Date, &post.AcceptedName, &post.AcceptedAuthRole, &valid, &admin, &isvalid)
		if err != nil {
			return err
		}

		post.Valid = valid.(bool)
		post.ValidAdmin = admin.(bool)
		post.IsValid = isvalid.(bool)

		f.Posts = append(f.Posts, post)
	}

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

//WriteReportToDb writes report for an admin to db
func (f *Forum) WriteReportToDb() error {
	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	INSERT INTO report (from_auth, post_id, body) VALUES (?, ?, ?)`, f.User.ID, f.Post.ID, f.Post.FullText)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `
	UPDATE posts SET valid_admin = 1 WHERE id = ?`, f.Post.ID)
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

//UpdatePostValid approves post
func (f *Forum) UpdatePostValid() error {
	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE posts 
	SET valid = 1, valid_admin = 0, isvalid = 1
	WHERE id = ?`, f.Post.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`INSERT OR REPLACE INTO posts_request (from_auth, to_auth, post_id) 
	VALUES (?,?,?)`, f.Post.AuthID, f.User.ID, f.Post.ID)
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

//CancelModerator demotes a moderator to user
func (f *Forum) CancelModerator() error {

	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE user SET role_id = 3 WHERE id = ?`, f.User.ID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	INSERT or REPLACE INTO make_request (id, from_auth, to_auth, makerequest, cancelrequest, iscancel) 
	VALUES ((SELECT Id FROM make_request WHERE from_auth = ?), ?, (SELECT id FROM user WHERE role_id = 1), ?, ?, ?)`, f.User.ID, f.User.ID, f.User.IsRequest, true, true)
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

//DeletePostAMFromDataBase deletes post from database for moderator or admin
func (p *Post) DeletePostAMFromDataBase() error {
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
	 DELETE FROM image_post WHERE image_post.post_id = $1;
	  `, p.ID)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	DELETE FROM posts WHERE id = $1;`, p.ID)
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
