package main

import(
  "fmt"
  "net/http"
  "github.com/sebber/go-wiki-core/entity"
  "github.com/sebber/go-wiki-core/repository"
  "github.com/sebber/go-wiki-core/usecase"
)

var repo = repository.MemoryWikipageRepository{}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/view/"):]
  p, _ := repo.LoadPage(title)

  fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func setup() {
  repo.Pages = make(map[string]*entity.Page)
  SaveWikipageUseCase := usecase.SaveWikipage{repo}
  SaveWikipageUseCase.Execute("A page", []byte("body"))
}

func main() {
  setup()

  http.HandleFunc("/", viewHandler)
  http.ListenAndServe(":8080", nil)
}

