package git

import (
	"log"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	libgit "github.com/libgit2/git2go/v34"
)

func CreateTestRepository(defaultRefName string, contents map[string]string, msg string) (*libgit.Repository, *libgit.Oid, error) {
	repo, err := CreateEmptyTestRepository(defaultRefName)
	if err != nil {
		return nil, nil, err
	}
	oid, err := UpdateAndCommit(repo, contents, msg)
	if err != nil {
		return nil, nil, err
	}
	return repo, oid, nil
}

func CreateEmptyTestRepository(defaultRefName string) (*libgit.Repository, error) {
	tmpdir, err := os.MkdirTemp("", "git-repo*")
	if err != nil {
		return nil, err
	}
	repo, err := libgit.InitRepository(tmpdir, false)
	if err != nil {
		return nil, err
	}

	if defaultRefName == "" {
		defaultRefName = "refs/heads/main"
	}
	if err = repo.SetHead(defaultRefName); err != nil {
		return nil, err
	}

	return repo, nil
}

// UpdateAndCommit accepts a map of filenames to contents and writes it to the repository's index
// this operation preserves unmodified files from the parent commit
func UpdateAndCommit(repo *libgit.Repository, contents map[string]string, msg string) (*libgit.Oid, error) {
	sig := &libgit.Signature{
		Name:  "Styra Dev",
		Email: "dev@styra.com",
		When:  time.Now(),
	}

	idx, err := repo.Index()
	if err != nil {
		return nil, err
	}
	for file, content := range contents {
		filepath := path.Join(repo.Workdir(), file)
		parentDir := filepath[:strings.LastIndex(filepath, "/")]
		if err := os.MkdirAll(parentDir, 0o750); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path.Join(repo.Workdir(), file), []byte(content), 0o600); err != nil {
			return nil, err
		}
		if err := idx.AddByPath(file); err != nil {
			return nil, err
		}
	}

	if err := idx.Write(); err != nil {
		return nil, err
	}

	treeId, err := idx.WriteTree()
	if err != nil {
		return nil, err
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		return nil, err
	}

	var currentTip *libgit.Commit
	if unborn, err := repo.IsHeadUnborn(); err != nil {
		return nil, err
	} else if !unborn {
		ref, err := repo.Head()
		if err != nil {
			return nil, err
		}
		currentTip, err = repo.LookupCommit(ref.Target())
		if err != nil {
			return nil, err
		}
	}
	var commitId *libgit.Oid
	if currentTip != nil {
		commitId, err = repo.CreateCommit("HEAD", sig, sig, msg, tree, currentTip)
	} else {
		commitId, err = repo.CreateCommit("HEAD", sig, sig, msg, tree)
	}
	if err != nil {
		return nil, err
	}

	return commitId, nil
}

func CloneRepository(t *testing.T, url string, checkoutBranch string) *libgit.Repository {
	tmpdir, err := os.MkdirTemp("", "git-clone*")
	if err != nil {
		t.Fatal(err)
	}
	repo, err := libgit.Clone(url, tmpdir, &libgit.CloneOptions{CheckoutBranch: checkoutBranch}) // default options
	if err != nil {
		t.Fatal(err)
	}

	return repo
}

func CleanupTestRepository(r *libgit.Repository) {
	dir := r.Workdir()
	r.Free()
	if err := os.RemoveAll(dir); err != nil {
		log.Fatal(err)
	}
}
