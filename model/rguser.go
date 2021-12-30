package model

import (
	"context"
	"net/http"
)

//ReadRequestsForUserFromDB gets all requests for a user from db
func (f *Forum) ReadRequestsForUserFromDB() error {
	// Reading requests of posts ...
	row, err := Db.Query(`
	SELECT posts.id, posts.title, posts.valid, posts.valid_admin, posts.isvalid, coalesce(user.username,'is-null') as to_authname, coalesce(user.role_id, 3) as roleid, posts.posted_on
	FROM posts
    LEFT join posts_request on posts.id = posts_request.post_id 
	LEFT JOIN user ON posts_request.to_auth = user.id
  WHERE posts.user_id = ?`, f.User.ID)
	if err != nil {
		return err
	}

	for row.Next() {
		post := Post{}
		var valid, admin, isvalid interface{}
		err := row.Scan(&post.ID, &post.Title, &valid, &admin, &isvalid, &post.AcceptedName, &post.AcceptedAuthRole, &post.Date)
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

//WriteRequestRegisteredUserToAdmin ...
func (f *Forum) WriteRequestRegisteredUserToAdmin(w http.ResponseWriter, r *http.Request) error {

	switch r.PostFormValue("makereq") {
	case "submit":
		f.User.IsRequest = true
		f.User.IsCancel = false
	case "cancel":
		f.User.IsRequest = false
		f.User.IsCancel = true
	default:
		f.User.IsRequest = false
		f.User.IsCancel = true
	}

	// Request want to be a moderator
	if err := f.WriteRequestToAdmin(); err != nil {
		return err
	}

	return nil
}

//WriteRequestToAdmin ...
func (f *Forum) WriteRequestToAdmin() error {

	ctx := context.Background()
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	cancel := false
	if f.User.IsRequest == false {
		cancel = true
	}

	_, err = tx.ExecContext(ctx, `
	INSERT or REPLACE INTO make_request (id, from_auth, to_auth, makerequest, cancelrequest, iscancel) 
	VALUES ((SELECT Id FROM make_request WHERE from_auth = ?), ?, (SELECT id FROM user WHERE role_id = 1), ?, ?, ?)`, f.User.ID, f.User.ID, f.User.IsRequest, cancel, f.User.IsCancel)
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
