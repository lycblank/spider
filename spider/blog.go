package spider

import (
	"bytes"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
)

type BlogSpider struct {
}

func NewBlogSpider() ISpider {
	return &BlogSpider{}
}

func (bs *BlogSpider) Spider(url string, header http.Header, next SpiderFunc) error {
	query := gorequest.New().Get(url)
	for k, v := range header {
		if len(v) > 0 {
			query = query.Set(k, v[0])
		}
	}
	if _, body, errs := query.EndBytes(); len(errs) > 0 {
		return errs[0]
	} else if doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body)); err != nil {
		return err
	} else {
		return next(doc)
	}
	return nil
}
