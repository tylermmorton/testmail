package landing

// /!\ THIS FILE IS GENERATED DO NOT EDIT /!\

import (
	"embed"
	"path/filepath"
	"strings"
)

func _tmpl(fsys embed.FS, path string) string {
	builder := &strings.Builder{}
	entries, err := fsys.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			builder.WriteString(_tmpl(fsys, filepath.Join(path, entry.Name())))
		} else {
			byt, err := fsys.ReadFile(filepath.Join(path, entry.Name()))
			if err != nil {
				panic(err)
			}
			builder.Write(byt)
		}
	}
	return builder.String()
}

//go:embed landing.tmpl.html
var LandingPageTmplFS embed.FS

func (t *LandingPage) TemplateText() string {
	return _tmpl(LandingPageTmplFS, ".")
}
