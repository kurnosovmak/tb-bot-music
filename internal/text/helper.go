package text

import "embed"

//go:embed texts/*.html
var FS embed.FS

func MustLoad(name string) string {
	b, err := FS.ReadFile("texts/" + name)
	if err != nil {
		panic(err)
	}
	return string(b)
}
