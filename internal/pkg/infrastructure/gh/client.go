package gh

import (
	"strings"

	gogh "github.com/cli/go-gh"
)

func Execute(args []string) (string, error) {
	result, _, err := gogh.Exec(args...)
	return result.String(), err
}

func RetrieveRepoInformation(repo string) (*RepoInfo, error) {
	if repo != "" {
		parts := strings.Split(repo, "/")
		if len(parts) == 2 {
			return &RepoInfo{Owner: parts[0], Name: parts[1]}, nil
		}
	}
	
	currentRepo, err := gogh.CurrentRepository()
	if err != nil {
		return nil, err
	}
	
	return &RepoInfo{
		Owner: currentRepo.Owner(),
		Name:  currentRepo.Name(),
	}, nil
}

type RepoInfo struct {
	Owner string
	Name  string
}







