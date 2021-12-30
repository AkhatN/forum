package model

import (
	"context"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//WriteAdmin writes new admin user into database
func WriteAdmin() error {
	var err error
	user := User{
		Username: "admin",
		Email:    "admin@mail.ru",
		Password: "admin",
	}
	user.HashPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = Db.Exec(`
	INSERT OR IGNORE INTO user (username, email, password_hash, role_id) 
	VALUES (?, ?, ?, 1) `, user.Username, user.Email, string(user.HashPassword))
	if err != nil {
		return err
	}

	return nil
}

//ReadRequestForAdmin gets specific post from database for an admin
func (f *Forum) ReadRequestForAdmin(w http.ResponseWriter, r *http.Request) error {
	if err := f.Post.ReadRequestFromDb(); err != nil {
		return nil
	}

	if err := f.Post.ReadReportForAdminFromDb(); err != nil {
		return nil
	}

	return nil
}

//ReadRequestsForAdminFromDb gets all requests from database for an admin
func (f *Forum) ReadRequestsForAdminFromDb() error {
	// Reading request posts
	row, err := Db.Query(`
	SELECT posts.id, posts.title, posts.posted_on, coalesce(user.username,'is-null') as to_authname, coalesce(user.role_id, 3) as roleid, posts.valid, posts.valid_admin, posts.isvalid
	FROM posts
    LEFT JOIN posts_request ON posts.id = posts_request.post_id 
	LEFT JOIN user ON posts_request.to_auth = user.id
	WHERE posts.isvalid = 1 and posts.valid_admin = 0`)
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

	row2, err := Db.Query(`
	SELECT posts.id, posts.title, posts.posted_on, user.username as to_authname, user.role_id, posts.valid, posts.valid_admin, posts.isvalid
	FROM posts
    inner JOIN report on report.post_id = posts.id 
	inner JOIN user ON report.from_auth = user.id
	WHERE posts.isvalid = 1 and posts.valid_admin = 1`)
	if err != nil {
		return err
	}

	for row2.Next() {
		post := Post{}
		var valid, admin, isvalid interface{}
		err := row2.Scan(&post.ID, &post.Title, &post.Date, &post.AcceptedName, &post.AcceptedAuthRole, &valid, &admin, &isvalid)
		if err != nil {
			return err
		}

		post.Valid = valid.(bool)
		post.ValidAdmin = admin.(bool)
		post.IsValid = isvalid.(bool)

		f.Posts = append(f.Posts, post)
	}

	if err := row2.Err(); err != nil {
		return err

	}

	// Reading request from user
	row, err = Db.Query(`
	SELECT user.id, user.username, make_request.makerequest FROM user 
    INNER JOIN make_request ON make_request.from_auth = user.id 
    WHERE make_request.makerequest = 1`)
	if err != nil {
		return err
	}

	for row.Next() {
		u := User{}
		err := row.Scan(&u.ID, &u.Username, &u.IsRequest)
		if err != nil {
			return err
		}

		f.Users = append(f.Users, u)
	}

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

//ReadReportForAdminFromDb gets report about specific request for an admin
func (p *Post) ReadReportForAdminFromDb() error {
	row := Db.QueryRow(`
	SELECT report.body, user.username 
	FROM report
	INNER JOIN user ON user.id = report.from_auth
	WHERE report.post_id = ?`, p.ID)

	err := row.Scan(&p.Report, &p.AcceptedName)
	if err != nil {
		return err
	}

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

//UpdateRoleAdmin updates role for a user
func (u *User) UpdateRoleAdmin() error {
	if err := u.UpdateRoleUser(); err != nil {
		return err
	}

	if err := u.UpdateMakeRequestAdmin(); err != nil {
		return err
	}

	return nil
}

//UpdateRoleUser sets new role a user
func (u *User) UpdateRoleUser() error {
	_, err := Db.Exec(`UPDATE user SET role_id = ? WHERE id = ?`, u.NewRole, u.ID)
	if err != nil {
		return err
	}

	return nil
}

//UpdateMakeRequestAdmin updates make_request table
func (u *User) UpdateMakeRequestAdmin() error {
	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	row := tx.QueryRow(`SELECT makerequest, cancelrequest, iscancel FROM make_request WHERE from_auth = ?`, u.ID)
	beforemake := false
	beforecancel := false
	beforeiscancel := false
	iscancer := false
	iscancel := false
	err = row.Scan(&beforemake, &beforecancel, &beforeiscancel)
	if err != nil {
	}

	if u.NewRole == 3 && beforemake && !beforeiscancel {
		iscancer = true
	}

	if err := row.Err(); err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
	INSERT OR REPLACE INTO make_request (from_auth, to_auth, makerequest, cancelrequest, iscancel) 
	VALUES (?, (SELECT id FROM user WHERE role_id = 1), ?, ?, ?)`, u.ID, false, iscancer, iscancel)
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

//DeleteCommentAdminFromDataBase admin deletes comments from database
func (c *Comment) DeleteCommentAdminFromDataBase() error {
	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM comments WHERE comments.id = $1;`, c.ID)
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
