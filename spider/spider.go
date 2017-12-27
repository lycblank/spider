package spider

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type SpiderFunc func(doc *goquery.Document) error

type ISpider interface {
	Spider(url string, header http.Header, next SpiderFunc) error
}
