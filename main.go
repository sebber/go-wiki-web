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
  vars := mux.Vars(r)
  title := vars["title"]

  p, _ := repo.GetByTitle(title)

  //fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
  renderTemplate(w, "view", p)
}

//func createHandler(w http.ResponseWriter, r *http.Request) {

//}

func setup() {
  repo.Pages = make(map[string]*entity.Page)
  SaveWikipageUseCase := usecase.SaveWikipage{repo}
  SaveWikipageUseCase.Execute("A page", []byte("body"))
}

func main() {
  setup()

  r := mux.NewRouter()
  r.HandleFunc("/view/{title}", viewHandler)
  http.Handle("/", r)
  http.ListenAndServe(":8080", nil)
}


type PageView struct {
  Content *entity.Page
}

func (page PageView) Title() string {
  return page.Content.Title
}

func (page PageView) Body() string {
  return string(page.Content.Body)
}


func renderTemplate(w http.ResponseWriter, tmpl string, p *entity.Page) {
  view := PageView{Content: p}
  output := mustache.RenderFile("templates/"+tmpl+".mustache", view)
  fmt.Fprintf(w, output)
}

