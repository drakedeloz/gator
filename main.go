package main

import (
	"fmt"

	"github.com/drakedeloz/gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("Drake")
	fmt.Println(config.Read())
}
