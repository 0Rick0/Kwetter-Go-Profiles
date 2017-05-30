package types

type User struct {
	Id int		`gorm:"primary_key"`
	Username string	`sql:"unique"`
	MimeType string
}

func (User) SwaggerDoc() map[string]string {
	return map[string]string{
		"":         "A user object",
		"id":	    "The id of the user",
		"username": "The username of the user",
	}
}

