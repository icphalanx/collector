package collector

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"io"
	"net/http"
	"os"
	"strconv"
)

type WebHost struct {
	ID            int    `json:"id"`
	HumanName     string `json:"humanName"`
	DidInvalidate bool   `json:"didInvalidate"`
}

func (c *Collector) setupWeb() {
	cookieStore := sessions.NewCookieStore([]byte("ex3Xu5J6nYYn22MPTJci15UsnRRh1FY2"))
	cookieStore.Options = &sessions.Options{
		Domain:   "collector.phalanx.lukegb.com",
		Path:     "/api/",
		MaxAge:   3600 * 8,
		HttpOnly: true,
	}

	s := func(r *http.Request) (*sessions.Session, error) {
		return cookieStore.Get(r, "PHALANXSESSID")
	}
	getUser := func(r *http.Request) (*DBUser, *sessions.Session, error) {
		sess, err := s(r)
		if err != nil {
			return nil, sess, err
		}

		u, ok := sess.Values["user"]
		if !ok {
			delete(sess.Values, "user")
			return nil, sess, nil
		}

		u2, ok := u.(int)
		if !ok {
			delete(sess.Values, "user")
			return nil, sess, nil
		}

		u3, err := c.ingestor.GetUserByID(u2)

		return u3, sess, err
	}

	r := pat.New()
	r.Get("/api/hosts/{hostID:[0-9]+}/logs", func(w http.ResponseWriter, r *http.Request) {
		user, _, err := getUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		j := json.NewEncoder(w)

		if len(user.HostIDs) == 0 && !user.SuperUser {
			j.Encode([]*DBLogLine{})
			return
		}

		vars := mux.Vars(r)
		hostIDStr := vars["hostID"]
		hostID, _ := strconv.Atoi(hostIDStr)

		logs, err := c.ingestor.GetLogsForHost(hostID, -100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		fmt.Println(j.Encode(logs))
	})
	r.Get("/api/hosts", func(w http.ResponseWriter, r *http.Request) {
		user, _, err := getUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hosts := map[string]WebHost{}
		j := json.NewEncoder(w)

		if len(user.HostIDs) == 0 && !user.SuperUser {
			j.Encode(hosts)
			return
		}

		hostsSlice, err := c.ingestor.GetHosts(user.HostIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, h := range hostsSlice {
			ho := WebHost{
				ID:            h.ID,
				HumanName:     h.HumanName,
				DidInvalidate: true,
			}
			hosts[fmt.Sprintf("%d", h.ID)] = ho
		}

		j.Encode(hosts)
	})
	r.Get("/api/user/me", func(w http.ResponseWriter, r *http.Request) {
		user, sess, err := getUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var ret struct {
			LoggedIn  bool   `json:"loggedIn"`
			Username  string `json:"username,omitempty"`
			SuperUser *bool  `json:"superuser,omitempty"`
		}
		sess.Save(r, w)

		j := json.NewEncoder(w)
		if user == nil {
			ret.LoggedIn = false
		} else {
			ret.LoggedIn = true
			ret.Username = user.Username
			ret.SuperUser = &user.SuperUser
		}
		j.Encode(ret)

	})
	r.Post("/api/user/login", func(w http.ResponseWriter, r *http.Request) {
		user, sess, err := getUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user != nil {
			http.Error(w, "already logged in", http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "something went wrong handling your request", http.StatusBadRequest)
			return
		}

		u2, err := c.ingestor.GetUser(r.Form.Get("username"))
		if err != nil {
			http.Error(w, "something went wrong handling your request", http.StatusBadRequest)
			return
		}

		sess.Values["user"] = u2.ID
		sessions.Save(r, w)

		w.WriteHeader(http.StatusNoContent)
		return
	})
	r.Post("/api/user/logout", func(w http.ResponseWriter, r *http.Request) {
		user, sess, err := getUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "already logged out", http.StatusBadRequest)
			return
		}

		delete(sess.Values, "user")
		sessions.Save(r, w)

		w.WriteHeader(http.StatusNoContent)
		return
	})

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/dist"))))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("web/dist/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		io.Copy(w, f)
	})

	go http.ListenAndServe(":8181", r)
}
