package model

func (Users) TableName() string {
	return "users"
}

type Users struct {
	BaseModel
	Username  string     `json:"username" gorm:"size:50;not null"`
	Password  string     `json:"password" gorm:"size:100;not null"`
	Email     string     `json:"email" gorm:"size:100;not null;"`
	UserRoles []UserRole `gorm:"foreignKey:UsersID"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

type UserRole struct {
	BaseModel
	UsersID        uint
	UserRoleTypeID uint
	Users          Users        `gorm:"foreignKey:UsersID"`
	UserRoleType   UserRoleType `gorm:"foreignKey:UserRoleTypeID"`
}

func (UserRoleType) TableName() string {
	return "user_role_type"
}

type UserRoleType struct {
	BaseModel
	Code      string
	UserRoles []UserRole `gorm:"foreignKey:UserRoleTypeID"`
}
