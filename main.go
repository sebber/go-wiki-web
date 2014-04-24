package main

import(
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/hoisie/mustache"

  "github.com/sebber/go-wiki-core/entity"
  "github.com/sebber/go-wiki-core/repository"
  "github.com/sebber/go-wiki-core/usecase"
)

var repo = repository.MemoryWikipageRepository{}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/view/"):]
  p, _ := repo.GetByTitle(title)

  fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func createHandler(w http.ResponseWriter, r *http.Request) {

}

func setup() {
  repo.Pages = make(map[string]*entity.Page)
  SaveWikipageUseCase := usecase.SaveWikipage{repo}
  SaveWikipageUseCase.Execute("A page", []byte("body"))
}

func main() {
  setup()

  r := mux.NewRouter()
  r.HandleFunc("/", viewHandler)
  http.ListenAndServe(":8080", nil)
  http.Handle("/", r)
}


type PageView struct {
  Content *entity.Page
}

func (page Page) Title() string {
  return page.Content.Title
}

func (page Page) Body() string {
  return string(page.Content.Body)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *entity.Page) {
  view := view.Page{Content: p}
  output := mustache.RenderFile("templates/"+tmpl+".mustache", view)
  fmt.Fprintf(w, output)
}

