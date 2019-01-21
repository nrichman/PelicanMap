package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

var wikiURL = "https://stardewvalleywiki.com"

func main() {
	//buildVillagerList()
	buildSchedule()
}

func getPageTokenizer(url string) *html.Tokenizer {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Bad request.")
	}

	return html.NewTokenizer(resp.Body)
}

func buildSchedule() {
	//url := wikiURL + parseFileToSlice("villagerList.txt")[0]
	url := wikiURL + "/Alex"

	z := getPageTokenizer(url)

	headers := map[string][]string{
		"Spring":   []string{},
		"Summer":   []string{},
		"Fall":     []string{},
		"Winter":   []string{},
		"Marriage": []string{},
	}

	header := ""
	collecting := false
	s := []string{}

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			break
		}

		if tt == html.StartTagToken {
			t := z.Token()
			for _, a := range t.Attr {
				if a.Key == "id" && a.Val == "Schedule" {
					collecting = true
				}

				if !collecting {
					continue
				}

				if a.Key == "title" {
					fmt.Println(a.Val)
					if _, ok := headers[a.Val]; ok {
						if header != a.Val {
							header = a.Val
							writeSliceToFile(header+".txt", s)
							s = []string{}
						}
					}
					if header == "Marriage" {
						return
					}
				}
			}

			// Gets the table constraint
			if header != "" && t.Data == "p" {
				z.Next()
				inner := z.Next()
				if inner == html.TextToken {
					s = append(s, (string)(z.Text()))
				}
			}

			// Gets the table values
			if header != "" && t.Data == "td" {
				inner := z.Next()
				if inner == html.TextToken {
					s = append(s, (string)(z.Text()))
				}
			}
		}
	}
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
					writeSliceToFile("villagerList.txt", keysToSlice(neighborSet))
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
		w.WriteString(v + "\n")
	}

	w.Flush()
	file.Close()
}

func parseFileToSlice(filename string) []string {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(file)
	s := []string{}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		s = append(s, line[:len(line)-1])
	}

	return s
}
