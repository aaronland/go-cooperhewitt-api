package template

import (
	gotemplate "html/template"
)

type ShoeboxIndex struct {
	Items []*ShoeboxIndexItem
}

type ShoeboxIndexItem struct {
	Id    int64
	URL   string
	Title string
}

type ShoeboxItem struct {
	Title  string
	Object ShoeboxObject
}

type ShoeboxObject struct {
	Title string
	URL   string
	// Images ShoeboxObjectImages
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
	{{- range .Items}}
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
	<title>{{.Title}}</title>
     </head>
     <body>
	<img src="" height="" width="" alt="" />
	<h3>{{.Object.Title}}</h3>
	<a href="{{.Object.URL}}">{{.Object.URL}}</a>
     </body>
</html>`)

}
