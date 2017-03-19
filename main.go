package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/dyatlov/go-opengraph/opengraph"
)

func main() {
	deps, err := GetAllDeps()

	parseError(err)

	for _, dep := range deps {
		_, err := url.ParseRequestURI(dep.HomePage)
		if err != nil {
			continue
		}
		resp, err := http.Get(dep.HomePage)
		if err != nil {
			continue
		}
		og := opengraph.NewOpenGraph()
		err = og.ProcessHTML(resp.Body)
		if err != nil {
			continue
		}
		if len(og.Images) > 0 {
			//if not valid url, try to add homepage and see if valid. else don't care
			if og.Images[0].SecureURL == "" {
				log.Println(dep.FullName + ": " + og.Images[0].URL)
			} else {
				log.Println(dep.FullName + ": " + og.Images[0].SecureURL)
			}
		}
	}
}
