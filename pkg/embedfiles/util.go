package embedfiles

import (
	"embed"
)

//go:generate rm -rf ./assets
//go:generate cp -r ../../assets ./assets
//go:embed assets/*
var assets embed.FS

func InnerReadFile(filename string) ([]byte, error) {
	return assets.ReadFile(filename)
}
