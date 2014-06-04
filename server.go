package main

import (
  // "fmt"
  "html/template"
  "net/http"
  "regexp"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9][a-zA-Z0-9_]*)$")

func renderTemplate(writer http.ResponseWriter, name string, page *Page) {
  // fmt.Println("Rendering template: ", name, " with page ", page)
  err := templates.ExecuteTemplate(writer, name + ".html", page)
  if err != nil {
    http.Error(writer, err.Error(), http.StatusInternalServerError)
  }
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    match := validPath.FindStringSubmatch(request.URL.Path)
    if match == nil {
      http.NotFound(writer, request)
      return
    }

    fn(writer, request, match[2])
  }
}

func viewHandler(writer http.ResponseWriter, request *http.Request, title string) {
  page, err := loadPage(title)
  if err != nil {
    http.Redirect(writer, request, "/edit/" + title, http.StatusFound)
    return
  }

  renderTemplate(writer, "view", page)
}

func editHandler(writer http.ResponseWriter, request *http.Request, title string) {
  page, err := loadPage(title)
  if err != nil {
    page = &Page{Title: title}
  }

  renderTemplate(writer, "edit", page)
}

func saveHandler(writer http.ResponseWriter, request *http.Request, title string) {
  body := request.FormValue("body")
  page := &Page{Title: title, Body: []byte(body)}
  err := page.save()
  if err != nil {
    http.Error(writer, err.Error(), http.StatusInternalServerError)
    return
  }

  http.Redirect(writer, request, "/view/" + title, http.StatusFound)
}

func main() {
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))

  http.ListenAndServe(":8888", nil)
}

