package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

func main() {
	fmt.Println("starting crawler...")
	srcUrl := flag.String("url", "", "web url to crawl")
	dir := flag.String("dir", "", "web url to crawl is required")
	flag.Parse()
	crawlWebpage(srcUrl, dir)
}

func crawlWebpage(srcUrl *string, dir *string) {
	parsedUrl, err := url.Parse(*srcUrl)
	if err != nil {
		fmt.Printf("%s", "Unable to parse url:"+*srcUrl)
		log.Fatal(err)
	}
	downloadedFile := downloadWebpage(srcUrl, dir)
	document, err := goquery.NewDocumentFromReader(downloadedFile)
	if err != nil {
		fmt.Printf("%s", "Unable to create query document from file:"+downloadedFile.Name())
		log.Fatal(err)
	}
	document.Find("a").Each(func(i int, selection *goquery.Selection) {
		hrefValue, isFound := selection.Attr("href")
		if isFound {
			anchorUrl, err := url.Parse(hrefValue)
			if err != nil {
				fmt.Printf("%s", "Unable to create query document from file:"+downloadedFile.Name())
				log.Fatal(err)
			}
			if parsedUrl.Host == anchorUrl.Host && parsedUrl.Scheme == anchorUrl.Scheme {
				c := make(chan os.Signal)
				signal.Notify(c, os.Interrupt, syscall.SIGTERM)
				go func() {
					sig := <-c
					crawlWebpage(&hrefValue, dir)
					fmt.Printf("Exiting webcrawler")
					fmt.Println(sig)
					os.Exit(1)
				}()
			}
		}
	})
}

func downloadWebpage(url *string, dir *string) *os.File {
	fmt.Println("Downloading web page..")
	urlRegex := "https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"
	match, err := regexp.MatchString(urlRegex, *url)
	if err != nil || !match {
		fmt.Printf("%s", "invalid url provided")
		return nil
	}
	resp, err := http.Get(*url)
	if err != nil {
		fmt.Printf("%s", "error occured reaching url:"+*url)
		log.Fatal(err)
	}
	defer resp.Body.Close()
	err = os.MkdirAll(*dir, os.ModeDir)
	if err != nil {
		fmt.Printf("%s", "unable to create directory:"+*dir)
		log.Fatal(err)
	}
	rand.Seed(int64(time.Now().Nanosecond()))
	filePrefix := fmt.Sprint(rand.Int())
	file, err := os.Create(*dir + "/" + filePrefix + ".html")
	if err != nil {
		fmt.Printf("%s", "unable to create file:"+*dir+"/1.html")
		log.Fatal(err)
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("%s", "unable to copy web content to file:"+file.Name())
		log.Fatal(err)
	}
	return file
}
