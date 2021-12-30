package pkg

import (
	"forum/model"
	"strings"
)

//CheckUniq ...
func CheckUniq(cats *[]model.Category, cs []string) {
	for _, v := range cs {
		isUnique := true
		cat := model.Category{Name: strings.ToLower(v)}
		for _, c := range *cats {
			if c.Name == cat.Name {
				isUnique = false
				break
			}
		}
		if isUnique {
			*cats = append(*cats, cat)
		}
	}

}
