package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type updates struct {
	latestVersion string
}

type Commit struct {
	Author struct {
		AvatarURL         string `json:"avatar_url"`
		EventsURL         string `json:"events_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		GravatarID        string `json:"gravatar_id"`
		HTMLURL           string `json:"html_url"`
		ID                int    `json:"id"`
		Login             string `json:"login"`
		OrganizationsURL  string `json:"organizations_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		ReposURL          string `json:"repos_url"`
		SiteAdmin         bool   `json:"site_admin"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		URL               string `json:"url"`
	} `json:"author"`
	CommentsURL string `json:"comments_url"`
	Commit      struct {
		Author struct {
			Date  string `json:"date"`
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"author"`
		CommentCount int `json:"comment_count"`
		Committer    struct {
			Date  string `json:"date"`
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			Sha string `json:"sha"`
			URL string `json:"url"`
		} `json:"tree"`
		URL string `json:"url"`
	} `json:"commit"`
	Committer struct {
		AvatarURL         string `json:"avatar_url"`
		EventsURL         string `json:"events_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		GravatarID        string `json:"gravatar_id"`
		HTMLURL           string `json:"html_url"`
		ID                int    `json:"id"`
		Login             string `json:"login"`
		OrganizationsURL  string `json:"organizations_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		ReposURL          string `json:"repos_url"`
		SiteAdmin         bool   `json:"site_admin"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		URL               string `json:"url"`
	} `json:"committer"`
	HTMLURL string `json:"html_url"`
	Parents []struct {
		HTMLURL string `json:"html_url"`
		Sha     string `json:"sha"`
		URL     string `json:"url"`
	} `json:"parents"`
	Sha string `json:"sha"`
	URL string `json:"url"`
}

var (
	errProviderNotSupported = fmt.Errorf("Provider not supported")
)

func getUpdates(g Gom) (*updates, error) {
	if !strings.HasPrefix(g.name, "github.com/") {
		return nil, errProviderNotSupported
	}

	projectName := strings.Replace(g.name, "github.com/", "", -1)
	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/commits", projectName))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	commits := []*Commit{}
	if err := json.Unmarshal(body, &commits); err != nil {
		return nil, err
	}

	return &updates{
		latestVersion: commits[0].Sha,
	}, nil
}

func update(args []string) error {
	allGoms, err := parseGomfile("Gomfile")
	if err != nil {
		return err
	}

	for _, g := range allGoms {
		fmt.Printf("%s\n", g.name)

		commit := g.options["commit"]
		if commit == "" {
			fmt.Printf("  \\_ No commit set. Please set a revion with :commit => 'SHA1'")
			continue
		}

		updates, err := getUpdates(g)
		if err != nil {
			return err
		}

		if err == errProviderNotSupported {
			fmt.Printf("  \\_ Unable to check on this provider. Only github.com is supported\n")
			continue
		}

		if commit == updates.latestVersion {
			fmt.Printf("  \\_ Up to date\n")
		} else {
			fmt.Printf("  \\_ Latest version: %s\n", updates.latestVersion)
		}
	}

	return nil
}
