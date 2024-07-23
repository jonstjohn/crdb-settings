package api

import (
	"encoding/json"
	"fmt"
	"github.com/jonstjohn/crdb-settings/pkg/settings"
	"net/http"
	"regexp"
)

var (
	SettingsReleaseReWithRelease  = regexp.MustCompile(`^/settings/release/(.+)$`)
	SettingsCompareReWithReleases = regexp.MustCompile(`^/settings/compare/(.+)\.\.(.+)$`)
)

func Serve(url string) {
	mux := http.NewServeMux()
	mux.Handle("/", &SettingsHandler{Url: url})
	http.ListenAndServe(":8080", mux)
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

type SettingsHandler struct {
	Url string
}

func (h *SettingsHandler) CompareSettingsForReleases(w http.ResponseWriter, r *http.Request) {
	matches := SettingsCompareReWithReleases.FindStringSubmatch(r.URL.Path)
	if len(matches) != 3 {
		w.WriteHeader(http.StatusOK) // TODO
		w.Write([]byte("Release must be included"))
		return
	}
	r1 := matches[1]
	r2 := matches[2]

	sm, err := settings.NewSettingsManager(h.Url)
	if err != nil {
		ErrorHandler(w, err)
		return
	}

	s, err := sm.CompareSettingsForReleases(r1, r2)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (h *SettingsHandler) ListSettingsForRelease(w http.ResponseWriter, r *http.Request) {
	matches := SettingsReleaseReWithRelease.FindStringSubmatch(r.URL.Path)
	if len(matches) != 2 {
		w.WriteHeader(http.StatusOK) // TODO
		w.Write([]byte("Release must be included"))
		return
	}
	release := matches[1]

	sm, err := settings.NewSettingsManager(h.Url)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	s, err := sm.GetSettingsForRelease(release)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func ErrorHandler(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadGateway) // TODO
	w.Write([]byte(fmt.Sprintf("%v", err)))
	return
}

func (h *SettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && SettingsReleaseReWithRelease.MatchString(r.URL.Path):
		h.ListSettingsForRelease(w, r)
		return
	case r.Method == http.MethodGet && SettingsCompareReWithReleases.MatchString(r.URL.Path):
		h.CompareSettingsForReleases(w, r)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
