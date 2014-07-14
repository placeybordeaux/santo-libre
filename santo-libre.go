package libre

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sort"

	"github.com/go-martini/martini"
)

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
	for _, route := range rs {
		s += fmt.Sprintf("## Default Name [%s]\n", route.Pattern())
		s += fmt.Sprintf("\n### Default Name [%s]\n", route.Method())
	}
	return s
}
