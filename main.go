package main

import (
    "fmt"
    "log"

    "flag"

    "bytes"
    "strings"

    "os/exec"

    "net/http"

    "encoding/json"
    "encoding/base64"

    "html/template"

    "strconv"

    "github.com/andreadipersio/pugo/lib/cache"
)

const (
    githubContentURL = "https://api.github.com/repos/%v/%v/contents/%v"
)

var (
    httpPort int

    repo     string
    owner    string
    token    string

    staticsPath  string
    templatesPath string

)

type GithubResponse struct {
    Content string `json:"content"`
}

func getTemplatePath(filename string) string {
    return fmt.Sprintf("%v/%v", templatesPath, filename)
}

func init() {
    flag.IntVar(&httpPort, "port", 8002, "HTTP Server port")

    flag.StringVar(&repo, "repo", "", "Github public repository name")
    flag.StringVar(&owner, "owner", "", "Github public repository owner")
    flag.StringVar(&staticsPath, "statics", "./statics", "path to static files")
    flag.StringVar(&templatesPath, "templates", "./templates", "path to html templates")
    flag.StringVar(&token, "token", "", "Github Personal Token")

    flag.Parse()
}

func getFilename(urlPath string) string {
    parts := strings.Split(urlPath, "/")
    filename := parts[len(parts) - 1]

    if filename == "" {
        return "README.md"
    }

    return fmt.Sprintf("%v.md", filename)
}

func rootHandler(w http.ResponseWriter, r *http.Request, localCache *cache.Cache) {
    var content string

    filename := getFilename(r.URL.Path)

    log.Printf("[REQ] %v", filename)

    refreshFlag := r.FormValue("refresh")

    getFromCache := func() string {
        localCache.GetChan <- filename
        return <-localCache.GetChan
    }

    getFromGithub := func() (content string) {
        githubResp := &GithubResponse{}

        fetch := func() error {
            remoteURL := fmt.Sprintf(githubContentURL, owner, repo, filename)

            log.Println(remoteURL)

            client := http.Client{}
            req, err := http.NewRequest("GET", remoteURL, nil)

            r.Header.Set("User-Agent", "andreadipersio/pugo")

            if token != "" {
                r.SetBasicAuth(owner, token)
            }

            if err != nil {
                panic(err)
            }

            resp, err := client.Do(req)

            if err != nil {
                panic(err)
            }

            defer resp.Body.Close()

            dec := json.NewDecoder(resp.Body)

            if err := dec.Decode(githubResp); err != nil {
                return err
            }

            return nil
        }

        decode := func() error {
            b64 := base64.StdEncoding

            if v, err := b64.DecodeString(githubResp.Content); err != nil {
                return err
            } else {
                content = string(v)
            }

            return nil
        }

        toHTML := func() error {
            cmd := exec.Command("markdown")

            if pipeWriter, err := cmd.StdinPipe(); err != nil {
                return err
            } else {
                pipeWriter.Write([]byte(content))
                pipeWriter.Close()
            }

            if output, err := cmd.Output(); err != nil {
                return err
            } else {
                content = string(output)
            }

            if tmpl, err := template.ParseFiles(getTemplatePath("base.html"), getTemplatePath("article.html")); err != nil {
                return err
            } else {
                context := &struct { Content template.HTML } { template.HTML(content), }

                wr := bytes.NewBufferString("")

                if parseErr := tmpl.Execute(wr, context); parseErr != nil {
                    return parseErr
                }

                content = wr.String()
            }

            localCache.PutChan <-cache.CacheRequest{filename, content}

            return nil
        }

        if err := fetch(); err != nil {
            log.Panicf("Cannot get data from github: %v", err)
        }

        if err := decode(); err != nil {
            log.Panicf("Cannot decode github response: %v", err)
        }

        if err := toHTML(); err != nil {
            log.Panicf("Cannot convert markdown to html: %v", err)
        }

        log.Println("[Github] OK")

        return
    }

    if content = getFromCache(); content == "" || refreshFlag != "" {
        content = getFromGithub()
    }

    w.Header().Set("Content-Type", "text/html")
    w.Header().Set("Content-Length", strconv.Itoa(len(content)))

    fmt.Fprint(w, content)
}

func main() {
    httpAddr := fmt.Sprintf(":%v", httpPort)

    log.Printf("Serving static files from %v", staticsPath)

    log.Printf("Starting server on %v", httpAddr)

    localCache := cache.NewCache()
    go localCache.Run()

    http.HandleFunc("/favicon.ico", http.NotFound)
    http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir(staticsPath))))

    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        rootHandler(w, r, localCache)
    })

    log.Fatal(http.ListenAndServe(httpAddr, nil))
}
