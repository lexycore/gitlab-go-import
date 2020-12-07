package impl

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/xanzy/go-gitlab"
)

var requests sync.Map

func getRequest(req string) (string, bool) {
	v, ok := requests.Load(req)
	if !ok {
		return "", false
	}
	res, ok2 := v.(string)
	return res, ok2
}

func setRequest(req, res string) {
	requests.Store(req, res)
}

func (server *Server) goGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	goGetValue, ok := query["go-get"]
	if !ok || len(goGetValue) < 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	path := strings.ReplaceAll(r.URL.Path, "%2F", "/")
	path = strings.ReplaceAll(path, "%2f", "/")
	path = strings.Trim(path, "/")
	log.Printf("url: %s", path)

	if response, ok2 := getRequest(path); ok2 {
		server.sendSuccessResult(w, response)
		return
	}

	pathItems := strings.Split(path, "/")

	if len(pathItems) < 1 {
		log.Printf("Wrong parameters: empty `url`")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i := len(pathItems); i > 0; i-- {
		projectPath := strings.Join(pathItems[0:i], "/")
		project, _, err := server.git.Projects.GetProject(projectPath, nil)
		if err == nil && project != nil {
			result := server.successResult(w, project)
			setRequest(path, result)
			return
		}
	}

	server.fallbackRequest(w, r.RequestURI)
}

func (server *Server) successResult(w http.ResponseWriter, project *gitlab.Project) string {
	webURL := project.WebURL
	webURLObj, err := url.Parse(project.WebURL)
	if err == nil {
		webURL = fmt.Sprintf("%s%s", webURLObj.Host, webURLObj.Path)
	}
	response := fmt.Sprintf(
		"<html><head><meta name=\"go-import\" content=\"%s git ssh://%s\"/></head></html>\n",
		webURL,
		strings.Replace(project.SSHURLToRepo, ":", "/", 1),
	)
	server.sendSuccessResult(w, response)
	return response
}

func (server *Server) sendSuccessResult(w http.ResponseWriter, response string) {
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Printf("Error writing response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *Server) fallbackRequest(w http.ResponseWriter, path string) {
	fallbackURL := fmt.Sprintf("%s%s?go-get=1", server.config.GitLabURL, path)
	log.Printf("sending fallback request to: %s", fallbackURL)
	resp, err := http.Get(fallbackURL)
	if err != nil {
		log.Printf("Error requesting GitLab fallback: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		log.Printf("Error reading GitLab response: %s", errRead)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		log.Printf("Error writing response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
