package pkg

import (
	"fmt"
	"net/mail"
	"strings"
)

//EmailValidate checks for valid email
func EmailValidate(email string) error {
	if len(email) == 0 {
		return fmt.Errorf("empty line")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}

	return nil
}

//CharactersValidate checks for valid charachters
func CharactersValidate(word string) error {
	if len(word) == 0 {
		return fmt.Errorf("empty line")
	}
	if len(word) > 40 {
		return fmt.Errorf("large amount of characters")
	}
	for i := range word {
		if word[i] < 33 || word[i] > 125 {
			return fmt.Errorf("invalid characters")
		}
	}

	return nil
}

func EmptySpaceCheck(s ...string) error {
	for i := range s {
		if len(strings.TrimSpace(s[i])) == 0 {
			return fmt.Errorf("Empty line is not allowed")
		}
	}
	return nil
}
