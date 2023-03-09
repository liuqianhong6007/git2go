package git

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestListWorktree(t *testing.T) {
	repo, err := OpenRepository("/Users/liuqianhong/temp/test")
	checkFatal(t, err)

	checkFatal(t, os.MkdirAll(repo.Workdir()+"/_worktree", 0755))

	wtNames := []string{
		"config",
		"language",
		"month_ticket",
		"pipeline",
		"sausage",
		"web-pay-mail",
	}

	for _, wtName := range wtNames {
		wt, err := repo.LookupWorktree(wtName)
		if err != nil {
			if IsErrorClass(err, ErrorClassInvalid) && IsErrorCode(err, ErrorCodeGeneric) {
				reference, err := repo.References.Lookup("refs/heads/" + wtName)
				if err != nil {
					if IsErrorCode(err, ErrorCodeNotFound) {
						remoteBranch, err := repo.LookupBranch("origin/"+wtName, BranchRemote)
						checkFatal(t, err)

						commit, err := repo.LookupCommit(remoteBranch.Target())
						checkFatal(t, err)

						localBranch, err := repo.CreateBranch(wtName, commit, false)
						checkFatal(t, err)

						err = localBranch.SetUpstream("origin/" + wtName)
						checkFatal(t, err)

						reference = localBranch.Reference

					} else {
						t.Fatal(err)
					}
				}
				wtPath := filepath.Join(repo.Workdir(), "/_worktree/"+wtName)
				wt, err = repo.AddWorktree(reference, wtName, wtPath)
				checkFatal(t, err)

			} else {
				t.Fatal(err)
			}
		}

		wtRepo, err := wt.OpenRepository()
		if err != nil {
			if IsErrorClass(err, ErrorClassOS) && IsErrorCode(err, ErrorCodeNotFound) {
				if err = wt.Prune(); err != nil {
					t.Fatal(err)
				}
				continue

			} else {
				t.Fatal(err)
			}
		}
		checkFatal(t, err)

		filename := filepath.Join(wtRepo.Workdir(), "README")
		content := strconv.FormatInt(rand.Int63(), 10)
		err = os.WriteFile(filename, []byte(content), os.ModePerm)
		checkFatal(t, err)
	}

	actualNames, err := repo.ListWorktree()
	checkFatal(t, err)
	t.Logf("valid worktree: %v", actualNames)
}
