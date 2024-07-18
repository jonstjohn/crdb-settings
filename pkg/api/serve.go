package api

import "net/http"

func Serve() {
	mux := http.NewServeMux()
	mux.Handle("/", &homeHandler{})
	http.ListenAndServe(":8080", mux)
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}
