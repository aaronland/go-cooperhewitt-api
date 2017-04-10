package template

import (
	gotemplate "html/template"
)

type ShoeboxIndex struct {
	Links []ShoeboxLink
}

type ShoeboxItem struct {
	Objects ShoeboxObject
}

type ShoeboxObject struct {
	Title  string
	URL    string
	Images ShoeboxObjectImages
}

type ShoeboxObjectImages map[string]ShoeboxObjectImage

type ShoeboxObjectImage struct {
	URL string
}

type ShoeboxLink struct {
	Title string
	URL   string
}

func NewShoeboxIndex(name string) (*gotemplate.Template, error) {

	return gotemplate.New(name).Parse(`
<html>
     <head>
	<title>Your Cooper Hewitt Shoebox</title>
     </head>
     <body>

	<ul>
	{{- range .Links}}
	    <li><a href="{{.URL}}">{{.Title}}</a></li>
	{{- end}}
	</ul>
	
     </body>
</html>`)

}

func NewShoeboxItem(name string) (*gotemplate.Template, error) {

	return gotemplate.New(name).Parse(`
<html>
     <head>
	<title>...</title>
     </head>
     <body>
	<img src="" height="" width="" alt="" />
	<a href="">...</a>
     </body>
</html>`)

}
