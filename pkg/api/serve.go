package api

import (
	"encoding/json"
	"fmt"
	"github.com/jonstjohn/crdb-settings/pkg/metrics"
	"github.com/jonstjohn/crdb-settings/pkg/releases"
	"github.com/jonstjohn/crdb-settings/pkg/settings"
	"net/http"
	"regexp"
)

var (
	SettingsReleaseReWithRelease  = regexp.MustCompile(`^/settings/release/(.+)$`)
	SettingsCompareReWithReleases = regexp.MustCompile(`^/settings/compare/(.+)\.\.(.+)$`)
	SettingsHistoryReWithSetting  = regexp.MustCompile(`^/settings/history/(.+)$`)
	SettingsDetailReWithSetting   = regexp.MustCompile(`^/settings/detail/(.+)$`)
	ReleasesRe                    = regexp.MustCompile(`^/releases/list$`)
	MetricsReleaseReWithRelease   = regexp.MustCompile(`^/metrics/release/(.+)$`)
	MetricsCompareReWithReleases  = regexp.MustCompile(`^/metrics/compare/(.+)\.\.(.+)$`)
	//	MetricsHistoryReWithSetting   = regexp.MustCompile(`^/metrics/history/(.+)$`)
	//	MetricsDetailReWithSetting    = regexp.MustCompile(`^/metrics/detail/(.+)$`)
)

func Serve(url string) {
	mux := http.NewServeMux()
	mux.Handle("/", &SettingsHandler{Url: url})
	http.ListenAndServe(":8080", mux)
}

type SettingsHandler struct {
	Url string
}

func (h *SettingsHandler) HistoryForSetting(w http.ResponseWriter, r *http.Request) {
	matches := SettingsHistoryReWithSetting.FindStringSubmatch(r.URL.Path)
	if len(matches) != 2 {
		w.WriteHeader(http.StatusOK) // TODO
		w.Write([]byte("Release must be included"))
		return
	}
	setting := matches[1]

	sm, err := settings.NewSettingsManager(h.Url)
	if err != nil {
		ErrorHandler(w, err)
		return
	}

	s, err := sm.HistoryForSetting(setting)
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

	//w.Write([]byte(fmt.Sprintf("History for '%s'", setting)))
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

func (h *SettingsHandler) ListReleases(w http.ResponseWriter, r *http.Request) {
	rm, err := releases.NewReleasesManager(h.Url)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	releases, err := rm.GetReleases()
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	jsonBytes, err := json.Marshal(releases)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *SettingsHandler) SettingDetail(w http.ResponseWriter, r *http.Request) {
	matches := SettingsDetailReWithSetting.FindStringSubmatch(r.URL.Path)
	if len(matches) != 2 {
		w.WriteHeader(http.StatusOK) // TODO
		w.Write([]byte("Release must be included"))
		return
	}
	setting := matches[1]

	sm, err := settings.NewSettingsManager(h.Url)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	s, err := sm.GetSettingDetail(setting)
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
	_, err = w.Write(jsonBytes)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
}

func (h *SettingsHandler) ListMetricsForRelease(w http.ResponseWriter, r *http.Request) {
	matches := MetricsReleaseReWithRelease.FindStringSubmatch(r.URL.Path)
	if len(matches) != 2 {
		w.WriteHeader(http.StatusOK) // TODO
		w.Write([]byte("Release must be included"))
		return
	}
	release := matches[1]

	m, err := metrics.NewManager(h.Url)
	if err != nil {
		ErrorHandler(w, err)
		return
	}

	ms, err := m.GetMetricsForRelease(release)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	jsonBytes, err := json.Marshal(ms)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *SettingsHandler) CompareMetricsForReleases(w http.ResponseWriter, r *http.Request) {
	matches := MetricsCompareReWithReleases.FindStringSubmatch(r.URL.Path)
	if len(matches) != 3 {
		w.WriteHeader(http.StatusOK) // TODO
		w.Write([]byte("Release must be included"))
		return
	}
	r1 := matches[1]
	r2 := matches[2]

	m, err := metrics.NewManager(h.Url)
	if err != nil {
		ErrorHandler(w, err)
		return
	}

	s, err := m.CompareMetricsForReleases(r1, r2)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch {
	case r.Method == http.MethodGet && SettingsReleaseReWithRelease.MatchString(r.URL.Path):
		h.ListSettingsForRelease(w, r)
		return
	case r.Method == http.MethodGet && SettingsCompareReWithReleases.MatchString(r.URL.Path):
		h.CompareSettingsForReleases(w, r)
		return
	case r.Method == http.MethodGet && SettingsHistoryReWithSetting.MatchString(r.URL.Path):
		h.HistoryForSetting(w, r)
		return
	case r.Method == http.MethodGet && ReleasesRe.MatchString(r.URL.Path):
		h.ListReleases(w, r)
	case r.Method == http.MethodGet && SettingsDetailReWithSetting.MatchString(r.URL.Path):
		h.SettingDetail(w, r)
	case r.Method == http.MethodGet && MetricsReleaseReWithRelease.MatchString(r.URL.Path):
		h.ListMetricsForRelease(w, r)
	case r.Method == http.MethodGet && MetricsCompareReWithReleases.MatchString(r.URL.Path):
		h.CompareMetricsForReleases(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
