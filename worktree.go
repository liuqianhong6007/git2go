package git

/*
#include <git2.h>
*/
import "C"
import "runtime"

type Worktree struct {
	doNotCompare
	cast_ptr *C.git_worktree
}

func newWorktree(wt *C.git_worktree) *Worktree {
	return &Worktree{
		cast_ptr: wt,
	}
}

func (wt *Worktree) Name() string {
	cname := C.git_worktree_name(wt.cast_ptr)
	return C.GoString(cname)
}

func (wt *Worktree) Path() string {
	cpath := C.git_worktree_path(wt.cast_ptr)
	return C.GoString(cpath)
}

func (wt *Worktree) OpenRepository() (*Repository, error) {
	var repo *C.git_repository

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_repository_open_from_worktree(&repo, wt.cast_ptr)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}
	return newRepositoryFromC(repo), nil
}

func (wt *Worktree) Prune() error {
	var opt C.git_worktree_prune_options

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_worktree_prune_options_init(&opt, C.GIT_WORKTREE_PRUNE_OPTIONS_VERSION)
	if ret < 0 {
		return MakeGitError(ret)
	}

	ret = C.git_worktree_is_prunable(wt.cast_ptr, &opt)
	if ret != 0 {
		return MakeGitError(ret)
	}

	ret = C.git_worktree_prune(wt.cast_ptr, &opt)
	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}
