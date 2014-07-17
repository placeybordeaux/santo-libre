package libre

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sort"
	"strings"

	"github.com/go-martini/martini"
)

var LastResponses map[string][]http.Response
var LastRequests map[string][]http.Request

func init() {
	LastResponses = make(map[string][]http.Response)
	LastRequests = make(map[string][]http.Request)
}

func RecordLastRequests(c martini.Context, r *http.Request, routes martini.Routes) {
	LastRequests[r.URL.Path] = append(LastRequests[r.URL.Path], *r)
}

func ExposeRoutesMD(m martini.Routes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		md := routes_to_md(m)
		_, err := w.Write([]byte(md))
		if err != nil {
			panic(err)
		}
	}
}

func ExposeRoutes(m martini.Routes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		md := routes_to_md(m)
		err := ioutil.WriteFile("/tmp/santo-libre.md", []byte(md), 0644)
		if err != nil {
			panic(err)
		}

		cmd := exec.Command("aglio", "-i", "/tmp/santo-libre.md", "-o", "-")
		out, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "text/html")
		_, err = w.Write(out)
		if err != nil {
			panic(err)
		}
	}
}

type routesSorter []martini.Route

func (rs routesSorter) Len() int {
	return len(rs)
}

func (rs routesSorter) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs routesSorter) Less(i, j int) bool {
	return rs[i].Pattern() < rs[j].Pattern()
}

func routes_to_md(routes martini.Routes) string {
	s := ""
	s += "FORMAT: 1A\nHOST: http://made.up\n"
	rs := routes.All()
	sort.Sort(routesSorter(rs))
	for i, route := range rs {
		if i == 0 || strings.Split(rs[i].Pattern(), "/")[1] != strings.Split(rs[i-1].Pattern(), "/")[1] {
			s += fmt.Sprintf("\n# group %s\n", strings.Split(route.Pattern(), "/")[1])
		}
		s += fmt.Sprintf("\n## %s %s\n", route.Method(), route.Pattern())
		reqs, ok := LastRequests[route.Pattern()]
		if ok {
			s += "+ <keyword> Payload\n\n"
			s += "\t Body section\n\n"
			for _, req := range reqs {
				s += fmt.Sprintf("\t+ Headers\n\n")
				for header, values := range req.Header {
					for _, v := range values {
						s += fmt.Sprintf("\t\t%s: %s\n", header, v)
					}
				}
				s += "\n"
			}
		}
	}
	return s
}
