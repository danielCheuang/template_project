package model

import "template_project/db/mysql"

type User struct {
	Id int			`json:"id"`
	Name string		`json:"name"`
	Age int			`json:"age"`	
}

func (*User) TableName() string {
	return "user"
}

func (this *User) QueryUserById (id int) ( *User, error) {
	user := User{}
	ret := mysql.DB.Model(&User{}).Where(" id = ? ", id ).Find( &user )
	if ret.Error != nil {
		return nil, ret.Error
	}

	return &user, nil
}





