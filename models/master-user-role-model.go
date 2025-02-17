package models

type MasterUserRoleStruct struct {
	Admin int
	User  int
}

var MasterUserRoleId = MasterUserRoleStruct{
	Admin: 1,
	User:  2,
}

type MasterUserRoleDbFieldStruct struct {
	Id   string
	Name string
}

var MasterUserRoleDbField = MasterUserRoleDbFieldStruct{
	Id:   "id",
	Name: "name",
}

type MasterUserRole struct {
	Id   int    `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}

func (MasterUserRole) TableName() string {
	return "nexmedis_master_user_roles"
}

type MasterUserRolePrimaryId struct {
	Id int `gorm:"primaryKey"`
}

func (MasterUserRolePrimaryId) TableName() string {
	return "nexmedis_master_user_roles"
}
