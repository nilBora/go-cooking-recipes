package recipe

import (
   "fmt"
)

type Recipe struct {
    Uuid        string
    Name        string
	Description string
	Text        string
    Image       string
	Labels      string
}

func (recipe Recipe) Create() {
    fmt.Println(recipe)
}
