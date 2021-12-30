package pkg

import (
	"html/template"
)

//NewView ...
func NewView(files ...string) (*template.Template, error) {
	files2 := []string{"views/bootstrap.html", "views/navbar.html"}
	files = append(files, files2...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return t, nil
}
