package main

import (
	"fmt"
	"github.com/newpanjing/seedgo/internal/model"
)

func main() {
	var id model.ID = 123

	id += 1
	idStr := id.String()
	fmt.Printf("id: %d, idStr: %s\n", id, idStr)
}
