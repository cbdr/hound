package main

import (
	"context"
	"flag"
	"strconv"

	"github.com/google/go-github/v39/github"
	"github.com/hound-search/hound/configurator"
)

func main() {

	username := flag.String("username", "", "")
	password := flag.String("password", "", "")
	output := flag.String("outputfile", "config.json", "")

	flag.Parse()

	if flag.NArg() != 0 {
		flag.Usage()
		return
	}

	configFile := configurator.Config{
		MaxConcurrentIndexers: 2,
		DbPath:                "data",
		Title:                 "Hound",
		HealthCheckURI:        "/healthz",
		VcsConfg:              make(map[string]*configurator.Vcs),
		Repos:                 make(map[string]*configurator.Repo),
	}

	var vcs = configurator.Vcs{
		DetectRef: true,
	}
	configFile.VcsConfg["git"] = &vcs

	ctx := context.Background()
	var tp github.BasicAuthTransport
	tp.Username = *username
	tp.Password = *password
	client := github.NewClient(tp.Client())
	opt := &github.RepositoryListByOrgOptions{Sort: "full_name", ListOptions: github.ListOptions{PerPage: 50}}

	var count int = 0
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, "cbdr", opt)
		if err != nil {
			println(err.Error())
			return
		}

		for _, s := range repos {
			if !*s.Archived {
				var repository = configurator.Repo{
					Url:             *s.SSHURL,
					MsBetweenPolls:  3600000,
					ExcludeDotFiles: true,
				}
				count = count + 1
				println(*s.Name + " (" + strconv.Itoa(count) + ")")
				configFile.Repos[*s.Name] = &repository
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	configFile.SaveToFile(*output)
}
