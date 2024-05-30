package users

type Users map[string]*User

func NewUsers() *Users{
	users := make(Users, 0);
	return &users
}

func (users *Users) GetUser(name string) *User{
	return (*users)[name]
}

func (users *Users) AddUser(name string) {
	if _, exist := (*users)[name]; !exist {
		(*users)[name] = NewUser()
	}
}

func (users *Users)GetMinDRFUser() string{
	for user := range *users{
		return user
	}
	return ""
}