package routes

import (
	"github.com/go-chi/chi"
	"net/http"
	"fmt"
	"github.com/kwhk/sync/api/web/session"
)

func TestRouter(globalSessions *session.Manager) http.Handler {
	r := chi.NewRouter()
	
	// To test session manager
	r.Get("/session", func(w http.ResponseWriter, r *http.Request) {
		sess := globalSessions.SessionStart(w, r)
		ct := sess.Get("countnum")
		if ct == nil {
			sess.Set("countnum", 1)
		} else {
			sess.Set("countnum", (ct.(int) + 1))
		}

		fmt.Println(sess.Get("countnum"))	
	})

	return r
}