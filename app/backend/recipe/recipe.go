package recipe

import (
   "fmt"
)

type Recipe struct {
    Name        string
	Description string
	Text        string
    Image       string
	Labels      string
}

func (recipe Recipe) Create() {
    fmt.Println(recipe)
}
