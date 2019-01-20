package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

var villagerURL = "https://stardewvalleywiki.com/Villagers"

func main() {
	resp, err := http.Get(villagerURL)
	if err != nil {
		log.Fatal(err)
	}

	z := html.NewTokenizer(resp.Body)

	collecting := false
	neighborSet := make(map[string]struct{})
	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			break
		}

		if tt == html.StartTagToken {
			t := z.Token()

			for _, a := range t.Attr {
				if a.Key == "id" && a.Val == "Bachelors" {
					collecting = true
				}

				if a.Key == "id" && a.Val == "Non-giftable_NPCs" {
					for key := range neighborSet {
						fmt.Println(key)
					}
				}

				if collecting {
					if a.Key == "href" {
						neighborSet[a.Val] = struct{}{}
					}
				}
			}
		}
	}
}
