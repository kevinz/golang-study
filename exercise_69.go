package main

import (
  "fmt"
  "time"
)

type Fetcher interface {
  // Fetch returns the body of URL and
  // a slice of URLs found on that page.
  Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(t *crawTask, fetcher Fetcher) {
  if t.depth <= 0 {
    return
  }
  body, urls, err := fetcher.Fetch(t.url)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Printf("found: %s %q\n", t.url, body)
  for _, u := range urls {
    allTasks.add(crawTask{u, t.depth - 1})
  }
}

type crawTask struct {
  url   string
  depth int
}

type CrawTasks struct {
  tasks chan *crawTask
  visit map[string]bool
}

func (ct *CrawTasks) add(t crawTask) {
  if ct.visit[t.url] {
    return
  }
  ct.visit[t.url] = true
  ct.tasks <- &t
}

var allTasks = CrawTasks{make(chan *crawTask, 10), make(map[string]bool)}

func main() {
  task := crawTask{"http://golang.org/", 4}
  allTasks.add(task)
  for {
    select {
    case t := <-allTasks.tasks:
      go Crawl(t, fetcher)
    case <-time.After(11e8):
      return
    }
  }

}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
  body string
  urls []string
}

func (f *fakeFetcher) Fetch(url string) (string, []string, error) {
  if res, ok := (*f)[url]; ok {
    return res.body, res.urls, nil
  }
  return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = &fakeFetcher{
  "http://golang.org/": &fakeResult{
    "The Go Programming Language",
    []string{
      "http://golang.org/pkg/",
      "http://golang.org/cmd/",
    },
  },
  "http://golang.org/pkg/": &fakeResult{
    "Packages",
    []string{
      "http://golang.org/",
      "http://golang.org/cmd/",
      "http://golang.org/pkg/fmt/",
      "http://golang.org/pkg/os/",
    },
  },
  "http://golang.org/pkg/fmt/": &fakeResult{
    "Package fmt",
    []string{
      "http://golang.org/",
      "http://golang.org/pkg/",
    },
  },
  "http://golang.org/pkg/os/": &fakeResult{
    "Package os",
    []string{
      "http://golang.org/",
      "http://golang.org/pkg/",
    },
  },
}
