package main

import (
    "errors"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "flag"
    //"fmt"
)

type AuthConfig struct {
    Locality string
    Username  string
    APIKey string
}

type Container []struct {
    Bytes int64  `json:"bytes"`
    Count int64  `json:"count"`
    Name  string `json:"name"`
}

type ObjectListing []struct {
    Bytes        int64  `json:"bytes"`
    ContentType  string `json:"content_type"`
    Hash         string `json:"hash"`
    LastModified string `json:"last_modified"`
    Name         string `json:"name"`
}

func authenticate(conf AuthConfig) (head http.Header, err error) {

    host := "auth.api.rackspacecloud.com"
    url := "https://auth.api.rackspacecloud.com/v1.0"
    if conf.Locality == "UK" {
        host = "lon.auth.api.rackspacecloud.com"
        url = "https://lon.auth.api.rackspacecloud.com/v1.0"
    }

    head = http.Header{}
    head.Set("Host", host)
    head.Set("X-Auth-User", conf.Username)
    head.Set("X-Auth-Key", conf.APIKey)

    client := http.Client{}

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return
    }
    req.Header = head

    resp, err := client.Do(req)
    if err != nil {
        return
    }

    if resp.StatusCode != 204 {
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }
        err = errors.New(string(body))
        return nil, err
    }

    head = resp.Header
    return head, nil
}

func putFile(filename, token, url string) (*http.Response, error) {
    fh, err := os.Open(filename)
    if err != nil {
        log.Print("PUT failed")
        return nil, err
    }

    url = url + "/test/file.txt" // where /[container]/[filename]
    log.Print("Uploading to: ", url)
    req, err := http.NewRequest("PUT", url, fh)
    if err != nil {
        log.Print("Failed @4")
        return nil, err
    }

    req.Header.Set("Content-Type", "text/plain")
    req.Header.Set("X-Auth-Token", token)
    req.Header.Set("Host", "storage.clouddrive.com")

    return http.DefaultClient.Do(req)
}

func headContainer(container, token, url string) (*http.Response, error) {
    url = url + "/" + container
    log.Print("Heading: ", url)
    req, err := http.NewRequest("HEAD", url, nil)
    if err != nil {
        log.Print("Failed @4")
        return nil, err
    }

    req.Header.Set("Content-Type", "text/plain")
    req.Header.Set("X-Auth-Token", token)
    req.Header.Set("Host", "storage101.dfw1.clouddrive.com")

    return http.DefaultClient.Do(req)
}

func listContainer(container, token, url string) (*http.Response, error) {
    url = url + "/" + container
    log.Print("Heading: ", url)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Print("Container listing failed.")
        return nil, err
    }

    req.Header.Set("Content-Type", "text/plain")
    req.Header.Set("X-Auth-Token", token)
    req.Header.Set("Host", "storage101.dfw1.clouddrive.com")

    return http.DefaultClient.Do(req)
}


var Listit bool

func main() {

    conf := AuthConfig{"", "", "US"}

    flag.StringVar(&conf.Username, "user", "", "username")
    flag.StringVar(&conf.APIKey, "key", "", "api key")
    flag.StringVar(&conf.Locality, "locality", "US", "location key: US|UK")
    //flag.BoolVar(&Listit, "list", false, "GET rather than HEAD")
    flag.Parse()
    /*if Listit != true {
        fmt.Println("not listing it")
    }*/

    auth, err := authenticate(conf)
    if err != nil {
        log.Fatal(err)
    }

    token := auth.Get("X-Storage-Token")
    url := auth.Get("X-Storage-Url")

    resp, err := headContainer("images", token, url)
    if err != nil {
        log.Print("Failed @5")
        log.Fatal(err)
    }

    log.Printf("Got %d status from CloudFiles", resp.StatusCode)
    for key, val := range resp.Header {
        log.Print(key, ": ", val)
    }

    resp, err = listContainer("images?format=json", token, url)
    if err != nil {
        log.Print("Failed @5")
        log.Fatal(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Print("Failed @6")
        log.Fatal(err)
    }

    log.Printf("Got %d status from CloudFiles", resp.StatusCode)
    log.Println("With content: \n" + string(body))
    for key, val := range resp.Header {
        log.Print(key, ": ", val)
    }
}
