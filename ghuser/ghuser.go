package ghuser

import (
	"context"

	"github.com/google/go-github/v69/github"
)

type Owner struct {
	Login       string
	Description string
}

type GHUser struct {
	Login  string
	Name   string
	Orgs   []string
	Owners []Owner

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
	// if len(orgs) == 0 {
	// 	panic("No organizations found")
	// }
	orgList := make([]string, len(orgs)+1)
	owners := make([]Owner, len(orgs)+1)
	for i, org := range orgs {
		orgList[i] = *org.Login
		owners[i] = Owner{Login: *org.Login, Description: *org.Description}
		// log.Printf("Org: %s, Description: %s\n", *org.Login, *org.Description)
	}
	orgList[len(orgs)] = login
	owners[len(orgs)] = Owner{Login: login, Description: "Current User"}
	return &GHUser{Login: login, Name: *user.Name, Orgs: orgList, Owners: owners, CurrentOrg: int64(len(orgs)), Client: client}
}
