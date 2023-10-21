package embedfiles

import (
	"embed"
)

//go:generate cp -r ../../assets ./assets
//go:embed assets/*
var assets embed.FS

func ReadFile(filename string) ([]byte, error) {
	return assets.ReadFile(filename)
}
