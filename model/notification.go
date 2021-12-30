package model

//ReadLikedNotification reads likes and dislikes notifications from database
func (f *Forum) ReadLikedNotification() error {
	row, err := Db.Query(`
	SELECT posts.id, posts.title, user.username, notific_liked.notific_date, liked.liked, liked.disliked FROM posts 
INNER JOIN notific_liked ON notific_liked.post_id = posts.id
INNER JOIN user ON user.id = notific_liked.from_auth
INNER JOIN liked ON posts.id = liked.post_id AND liked.auth_id = notific_liked.from_auth 
WHERE posts.user_id = ? AND notific_liked.from_auth != ?`,
		f.User.ID, f.User.ID)
	if err != nil {
		return err
	}

	for row.Next() {
		post := Post{}
		var a, b interface{}
		err := row.Scan(&post.ID, &post.Title, &post.NotificUsername, &post.NotificDate, &a, &b)
		if err != nil {
			return err
		}
		if a == nil && b == nil {
			post.Sel.IsLike = false
			post.Sel.IsDislike = false
		} else {
			post.Sel.IsLike = a.(bool)
			post.Sel.IsDislike = b.(bool)

		}

		f.Posts = append(f.Posts, post)
	}
	if err := row.Err(); err != nil {
		return err

	}

	return nil
}

//ReadCommentNotification reads comments notifications from database
func (f *Forum) ReadCommentNotification() error {
	row, err := Db.Query(`
	SELECT DISTINCT posts.id, posts.title, user.username, notific_comment.notific_date from posts 
INNER JOIN notific_comment ON notific_comment.post_id = posts.id
INNER JOIN user ON user.id = notific_comment.from_auth
INNER JOIN comments ON posts.id = comments.post_id AND comments.auth_id = notific_comment.from_auth 
WHERE posts.user_id = ? AND notific_comment.from_auth != ?
	`, f.User.ID, f.User.ID)
	if err != nil {
		return err
	}

	for row.Next() {
		c := Comment{}
		err := row.Scan(&c.PostID, &c.PostName, &c.NotificUsername, &c.NotificDate)
		if err != nil {
			return err
		}

		f.Comments = append(f.Comments, c)
	}
	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

//UpdateNotifications updates new notifications to old
func (f *Forum) UpdateNotifications() error {
	_, err := Db.Exec(`
	UPDATE notific_liked SET valid = 0 WHERE to_auth = ?;
	UPDATE notific_comment SET valid = 0 WHERE to_auth = ?`, f.User.ID, f.User.ID)
	if err != nil {
		return err
	}

	return nil
}

//ReadCountNotification reads number of new notifications
func (f *Forum) ReadCountNotification() error {
	row := Db.QueryRow(`SELECT count(notific_liked.valid) FROM notific_liked 
	INNER JOIN posts ON notific_liked.post_id = posts.id 
	WHERE notific_liked.valid = 1 AND posts.user_id = ? AND notific_liked.from_auth != ?
	`, f.User.ID, f.User.ID)
	err := row.Scan(&f.AmountNotification)
	if err != nil {
		return err
	}

	var amountComment int
	row2 := Db.QueryRow(`SELECT count(notific_comment.valid) FROM notific_comment 
	INNER JOIN posts ON notific_comment.post_id = posts.id 
	WHERE notific_comment.valid = 1 AND posts.user_id = ? AND notific_comment.from_auth != ?
	`, f.User.ID, f.User.ID)
	err = row2.Scan(&amountComment)
	if err != nil {
		return err
	}

	f.AmountNotification += amountComment

	return nil
}
