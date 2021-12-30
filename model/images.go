package model

import "errors"

//WriteImagetoDb writing an image path to Sqlite DataBase, returning error if it couldn't insert ...
func (p *Post) WriteImagetoDb() error {
	_, err := Db.Exec(`
	INSERT OR REPLACE INTO image_post (post_id, pimage) 
	VALUES (?, ?)`, p.ID, p.Image)
	if err != nil {
		return err
	}

	return nil
}

//ReadImagefromDb reading image from Sqlite DataBase, returning error if it couldn't read ...
func (p *Post) ReadImagefromDb() (string, error) {
	row := Db.QueryRow(`SELECT pimage FROM image_post WHERE post_id = ?`, p.ID)
	image := ""
	err := row.Scan(&image)
	if err != nil {
		return image, err
	}

	if image == "" {
		return image, errors.New("No image for the post")
	}

	p.IsImage = true
	return image, nil
}

//DeleteImageFromDb reading image from Sqlite DataBase, returning error if it couldn't insert ...
func (p *Post) DeleteImageFromDb() error {
	_, err := Db.Exec(`UPDATE image_post SET pimage = ? WHERE post_id = ?`, "", p.ID)
	if err != nil {
		return err
	}

	return nil
}
