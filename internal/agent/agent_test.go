package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"newsplatform/internal/models"
)

var fakeResponse = `
{
	"status": "ok",
	"totalResults": 12345,
	"articles": [{
		"source": {
			"id": "aaa",
			"name": "bbb"
		},
		"author": "fake name",
		"title": "This is title",
		"description": "This is description",
		"url": "https://www.abc/edf",
		"urlToImage": "https://abc/def/g.jpg",
		"content": "content... here[+4140 chars]"
	}]
}
`

const htmlPage = `<!DOCTYPE html>
<html>
<head>
    <title>title</title>
</head>
<body>
    <h1>h1 text</h1>
    <div>
		<p>p1</p>
		<p>p2</p>
	</div>
    <em>em text</em>
	<p>p outer</p>
    <footer>
		footer text
		<p>p in footer</p>
	</footer>
    copyright
    <p></p>
</body>
</html>`

func setupSrv(fakeResponse string) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(fakeResponse))
	}))
	return srv
}

// TODO: add more test cases
func Test_request(t *testing.T) {
	srv := setupSrv(fakeResponse)
	defer srv.Close()
	articles, err := request(srv.URL)
	want := []models.News{models.News{ID: 0, Title: "This is title", Author: "fake name", Description: "This is description", URL: "https://www.abc/edf", Content: "content... here[+4140 chars]"}}
	assert.NoError(t, err)
	assert.Equal(t, want, articles)
}

func Test_getFullContent(t *testing.T) {
	srv := setupSrv(htmlPage)
	defer srv.Close()
	content, err := getFullContent(srv.URL)
	assert.NoError(t, err)
	assert.Equal(t, "title\n\nh1 text\n\np1\n\np2\n\np outer\n\np in footer\n\n\n\n", content)
}
