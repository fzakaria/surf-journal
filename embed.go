package surf_journal

import "embed"

//go:generate tailwindcss -i ./static/css/input.css -o ./static/css/output.css

//go:embed templates/*
var TemplateFS embed.FS

//go:embed static/css/*
var StaticFS embed.FS
