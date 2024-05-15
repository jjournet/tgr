package repository

import (
	"fmt"

	"github.com/jjournet/tgr/repository"
)

type RepoElement struct {
	Type       int
	NbElements int
}

func (r *RepoElement) FilterValue() string {
	return repository.ConvertRepoElementType(r.Type)
}

func (r *RepoElement) Title() string {
	return repository.ConvertRepoElementType(r.Type)
}

func (r *RepoElement) Description() string {
	return fmt.Sprintf("%d elements", r.NbElements)
}
