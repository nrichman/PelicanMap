package main

import (
	"bufio"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

var wikiURL = "https://stardewvalleywiki.com"

func main() {
	buildVillagerList()
}

func getPageTokenizer(url string) *html.Tokenizer {
	resp, err := http.Get(wikiURL + "/Villagers")
	if err != nil {
		log.Fatal(err)
	}

	return html.NewTokenizer(resp.Body)
}

func buildVillagerList() {
	z := getPageTokenizer(wikiURL + "/Villagers")

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
					writeSliceToFile("test.txt", keysToSlice(neighborSet))
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

func keysToSlice(m map[string]struct{}) []string {
	s := []string{}

	for key := range m {
		s = append(s, key)
	}

	return s
}

func writeSliceToFile(filename string, s []string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(file)

	for _, v := range s {
		_, _ = w.WriteString(v + "\n")
	}

	w.Flush()
	file.Close()
}
