package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
)

func main() {
	// /-/css/strle.css => publ;ic/css/style.css
	// / index
	http.HandleFunc("/", indexHandler)
	http.Handle("/-/", http.StripPrefix("/-", http.FileServer(noDir{http.Dir("public")})))
	http.ListenAndServe(":8080", nil)

}

type noDir struct {
	http.Dir
}

func (d noDir) Open(name string) (http.File, error) {
	f, err := d.Dir.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}

type indexData struct {
	Name string
	List []string
}

var t = template.Must(template.ParseFiles("index.tmpl")) // create template outside handler (for speedup request-response)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := indexData{
		Name: "nuSapb",
		List: []string{
			"Go",
			"C",
			"JS",
		},
	}

	log.Println(r.URL.Path)

	err := t.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
