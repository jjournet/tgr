package ghuser

import (
	"context"

	"github.com/google/go-github/v61/github"
)

type GHUser struct {
	Login      string
	Name       string
	Orgs       []string
	CurrentOrg int64
	Client     *github.Client
}

func NewUser(client *github.Client) *GHUser {
	// get current user login
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		panic(err)
	}
	login := *user.Login
	orgs, _, err := client.Organizations.List(context.Background(), "", nil)
	if err != nil {
		panic(err)
	}
	if len(orgs) == 0 {
		panic("No organizations found")
	}
	orgList := make([]string, len(orgs)+1)
	for i, org := range orgs {
		orgList[i] = *org.Login
	}
	orgList[len(orgs)] = login
	return &GHUser{Login: login, Name: *user.Name, Orgs: orgList, CurrentOrg: int64(len(orgs)), Client: client}
}
