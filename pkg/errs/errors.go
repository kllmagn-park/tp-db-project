package errs

import "errors"

var (
	ErrorForumAlreadyExist = errors.New("forum already exist")
	ErrorForumDoesNotExist = errors.New("forum does not exist")

	ErrorPostDoesNotExist       = errors.New("post does not exist")
	ErrorAuthorDoesNotExist     = errors.New("author does not exist")
	ErrorParentPostDoesNotExist = errors.New("parent post does not exist")

	ErrorNoAuthorOrForum    = errors.New("author or forum does not exist")
	ErrorThreadAlreadyExist = errors.New("thread already exist")
	ErrorThreadDoesNotExist = errors.New("thread does not exist")

	ErrorUserAlreadyExist   = errors.New("user already exist")
	ErrorUserDoesNotExist   = errors.New("user does not exist")
	ErrorConflictUpdateUser = errors.New("data conflicts with existing users")
)