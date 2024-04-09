package database

type UserModel struct {
	db *DataBase
}

func (um *UserModel) CreateNewUser() (*UserModel, error) {
	um.db.sqlDB.QueryRow("")
	return &UserModel{}, nil
}

func (um *UserModel) PartUpdateUser() (*UserModel, error) {
	um.db.sqlDB.QueryRow("")
	return &UserModel{}, nil
}

func (um *UserModel) UpdateUser() (*UserModel, error) {
	um.db.sqlDB.Query("")
	return &UserModel{}, nil
}

func (um *UserModel) DeleteUser() (*UserModel, error) {
	um.db.sqlDB.Query("")
	return &UserModel{}, nil
}
