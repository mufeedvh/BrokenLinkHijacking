package main

import (
  "crypto/tls"
  "flag"
  "fmt"
  "github.com/jackdanger/collectlinks"
  "net/http"
  "net/url"
  "strings"
  "os"
)


var visited = make(map[string]bool)






func main() {
  flag.Parse()

  args := flag.Args()
  fmt.Println(args)
  if len(args) < 1 {
    fmt.Println("specify the target")
    os.Exit(1)
  }

  queue := make(chan string)

  go func() { queue <- args[0] }()

  for uri := range queue {
    enqueue(uri, queue)
  }
}

func enqueue(uri string, queue chan string) {

  visited[uri] = true
  transport := &http.Transport{
    TLSClientConfig: &tls.Config{
      InsecureSkipVerify: true,
    },
  }
  client := http.Client{Transport: transport}
  resp, err := client.Get(uri)
  if err != nil {
    return
  }
  defer resp.Body.Close()

  links := collectlinks.All(resp.Body)
  status_codes:=map[int]string{
    404:"Resource Not FOUND",
  }
  domain_list:=map[string]string{
    "linkedin.com":"Linkedin",
    "facebook.com":"Facebook",
    "twittter.com":"twitter",
    "youtube.com":"youtube",
    "twitch.com":"twitch",
    "discord.com":"discord",
             }
  for _, link := range links {

    absolute := fixUrl(link, uri)
    if uri != "" {
      if !visited[absolute] {
        response, err := client.Get(absolute)
       if err!=nil{
         return
       }
        u, err := url.Parse(absolute)
             if err != nil {
                 panic(err)
             }
             parts := strings.Split(u.Hostname(), ".")
             domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
            _, exists := status_codes[response.StatusCode]

        _,domain_exists:= domain_list[domain]
        if exists && domain_exists{

            fmt.Println("Seems to be vulnerable",link)
        }
          fmt.Println(absolute)
        go func() { queue <- absolute }()
      }
    }
  }
}

func fixUrl(href, base string) (string) {
  uri, err := url.Parse(href)
  if err != nil {
    return ""
  }
  baseUrl, err := url.Parse(base)
  if err != nil {
    return ""
  }
  uri = baseUrl.ResolveReference(uri)
  return uri.String()
}