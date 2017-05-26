package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var flags = flagss{}

type flagss struct {
	BaseURL, Username, Password, Action, DataFile string
}

type teamcity struct {
	BaseURL   *url.URL
	CookieJar *cookiejar.Jar
	Client    *http.Client
	DataFile  string
}

func (tc *teamcity) URL(uri string) *url.URL {
	rel, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}
	return tc.BaseURL.ResolveReference(rel)
}

func (tc *teamcity) PostForm(uri string, values url.Values) (*http.Response, error) {
	u := tc.URL(uri)
	log.Printf("POST %s", u)
	return tc.Client.PostForm(u.String(), values)
}

var requiredFields = map[string]*string{}

// req defines a required flag.
func req(ref *string, name, defaultValue, desc string) {
	flag.StringVar(ref, name, defaultValue, desc)
	requiredFields[name] = ref
}

// opt defines an optional flag
func opt(ref *string, name, defaultValue, desc string) {
	flag.StringVar(ref, name, defaultValue, desc)
}

func init() {
	req(&flags.BaseURL, "baseurl", "", "Base URL of TeamCity instance")
	req(&flags.Username, "user", "", "TeamCity username")
	req(&flags.Password, "password", "", "TeamCity password")
	req(&flags.Action, "action", "", "Action to perform")
	opt(&flags.DataFile, "data", "", "File to send as body")
}

func main() {
	flag.Parse()
	for name, ref := range requiredFields {
		if *ref == "" {
			log.Fatalf("required flag %s missing", name)
		}
	}
	action, ok := actions[flags.Action]
	if !ok {
		log.Fatalf("no action named %q", flags.Action)
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	baseURL, err := url.Parse(flags.BaseURL)
	if err != nil {
		log.Fatalf("baseurl %q not valid: %s", flags.BaseURL, err)
	}
	baseURL.User = url.UserPassword(flags.Username, flags.Password)
	tc := teamcity{
		CookieJar: jar,
		Client: &http.Client{
			Jar: jar,
		},
		DataFile: flags.DataFile,
		BaseURL:  baseURL,
	}
	action(tc)
}

var actions = map[string]func(teamcity){
	"set-meta-runner": setMetaRunner,
}

func setMetaRunner(tc teamcity) {
	if tc.DataFile == "" {
		log.Fatal("set-meta-runner requires the -data flag")
	}
	fileName := tc.DataFile + ".xml"
	body, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	response, err := tc.PostForm("plugins/metarunner/runner-edit.html", url.Values{
		"projectId":         []string{"_Root"},
		"editRunnerId":      []string{tc.DataFile},
		"metaRunnerContent": []string{string(body)},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("unable to read response body: %s", err)
	}
	responseBody := string(responseBodyBytes)
	if response.StatusCode != 200 {
		log.Fatalf("Got http status code %d; body: %s", response.StatusCode, string(responseBody))
	}
}
