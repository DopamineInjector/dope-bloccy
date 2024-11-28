package repository

type UserExistsError struct{}

func (t *UserExistsError) Error() string {
	return "user already exists"
}
