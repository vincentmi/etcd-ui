package main

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io"
	"html/template"
	"log"
	"flag"
	"path/filepath"
)


type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}


func parseTemplates(pattern string) *template.Template{
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}
	if len(filenames) == 0 {
		log.Fatal("no template")
	}
	t := template.New("")
	t.Delims("<%" ,"%>")

	t1 , err  := t.ParseFiles(filenames...)

	if err != nil {
		panic(err)

	}else{
		return t1
	}
}


func main() {

	listen := flag.String("listen" , ":9600" , "http listen")

	etcd := flag.String("etcd" , "http://127.0.0.1:2379" , "ETCD endpoints,separator with comma.")

	e := echo.New()

	log.Print(*etcd)

	t := &Template{
		templates: parseTemplates("views/*.html"),
	}



	e.Renderer = t

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/js" , "public/js")

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", "Vincent")
	})

	// Start server
	e.Logger.Fatal(e.Start(*listen))
}
