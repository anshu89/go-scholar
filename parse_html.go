package main

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
	"log"
	"fmt"
	"encoding/json"
)

type Article struct {
	Title             string
	Year              string
	URL               string
	CitationURL       string
	ClusterId         string
	NumberOfCitations string
	NumberOfVersions  string
	InfoId            string
	PDFLink           string
	PDFSource         string
	Bibtex            string
}


type Articles struct {
	n        int
	articles []Article
}


func (a *Article) Parse(s *goquery.Selection) {
	// title
	h3Title := s.Find(ARTICLE_TITLE_SELECTOR)
	a.URL, _ = h3Title.Attr("href") // TODO: doesn't work
	a.Title = h3Title.Text()

	// header
	a.Year = s.Find(ARTICLE_HEADER_SELECTOR).Text() // TODO: parse

	// footer
	divFooter := s.Find(ARTICLE_FOOTER_SELECTOR)
	parseFooter := func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		text := s.Text()
		if strings.HasPrefix(href, "/scholar?cites") {
			a.CitationURL = href
			a.NumberOfCitations = text
		}
		if strings.HasPrefix(href, "/scholar?cluster") {
			a.ClusterId = href
			a.NumberOfVersions = text
		}
		if strings.HasPrefix(href, "/scholar?q=related") {
			a.InfoId = href
		}

	}
	divFooter.Find("a").Each(parseFooter)

	// sideBar
	sideBarA := s.Find(ARTICLE_SIDEBAR_SELECTOR)
	a.PDFLink, _ = sideBarA.Attr("href")
	a.PDFSource = sideBarA.Text()

	a.parseBibTeX()
}

func (a *Article) parseBibTeX() {
	popURL := "https://scholar.google.co.jp/scholar?q=info:" + a.InfoId + ":scholar.google.com/&output=cite&scirp=0&hl=en"
	popDoc, err := goquery.NewDocument(popURL)
	if err != nil {
		log.Fatal(err)
	}
	bibURL, _ := popDoc.Find("#gs_citi > a:first-child").Attr("href")
	bibDoc, err := goquery.NewDocument("https://scholar.google.co.jp/" + bibURL)
	if err != nil {
		log.Fatal(err)
	}
	a.Bibtex = bibDoc.Text()
}

func (a *Article) dump() {
	fmt.Println("title :", a.Title)
	fmt.Println("year :", a.Year)
	fmt.Println("url: ", a.URL)
	fmt.Println("cluster_id: ", a.ClusterId)
	fmt.Println("# of citations: ", a.NumberOfVersions)
	fmt.Println("# of versions: ", a.NumberOfCitations)
	fmt.Println("infor id: ", a.InfoId)
	fmt.Println("pdfLink: ", a.PDFLink)
	fmt.Println("pdfSource: ", a.PDFSource)
	fmt.Println("citationURL: ", a.CitationURL)
	fmt.Println("BibTeX: ", a.Bibtex)
}


func (as *Articles) ParseAllArticles(doc *goquery.Document) {
	as.articles = make([]Article, as.n)

	parse := func(i int, s *goquery.Selection) {
		as.articles[i].Parse(s)
	}
	doc.Find(WHOLE_ARTICLE_SELECTOR).Each(parse)
}

func (as *Articles) Json() string {
	bytes, err := json.Marshal(as.articles)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}
