package logic

import "net/http"

func InitHTTP() (err error) {
	httpServerMux := http.NewServeMux()
	httpServerMux.HandleFunc("/api/v1/push", Push)
	return err
}

func Push(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
	}

	var (
		auth = r.URL.Query().Get("auth")
	)

	_, _ = getRouter(auth)

}
