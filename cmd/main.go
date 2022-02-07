package main

import "github.com/SukiEva/shortener"

func main() {
	if s, err := shortener.New(); err == nil {
		s.Serve()
	}
}
