package model

import (
	"fmt"
)

//ReadPostsByCatName reads all posts dedicated to a specific topic to a guest
func (f *Forum) ReadPostsByCatName() error {
	rows, err := Db.Query(`
		SELECT posts.id, posts.title, posts.body, posts.posted_on, user.username, amount_likes, amount_dislikes 
		FROM categories 
		INNER JOIN cat_posts ON cat_posts.cat_id = (select categories.id WHERE categories.name = ?)
		INNER JOIN posts ON posts.id = cat_posts.post_id
		INNER JOIN user ON user.id = posts.user_id
		WHERE posts.valid = 1 and posts.isvalid = 1
		`, f.Post.Title)
	if err != nil {
		return err
	}

	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.FullText, &post.Date, &post.AuthName, &post.AmountOfLikes, &post.AmountOfDislikes)
		if err != nil {
			return err
		}
		f.Posts = append(f.Posts, post)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for i := range f.Posts {
		if err = f.Posts[i].ReadCategories(); err != nil {
			return err
		}
	}
	return nil
}

//ReadCatPostForRegisteredUserFromDB reads all posts dedicated to a specific topic to a registered user
func (f *Forum) ReadCatPostForRegisteredUserFromDB() error {
	row, err := Db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.posted_on, user.id as userid,user.username, amount_likes, amount_dislikes, liked.liked, liked.disliked
	FROM categories 
	INNER JOIN cat_posts ON cat_posts.cat_id = (SELECT categories.id WHERE categories.name = ?)
	INNER JOIN posts ON posts.id = cat_posts.post_id
	INNER JOIN user ON user.id = posts.user_id
	LEFT JOIN liked ON CASE
	WHEN liked.post_id = posts.id AND liked.auth_id = ? THEN 1
	ELSE 0
	END
	WHERE posts.valid = 1`, f.Post.Title, f.User.ID)
	if err != nil {
		return err
	}

	for row.Next() {
		post := Post{}
		var a, b interface{}
		err := row.Scan(&post.ID, &post.Title, &post.FullText, &post.Date, &post.AuthID, &post.AuthName, &post.AmountOfLikes, &post.AmountOfDislikes, &a, &b)
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

	for i := range f.Posts {
		if err = f.Posts[i].ReadCategories(); err != nil {
			return err
		}
	}

	return nil
}

//ReadLDmyPostsByID reads all liked and disliked posts and my own posts ...
func (f *Forum) ReadLDmyPostsByID() error {
	rows, err := Db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.posted_on, posts.user_id, user.username, posts.amount_likes, posts.amount_dislikes, liked.liked, liked.disliked 
FROM posts
LEFT JOIN
liked ON posts.id = liked.post_id and liked.auth_id = ?
INNER JOIN user ON user.id = posts.user_id
WHERE ((liked.liked = 1 and liked.disliked = 0) OR (liked.liked = 0 and liked.disliked = 1) OR posts.user_id = ?) AND posts.isvalid = 1
	`, f.User.ID, f.User.ID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var a, b interface{}
		post := Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.FullText, &post.Date, &post.AuthID, &post.AuthName, &post.AmountOfLikes, &post.AmountOfDislikes, &a, &b)
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

	if err := rows.Err(); err != nil {
		return err
	}

	for i := range f.Posts {
		if err = f.Posts[i].ReadCategories(); err != nil {
			return err
		}
	}
	return nil
}

//ReadComLikesDislikes reads comments with their likes and disliked : registerd user
func (f *Forum) ReadComLikesDislikes() error {
	row, err := Db.Query(`SELECT comments.id, comments.comment, comments.auth_id, comments.post_id, posts.title, comlikes.liked, comlikes.disliked 
	FROM comments 
	LEFT JOIN comlikes ON comments.id = comlikes.com_id AND comlikes.auth_id = ? inner join posts ON posts.id = comments.post_id 
	WHERE ((comlikes.liked = 1 AND comlikes.disliked = 0) OR (comlikes.liked = 0 AND comlikes.disliked = 1) OR comments.auth_id = ?) AND posts.isvalid = 1 
	`, f.User.ID, f.User.ID)
	if err != nil {
		return err
	}

	for row.Next() {
		comm := Comment{}
		var a, b interface{}
		err := row.Scan(&comm.ID, &comm.Comm, &comm.AuthID, &comm.PostID, &comm.PostName, &a, &b)
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
		f.Comments = append(f.Comments, comm)
	}

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

//ReadLPostsByID gets all liked posts from database
func (f *Forum) ReadLPostsByID() error {
	rows, err := Db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.posted_on, user.id as userid, user.username, posts.amount_likes, posts.amount_dislikes, liked.liked, liked.disliked 
FROM posts
INNER JOIN
liked ON posts.id = liked.post_id 
INNER JOIN user ON user.id = posts.user_id
WHERE liked.liked = 1 AND liked.disliked = 0 AND liked.auth_id = ? and posts.valid = 1 and posts.isvalid = 1
	`, f.User.ID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var a interface{}
		var b interface{}
		post := Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.FullText, &post.Date, &post.AuthID, &post.AuthName, &post.AmountOfLikes, &post.AmountOfDislikes, &a, &b)
		if err != nil {
			return err
		}
		post.Sel.IsLike = a.(bool)
		post.Sel.IsDislike = b.(bool)
		f.Posts = append(f.Posts, post)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for i := range f.Posts {
		if err = f.Posts[i].ReadCategories(); err != nil {
			return err
		}
	}
	return nil
}

//ReadMyPostsByID gets all created posts from database
func (f *Forum) ReadMyPostsByID() error {
	rows, err := Db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.posted_on, user.id as userid, user.username,  posts.amount_likes, posts.amount_dislikes, liked.liked, liked.disliked 
FROM posts
INNER JOIN user ON user.id = posts.user_id
Left JOIN liked ON CASE 
WHEN liked.post_id = posts.id AND liked.auth_id = ? THEN 1
	ELSE 0
	END
 WHERE posts.user_id = ? and posts.valid = 1 and posts.isvalid = 1
	`, f.User.ID, f.User.ID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var a interface{}
		var b interface{}
		post := Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.FullText, &post.Date, &post.AuthID, &post.AuthName, &post.AmountOfLikes, &post.AmountOfDislikes, &a, &b)
		if err != nil {
			fmt.Println(err)
			continue
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

	if err := rows.Err(); err != nil {
		return err
	}

	for i := range f.Posts {
		if err = f.Posts[i].ReadCategories(); err != nil {
			return err
		}
	}
	return nil
}
