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

func printMess(m map[string]map[string][]string) {
	for key, val := range m {
		fmt.Println(key)
		for key2, val2 := range val {
			fmt.Println("  " + key2)
			for _, val3 := range val2 {
				fmt.Println("    ", val3)
			}
		}
	}
}

func buildSchedule() {
	//url := wikiURL + parseFileToSlice("villagerList.txt")[0]
	url := wikiURL + "/Alex"

	z := getPageTokenizer(url)

	seasons := map[string]map[string][]string{
		"Spring":   map[string][]string{},
		"Summer":   map[string][]string{},
		"Fall":     map[string][]string{},
		"Winter":   map[string][]string{},
		"Marriage": map[string][]string{},
	}
	season := ""
	constraint := ""
	time := ""
	s := []string{}

	collecting := false

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

				// Gets the current season
				if a.Key == "title" {
					if _, ok := seasons[a.Val]; ok {
						// Selected season changes
						if season != a.Val {
							// Fill in the last constraint of the season
							if constraint != "" {
								seasons[season][constraint] = s
							}
							season = a.Val
							s = []string{}
						}
					}

					// Exit condition
					if season == "Marriage" {
						printMess(seasons)
						return
					}
				}
			}

			// Only parse data if a season is selected
			if season != "" {
				inner := z.Next()
				switch t.Data {
				// Gets
				case "p":
					inner = z.Next()
					if inner == html.TextToken {
						if constraint != "" {
							seasons[season][constraint] = s
							s = []string{}
						}
						constraint = (string)(z.Text())
					}
				// Builds a time/location string from a table
				case "td":
					if inner == html.TextToken {
						token := (string)(z.Text())
						token = token[:len(token)-1]
						if token == "" {
							continue
						}

						if time == "" {
							time = token
						} else {
							location := token
							s = append(s, time+";"+location)
							time = ""
						}
					}
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
