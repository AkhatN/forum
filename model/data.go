package model

import (
	"database/sql"
	"net/http"
)

//Db is database
var Db *sql.DB

//Forum struct
type Forum struct {
	User               User
	Users              []User
	Posts              []Post
	Post               Post
	IsErr              bool
	ErrMsg             string
	ErrStatus          int
	Comments           []Comment
	AmountNotification int
}

//ReadDataFromDataBase gets all posts from database
func (f *Forum) ReadDataFromDataBase(w http.ResponseWriter, r *http.Request) error {
	if !f.User.IsStatus {
		if err := f.ReadPostFromDB(); err != nil {
			return err
		}
		return nil
	}

	if err := f.ReadPostsForRegisteredUserFromDB(); err != nil {
		return err
	}

	if err := f.ReadCountNotification(); err != nil {
		return err
	}

	return nil
}

//ReadPostFromDB gets all posts for a guest
func (f *Forum) ReadPostFromDB() error {
	row, err := Db.Query(`SELECT posts.id, posts.title, posts.body, posts.posted_on, user.username, amount_likes, amount_dislikes FROM posts
	INNER JOIN user ON posts.user_id = user.id 
	WHERE posts.valid = 1 and posts.isvalid = 1`)
	if err != nil {
		return err
	}

	for row.Next() {
		post := Post{}
		err := row.Scan(&post.ID, &post.Title, &post.FullText, &post.Date, &post.AuthName, &post.AmountOfLikes, &post.AmountOfDislikes)
		if err != nil {
			return err
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

//ReadUsersFromDB ...
func (f *Forum) ReadUsersFromDB() error {
	row, err := Db.Query(`SELECT id, username from user where id != ?`, f.User.ID)
	if err != nil {
		return err
	}

	for row.Next() {
		u := User{}
		err := row.Scan(&u.ID, &u.Username)
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

//ReadPostFromDataBase gets specific post with comments and likes from database
func (f *Forum) ReadPostFromDataBase(w http.ResponseWriter, r *http.Request) error {
	//if not logged user
	if !f.User.IsStatus {
		//if return nil page not found
		if err := f.Post.ReadPostByID(); err != nil {
			return nil
		}

		if err := f.Post.ReadCommensByID(); err != nil {
			return err
		}

		return nil
	}

	//if logged user
	//if return nil page not found
	if err := f.Post.ReadPostForUserFromDB(f.User.ID); err != nil {
		return nil
	}

	if err := f.Post.ReadComLikesForUser(f.User.ID); err != nil {
		return err
	}

	return nil
}

//ReadPostsForRegisteredUserFromDB gets all posts for a registered user
func (f *Forum) ReadPostsForRegisteredUserFromDB() error {
	row, err := Db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.posted_on, user.id as userid, user.username, amount_likes, amount_dislikes, liked.liked, liked.disliked FROM posts
		INNER JOIN user ON posts.user_id = user.id
		LEFT JOIN liked ON CASE
		WHEN liked.post_id = posts.id and liked.auth_id = ? THEN 1
		ELSE 0
		END
	WHERE posts.valid = 1 and posts.isvalid = 1 `, f.User.ID)
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

//ReadCatPostsFromDataBase reads all posts dedicated to a specific topic
func (f *Forum) ReadCatPostsFromDataBase(w http.ResponseWriter, r *http.Request) error {
	//if not logged user
	if !f.User.IsStatus {
		if err := f.ReadPostsByCatName(); err != nil {
			return err
		}
		return nil
	}

	//if logged user
	if err := f.ReadCatPostForRegisteredUserFromDB(); err != nil {
		return err
	}

	if err := f.ReadCountNotification(); err != nil {
		return err
	}

	return nil
}

//ReadLikedDislikedmyPostsFromDataBase Reads liked and disliked posts and including my own posts : registered user
func (f *Forum) ReadLikedDislikedmyPostsFromDataBase(w http.ResponseWriter, r *http.Request) error {
	//if not logged user
	if !f.User.IsStatus {
		return nil
	}

	//if logged user
	if err := f.ReadLDmyPostsByID(); err != nil {
		return err
	}

	// Read Comment with their likes and dislikes ...
	if err := f.ReadComLikesDislikes(); err != nil {
		return err
	}

	return nil
}

//ReadLikedPostsFromDataBase reads all liked posts of a registered user
func (f *Forum) ReadLikedPostsFromDataBase(w http.ResponseWriter, r *http.Request) error {
	//if not logged user
	if !f.User.IsStatus {
		return nil
	}

	//if logged user
	if err := f.ReadLPostsByID(); err != nil {
		return err
	}

	if err := f.ReadCountNotification(); err != nil {
		return err
	}

	return nil
}

//ReadMyPostsFromDataBase reads all created posts of a registered user
func (f *Forum) ReadMyPostsFromDataBase(w http.ResponseWriter, r *http.Request) error {
	//if not logged user
	if !f.User.IsStatus {
		return nil
	}

	//if logged user
	if err := f.ReadMyPostsByID(); err != nil {
		return err
	}

	return nil
}
