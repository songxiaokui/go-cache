package internal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

/*
@Time    : 2021/3/3 22:31
@Author  : austsxk
@Email   : austsxk@163.com
@File    : http.go
@Software: GoLand
*/

type CacheServers struct {
	// embedding Cache interface
	Cache
}

func (cs *CacheServers) Listen() {
	// add funcHandler /cache/
	http.Handle("/cache/", cs.cacheHandler())

	// add funcHandler /state/
	http.Handle("/status/", cs.statusHandler())

	// listen server
	http.ListenAndServe("127.0.0.1:9999", nil)
}

func New(cache *Cache) *CacheServers {
	return &CacheServers{*cache}
}

type CacheHandler struct {
	*CacheServers
}

func (ch *CacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get request url
	requestUrl := r.URL.EscapedPath()
	key := strings.Split(requestUrl, "/")[2]
	// 400 bad request
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get request methods
	methods := r.Method

	// append method to deal with
	// add
	if methods == http.MethodPut {
		// read body info
		b, _ := ioutil.ReadAll(r.Body)
		if len(b) != 0 {
			body := string(b)
			log.Print("add: ", body)
			err := ch.Set(key, body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		return
	}
	// get
	if methods == http.MethodGet {
		r, err := ch.Get(key)
		log.Print("get: ", r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(r) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(r))
		return
	}
	// delete
	if methods == http.MethodDelete {
		err := ch.Delete(key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	// any methods
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func (cs *CacheServers) cacheHandler() http.Handler {
	return &CacheHandler{cs}
}

// status router
type StatusHandler struct {
	*CacheServers
}

func (sh *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// convert data to json
	d, err := json.Marshal(sh.GetStat())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(d)
	return
}

func (cs *CacheServers) statusHandler() http.Handler {
	return &StatusHandler{cs}
}
