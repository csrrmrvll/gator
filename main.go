package main

import (
	"fmt"

	"github.com/csrrmrvll/gator/internal/config"
)

func main() {
	cfg := config.Read()
	config.SetUser(cfg)
	cfg = config.Read()
	fmt.Println(cfg)
}
