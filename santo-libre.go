package libre

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

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

func routes_to_md(routes martini.Routes) string {
	s := ""
	s += "FORMAT: 1A\nHOST: http://made.up\n"
	paths := make(map[string][]martini.Route)
	for _, route := range routes.All() {
		if _, ok := paths[route.Pattern()]; ok {
			paths[route.Pattern()] = make([]martini.Route, 1)
		}
		paths[route.Pattern()] = append(paths[route.Pattern()], route)
	}
	for pattern, routes := range paths {
		s += fmt.Sprintf("## Default Name [%s]\n", pattern)
		for _, route := range routes {
			s += fmt.Sprintf("\n### Default Name [%s]\n", route.Method())
		}
	}
	return s
}
