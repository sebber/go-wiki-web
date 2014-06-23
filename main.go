package main

import(
  "fmt"
  "log"
  "net/http"
  
  "github.com/gorilla/mux"
  "github.com/gorilla/schema"
  "github.com/hoisie/mustache"  

  "github.com/sebber/go-wiki-core/entity"
  "github.com/sebber/go-wiki-core/repository"
  "github.com/sebber/go-wiki-core/usecase"
)

type PageInput struct {
  Title string
  Body string
}

var (
  repo = repository.MemoryWikipageRepository{}
  decoder = schema.NewDecoder()
)

func listHandler(w http.ResponseWriter, r *http.Request) {
  pages, _ := repo.All()

  renderPageTemplate(w, "list", pages)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  title := vars["title"]

  p, _ := repo.GetByTitle(title)

  renderPageTemplate(w, "view", p)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
  renderTemplate(w, "create")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()

  if err != nil {
    log.Fatal(err)
    http.Redirect(w, r, "/", 302)
  }

  page := new(PageInput)
  err = decoder.Decode(page, r.PostForm)

  if err != nil {
    log.Fatal(err)
    http.Redirect(w, r, "/", 302)
  }

  SaveWiki := usecase.SaveWikipage{repo}
  SaveWiki.Execute(page.Title, []byte(page.Body))

  http.Redirect(w, r, "/view/"+ page.Title, 302)
}

func setup() {
  repo.Pages = make(map[string]*entity.Page)
  SaveWikipageUseCase := usecase.SaveWikipage{repo}
  SaveWikipageUseCase.Execute("A page", []byte("body"))
}

func main() {
  setup()

  r := mux.NewRouter()
  r.HandleFunc("/view/create",  createHandler)
  r.HandleFunc("/view/save",    saveHandler)
  r.HandleFunc("/view/{title}", viewHandler)
  r.HandleFunc("/", listHandler)
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


func renderPageTemplate(w http.ResponseWriter, tmpl string, p *entity.Page) {
  view := PageView{Content: p}
  output := mustache.RenderFile("templates/"+tmpl+".mustache", view)
  fmt.Fprintf(w, output)
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
  output := mustache.RenderFile("templates/"+tmpl+".mustache")
  fmt.Fprintf(w, output)
}
