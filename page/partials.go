package page

import "github.com/biz/templates"

func partials() {
	templates.AddPartial("messages", `
		{{ template "infoMessage" . }}
		{{ template "errorMessage" . }}
	`)

	templates.AddPartial("errorMessage", `
{{ if .Page.ErrorMessage }}
	<div class="alert alert__error">
	{{ .Page.ErrorMessage }}
	</div>
{{ end }}
	`)

	templates.AddPartial("infoMessage", `
{{ if .Page.InfoMessage }}
	<div class="alert alert__info">
	{{ .Page.InfoMessage }}
	</div>
{{ end }}
	`)

	// template used to add css links to the page
	templates.AddPartial("links", `
{{ range .Page.Links }}
	<link rel="stylesheet" href="{{.}}" type="text/css" \>
{{ end }}
	`)

	// template used to add script tags to the page
	templates.AddPartial("scripts", `
{{ range .Page.Scripts }}
	<script src="{{.}}"></script>
{{ end }}
	`)

	// template used to add script tags to the page
	templates.AddPartial("scripts-no-bust", `
{{ range .Page.ScriptsNoBust }}
	<script src="{{.}}"></script>
{{ end }}
	`)

	// template used to add script tags to the page
	templates.AddPartial("scripts-no-bust-post-scripts", `
{{ range .Page.ScriptsNoBustPostScripts }}
	<script async defer src="{{.}}"></script>
{{ end }}
	`)
}
