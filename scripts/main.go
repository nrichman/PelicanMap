package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/net/html"
)

var wikiURL = "https://stardewvalleywiki.com"

func main() {

	neighbors := retrieveNeighbors()

	var wg sync.WaitGroup
	wg.Add(len(neighbors))

	for _, l := range neighbors {
		go func(link string) {
			defer wg.Done()
			sch := retrieveSchedule(link)
			printMess(link, sch)
		}(l)
	}

	wg.Wait()
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

func printMess(neighbor string, m map[string]map[string][]string) {
	file := "schedules/" + neighbor[1:] + ".txt"
	for key, val := range m {
		//fmt.Println(key)
		writeSliceToFile(file, append([]string{}, key), false, "")
		for key2, val2 := range val {
			//fmt.Println("  " + key2)
			writeSliceToFile(file, append([]string{}, key2), false, "  ")
			for _, val3 := range val2 {
				//fmt.Println("    ", val3)
				writeSliceToFile(file, append([]string{}, val3), false, "    ")
			}
		}
	}
}

func retrieveSchedule(neighbor string) map[string]map[string][]string {
	schedule := map[string]map[string][]string{
		"Spring":     map[string][]string{},
		"Summer":     map[string][]string{},
		"Fall":       map[string][]string{},
		"Winter":     map[string][]string{},
		"Marriage":   map[string][]string{},
		"Deviations": map[string][]string{},
	}

	var season string
	var constraint string
	var time string
	var s []string

	collecting := false
	z := getPageTokenizer(wikiURL + neighbor)
	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			break
		}

		if tt == html.StartTagToken {
			t := z.Token()
			for _, a := range t.Attr {
				if a.Key == "id" && a.Val == "Relationships" {
					if constraint != "" {
						schedule[season][constraint] = s
					}
					return schedule
				}

				// Don't start building the schedule until we've passed the tag
				if a.Key == "id" && a.Val == "Schedule" {
					collecting = true
				}

				if !collecting {
					continue
				}

				// Gets the current season
				if a.Key == "title" {
					if _, ok := schedule[a.Val]; ok {
						// Selected season changes
						if season != a.Val {
							// Fill in the last constraint before changing season
							if constraint != "" {
								schedule[season][constraint] = s
								constraint = ""
							}
							season = a.Val
							s = []string{}
						}
					}
				}
			}

			// Only parse data if a season is selected
			if season != "" {
				switch t.Data {
				// Gets a constrant from <p><b>{CONSTRAINT}</b></p>
				case "p":
					inner := z.Next()
					inner = z.Next()
					if inner == html.TextToken {
						if constraint != "" {
							schedule[season][constraint] = s
							s = []string{}
						}
						constraint = (string)(z.Text())
					}
				// Builds a time/location string from a table
				case "td":
					inner := z.Next()
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
	return nil
}

func retrieveNeighbors() []string {
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
					writeSliceToFile("villagerList.txt", keysToSlice(neighborSet), true, "")
					return keysToSlice(neighborSet)
				}

				if collecting {
					if a.Key == "href" {
						neighborSet[a.Val] = struct{}{}
					}
				}
			}
		}
	}
	return nil
}

func keysToSlice(m map[string]struct{}) []string {
	s := []string{}

	for key := range m {
		s = append(s, key)
	}

	return s
}

func writeSliceToFile(filename string, s []string, create bool, head string) {
	flags := os.O_CREATE | os.O_WRONLY

	if !create {
		flags = flags | os.O_APPEND
	}

	file, err := os.OpenFile(filename, flags, 0644)
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(file)

	for _, v := range s {
		w.WriteString(head + v + "\n")
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
