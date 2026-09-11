package api

import "net/http"

// Build identifies the running build. The UI shows it in the footer so you
// can tell at a glance which commit a preview environment is serving.
type Build struct {
	Env    string `json:"env"`
	GitSHA string `json:"gitSha"`
}

func (s *Server) version(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.build)
}
