package exception

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidAccessToken = errors.New("invalid access token")
)

func ServiceException(url string) error {
	return errors.New("Response Error from: " + url)
}
