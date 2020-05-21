package model

func migration() {
	DB.AutoMigrate(
		&User{},
		&Role{},
		&Menu{},
		&RoleMenu{},
		&Server{},
		&ServerCollaborator{},
	)
}
