package managers

import (
	"io"
	"log"
	"net/http"
	_ "net/http/httputil" 
	_ "os"
	"regexp"
	"sync"
)

type ImageSearcher struct {
	cfg *Config
}

func GetImageSearcher(cfg *Config) *ImageSearcher {
	return &ImageSearcher{cfg: cfg}
}

type QueryLink struct {
	Query, Link string
}

func (imser *ImageSearcher) GetQueryImageLinks(queries []string) ([]QueryLink, error) {
	resChan := make(chan QueryLink, len(queries))

	var wg sync.WaitGroup
	wg.Add(len(queries))
	for _, query := range queries {
		go func(q string) {
			defer wg.Done()
			imser.getImageLink(q, resChan)
		}(query)
	}
	go func() {
		wg.Wait()
		close(resChan)
	}()

	result := make([]QueryLink, 0, len(queries))

	for queryLink := range resChan {
		result = append(result, queryLink)
	}

	return result, nil
}

var (
	// matching link to the image in position one
	imageLinkRe = regexp.MustCompile(`(?s){.*"position":\s?1,\n.*"original":\s?"(.*)",\n.*}`)
)

func (imser *ImageSearcher) getImageLink(query string, res chan QueryLink) {
	req, err := http.NewRequest("GET", "https://serpapi.com/search", nil)
	if err != nil {
		log.Println(err)
		res <- QueryLink{Query: query, Link: ""}
		return
	}

	q := req.URL.Query()
	q.Add("engine", "google_images_light")
	q.Add("num", "1")
	q.Add("q", query)
	q.Add("api_key", imser.cfg.SerpApiKey)

	req.URL.RawQuery = q.Encode()

	// // Log request dump
	// dump, err := httputil.DumpRequestOut(req, true)
	// if err != nil {
	// 	log.Println("Failed to dump request:", err)
	// } else {
	// 	log.Println("Request dump:\n", string(dump))
	// }

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		res <- QueryLink{Query: query, Link: ""}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		res <- QueryLink{Query: query, Link: ""}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		res <- QueryLink{Query: query, Link: ""}
		return
	}

	match := imageLinkRe.FindStringSubmatch(string(body))
	if len(match) != 2 {
		log.Println("Did not match")
		// f, err := os.Create("log")
		// if err != nil {
		// 	log.Println("Failed to create log file:", err)
		// } else {
		// 	_, err = f.Write(body)
		// 	if err != nil {
		// 		log.Println("Failed to write body to log file:", err)
		// 	}
		// 	f.Close()
		// }
		res <- QueryLink{Query: query, Link: ""}
		return
	}

	link := string(match[1])

	res <- QueryLink{Query: query, Link: link}
}
