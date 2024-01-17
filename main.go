package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	ccsv "github.com/tsak/concurrent-csv-writer"
)

type extractedJob struct {
	id       string
	title    string
	location string
	summary  string
}

var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=python"

func main() {
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJob := <-c
		jobs = append(jobs, extractedJob...)
	}

	writeJobs(jobs)

	fmt.Println("Done, extracted", len(jobs))
}

func getPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL + "&recruitPage=" + strconv.Itoa(page+1)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".item_recruit")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extracteJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func extracteJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("value")
	title := cleanString(card.Find(".area_job>h2>a").Text())
	location := cleanString(card.Find(".area_job>.job_condition>span>a").Text())
	summary := cleanString(card.Find(".area_job>.job_sector").Text())
	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		summary:  summary,
	}
}

func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func writeJobs(jobs []extractedJob) {
	csv, err := ccsv.NewCsvWriter("Jobs.csv")
	checkErr(err)

	defer csv.Close()
	headers := []string{"Link", "Title", "Location", "Summary"}
	wErr := csv.Write(headers)
	checkErr(wErr)

	done := make(chan bool)
	for _, job := range jobs {
		go func(job extractedJob) {
			jobSlice := []string{"https://www.saramin.co.kr/zf_user/jobs/relay/view?isMypage=no&rec_idx=" + job.id, job.title, job.location, job.summary}
			jwErr := csv.Write(jobSlice)
			checkErr(jwErr)
			done <- true
		}(job)
	}

	for i := 0; i < len(jobs); i++ {
		<-done
	}

	// file, err := os.Create("Jobs.csv")
	// checkErr(err)

	// w := csv.NewWriter(file)
	// defer w.Flush()

	// headers := []string{"Link", "Title", "Location", "Summary"}

	// wErr := w.Write(headers)
	// checkErr(wErr)

	// for _, job := range jobs {
	// 	jobSlice := []string{"https://www.saramin.co.kr/zf_user/jobs/relay/view?isMypage=no&rec_idx=" + job.id, job.title, job.location, job.summary}
	// 	jwErr := w.Write(jobSlice)
	// 	checkErr(jwErr)
	// }
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
