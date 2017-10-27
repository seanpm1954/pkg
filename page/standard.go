package page

import "github.com/biz/templates"

func standard() {
	// template used as the main base view
	templates.AddPartial("standard.wrapper", `
<!DOCTYPE html>
<html>
	<head>
		{{ CacheLinks .Page.Links }}
		<link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
		<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700" type="text/css">
		<meta name=viewport content="width=device-width, initial-scale=1">
		<title>{{.Page.Title}}</title>

		<!-- TODO(move) -->
		{{ if .Page.FaviconHTML }}
			{{ .Page.FaviconHTML }}
		{{ end }}

		<meta name="apple-mobile-web-app-capable" content="yes">
		<meta name="apple-mobile-web-app-status-bar-style" content="black">

		{{ if .Config }}
		<script>
			var config = {{ Json .Config }};
		</script>
		{{ end }}
	</head>
	<body>
		<div class="mdl-layout mdl-js-layout mdl-layout--fixed-header {{ if (and .Menu (not .Page.CollapseMenu)) }}mdl-layout--fixed-drawer{{end}}">
			{{if .Page.Header}}{{template "standard.header" .}}{{end}}
			<main class="mdl-layout__content">
				{{ template "messages" . }}
				{{ if .Page.BreadCrumbs }}
				<div id="breadcrumbs">
					{{ range $bc := .Page.BreadCrumbs }}
						{{ if $bc.Link }}
						<a href="{{SafeURL $bc.Link}}" class="mdl-button mdl-js-button mdl-button--raised">{{$bc.Label}}</a>
						{{ else }}
						<span class="breadcrumb-tail">{{ $bc.Label }}</span>
						{{ end }}
					{{ end }}
				</div>
				{{ end }}
				<div class="content">
					{{template "body" .}}
				</div>
			</main>
		</div>
		<script>
		(function() {
			var back = document.querySelector(".back-button");
			if (back) {
				back.addEventListener("click", function() {
					window.history.back();
				});
			}
		}())
		</script>

		{{ template "scripts-no-bust" . }}
		{{ CacheScripts .Page.Scripts }}
		{{ template "scripts-no-bust-post-scripts" . }}

		{{ block "footer" . }}{{end}}
		{{ block "footerMisc" . }}{{end}}

		<div class="mdl-js-snackbar mdl-snackbar">
			<div class="mdl-snackbar__text"></div>
			<button class="mdl-snackbar__action" type="button"></button>
		</div>
	</body>
</html>
	`)

	// template will render the header and navigation
	templates.AddPartial("standard.header", `
{{ if .Page.Header }}
<header class="mdl-layout__header">
	<div class="mdl-layout__header-row">
        {{ if .Page.GoBack}}
		{{ GoBack .Req }} &nbsp;&nbsp {{.Page.Title}}
        {{ else }}
        {{.Page.Title}}
        {{ end }}
		<!-- Add spacer, to align navigation to the right -->

		<div class="mdl-layout-spacer"></div>

		<nav class="mdl-navigation">
			{{ if eq .Session.UserID 0 }}
				<a class="mdl-navigation__link" href="/login">Login</a>
			{{ else }}
				<a class="mdl-navigation__link" href="/logout">{{ .Session.Username }} -- Logout</a>
			{{ end }}
		</nav>
	</div>
</header>
{{end}}

{{ if .Menu }}
<div class="mdl-layout__drawer">
	{{ if .Page.Header.Logo.Img }}
		<div class="logo-image">
			<a href="/" title="Home">
				<img src="{{ .Page.Header.Logo.Img }}" alt="{{ .Page.Header.Title }} Logo" {{ if .Page.Header.Logo.Height }}height="{{ .Page.Header.Logo.Height }}"{{end}}{{ if .Page.Header.Logo.Width }}width="{{ .Page.Header.Logo.Width }}"{{end}}></img>
			</a>
		</div>
	{{ else }}
		<span class="mdl-layout-title logo-title"><a class="mdl-navigation__link" href="/" title="home">{{.Page.Header.Title}}</a></span>
	{{ end }}

	{{ if .Menu.Title }}
		<div class="menu-title">
			{{.Menu.Title}}
		</div>
	{{ end }}

	<nav class="mdl-navigation">
	{{ range $item := .Menu.Items }}
		<a class="mdl-navigation__link {{if $item.IsActive}}active-link{{end}}" href="{{$item.Href}}" title="{{ $item.Title }}">
			{{$item.Name}}
		</a>
		{{ range $child := $item.SubItems }}
			<a {{$child.PrintAttr "class" "mdl-navigation__link" }} href="{{$child.Href}}" title="{{ $child.Title }}">
				{{$child.Name}}
			</a>
		{{ end }}
	{{ end }}
	</nav>
</div>
{{ end }}
	`)

	templates.AddPartial("skeleton.base", `
<!DOCTYPE html>
<html>
	<head>
		{{ CacheLinks .Page.Links }}
		{{ CacheScripts .Page.Scripts }}
		<title>{{.Page.Title}}</title>
		<link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
		<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700" type="text/css">
		<meta name=viewport content="width=device-width, initial-scale=1">
		<meta name="apple-mobile-web-app-capable" content="yes">
		<meta name="apple-mobile-web-app-status-bar-style" content="black">
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
	</head>
	<body {{ if .Page.BodyClass }}class="{{.Page.BodyClass}}"{{end}}>
	{{ template "body" . }}
	</body>
</html>
	`)
}
