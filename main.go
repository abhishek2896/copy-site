package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Metadata struct {
	Site     string
	NumLinks int
	NumImg   int
	Time     time.Time
}

func removeIP(url string, isFile bool) string {
	url = strings.ReplaceAll(url, "http://", "")
	url = strings.ReplaceAll(url, "https://", "")
	if isFile {
		url += ".html"
	}
	return url
}

func getURLContent(args []string, clone bool) error {
	var returningErr error = nil
	for i := 0; i < len(args); i++ {
		resp, err := http.Get(args[i])

		if err != nil {
			returningErr = err
			break
		}

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			returningErr = err
			break
		}

		if clone {
			err = os.WriteFile(removeIP(args[i], false)+"/index.html", body, 0644)
		} else {
			err = os.WriteFile(removeIP(args[i], true), body, 0644)
		}

		if err != nil {
			returningErr = err
			break
		}
	}
	return returningErr
}

func getURLMetadata(args []string) error {
	var returningErr error = nil
	var metadata []Metadata = []Metadata{}
	for i := 0; i < len(args); i++ {
		c := colly.NewCollector()
		imageCount := 0
		linkCount := 0

		c.OnHTML("img[src]", func(e *colly.HTMLElement) {
			imageCount++
		})

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			linkCount++
		})

		c.OnRequest(func(r *colly.Request) {
			// do nothing
		})

		returningErr = c.Visit(args[i])

		if returningErr != nil {
			break
		}

		siteData := Metadata{
			Site:     args[i],
			NumLinks: linkCount,
			NumImg:   imageCount,
			Time:     time.Now(),
		}

		metadata = append(metadata, siteData)
	}
	if returningErr == nil {
		for i := 0; i < len(metadata); i++ {
			fmt.Print("\n")
			fmt.Println("site:", metadata[i].Site)
			fmt.Println("num_links:", metadata[i].NumLinks)
			fmt.Println("images:", metadata[i].NumImg)
			fmt.Println("last_fetch:", metadata[i].Time)
		}
	}
	return returningErr
}

func downloadImgToPath(fileName string, url string, path string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	//open a file for writing
	file, err := os.Create(path + "/" + fileName)

	if err != nil {
		return err
	}

	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func downloadFileToPath(fileName string, url string, path string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/"+fileName, body, 0644)

	if err != nil {
		return err
	}

	return nil
}

func isValidLink(urlStr string) bool {
	if _, err := http.Get(urlStr); err != nil {
		return false
	}
	return true
}

func getClone(args []string) error {
	var returningErr error = nil
	for i := 0; i < len(args); i++ {
		c := colly.NewCollector()

		c.OnHTML("img[src]", func(e *colly.HTMLElement) {
			link := e.Attr("src")
			if strings.HasPrefix(link, "data:image") || strings.HasPrefix(link, "blob:") {
				return
			}

			if !isValidLink(link) {
				links := strings.Split(link, "/")
				fileName := links[len(links)-1]
				links = links[:len(links)-2]
				if err := os.MkdirAll(removeIP(args[i], false)+strings.Join(links, "/"), os.ModePerm); err != nil {
					log.Fatal(err)
				}
				downloadImgToPath(fileName, args[i]+link, removeIP(args[i], false)+strings.Join(links, "/"))
			}
		})

		c.OnHTML("script[src]", func(e *colly.HTMLElement) {
			link := e.Attr("src")

			if !isValidLink(link) {
				link = removeIP(args[i], false)
				links := strings.Split(link, "/")
				fileName := links[len(links)-1]
				links = links[:len(links)-2]
				if err := os.MkdirAll(removeIP(args[i], false)+strings.Join(links, "/"), os.ModePerm); err != nil {
					log.Fatal(err)
				}
				downloadFileToPath(fileName, args[i]+link, removeIP(args[i], false)+strings.Join(links, "/"))
			}
		})

		c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
			link := e.Attr("href")

			if !isValidLink(link) {
				link = removeIP(args[i], false)
				links := strings.Split(link, "/")
				fileName := links[len(links)-1]
				links = links[:len(links)-2]
				if err := os.MkdirAll(removeIP(args[i], false)+strings.Join(links, "/"), os.ModePerm); err != nil {
					log.Fatal(err)
				}
				downloadFileToPath(fileName, args[i]+link, removeIP(args[i], false)+strings.Join(links, "/"))
			}
		})

		c.OnRequest(func(r *colly.Request) {
			if err := os.Mkdir(removeIP(args[i], false), os.ModePerm); err != nil && !strings.Contains(err.Error(), "file exists") {
				log.Fatal(err)
			}
		})

		c.Visit(args[i])

		getURLContent([]string{args[i]}, true)

	}
	return returningErr
}

func main() {
	args := os.Args[1:]
	var err error
	if args[0] == "--metadata" {
		err = getURLMetadata(args[1:])
	} else if args[0] == "--clone" {
		// incomplete
		err = getClone(args[1:])
	} else {
		err = getURLContent(args, false)
	}

	if err != nil {
		fmt.Println(err)
	}
}
