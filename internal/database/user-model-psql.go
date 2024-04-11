package database

type UserModel struct {
	db *DataBase
}

func (um *UserModel) CreateNewUser(user *User) (*User, error) {
	if err := um.db.sqlDB.QueryRow(
		"INSERT INTO users (username, email, encrypt_password) VALUES ($2, $3, $4) RETURNING id",
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Encrypt_Password,
	).Scan(&user.ID); err != nil {
		um.db.logger.Errorln(err)
		return nil, err
	}

	um.db.logger.Infoln("User successfully added")
	return user, nil
}

func (um *UserModel) DeleteUser(id int) (*User, error) {
	user := &User{}
	if err := um.db.sqlDB.QueryRow(
		"DELETE FROM users WHERE id = $1 RETURNING id, username, email, encrypt_password", id,
	).Scan(&user.ID, &user.Name, &user.Encrypt_Password); err != nil {
		um.db.logger.Errorln(err)
		return nil, err
	}

	um.db.logger.Infoln("User successfully deleted")
	return user, nil
}

func (um *UserModel) UpdateUserFully(id int, name, email, pass string) (*User, error) {
	user := &User{}
	if err := um.db.sqlDB.QueryRow(
		"UPDATE users SET username = $2, email = $3, encrypt_password = $4 WHERE id = $1 RETURNING id, username, email, encrypt_password", name, email, pass, id,
	).Scan(&user.ID, &user.Name, &user.Encrypt_Password); err != nil {
		um.db.logger.Errorln(err)
		return nil, err
	}

	um.db.logger.Infoln("User successfully updated")
	return user, nil
}

func (um *UserModel) PartUpdateUserName(id int, name string) (*User, error) {
	user := &User{}
	if err := um.db.sqlDB.QueryRow(
		"UPDATE users SET username = $2 WHERE id = $1 RETURNING id, username", id, name,
	).Scan(&user.ID, &user.Name); err != nil {
		um.db.logger.Infoln(err)
		return nil, err
	}

	um.db.logger.Infof("Update the User '%d' with his new name '%s'\n", id, name)
	return user, nil
}

func (um *UserModel) PartUpdateUserEmail(id int, email string) (*User, error) {
	user := &User{}
	if err := um.db.sqlDB.QueryRow(
		"UPDATE users SET email = $3 WHERE id = $1 RETURNING id, email", id, email,
	).Scan(&user.ID, &user.Email); err != nil {
		um.db.logger.Infoln(err)
		return nil, err
	}

	um.db.logger.Infof("Update the User '%d' and new email '%s'\n", id, email)
	return user, nil
}

func (um *UserModel) PartUpdateUserPass(id int, pass string) (*User, error) {
	user := &User{}
	if err := um.db.sqlDB.QueryRow(
		"UPDATE users SET username = $4 WHERE id = $1 RETURNING id, encrypt_password", id, pass,
	).Scan(&user.ID, &user.Encrypt_Password); err != nil {
		um.db.logger.Infoln(err)
		return nil, err
	}

	um.db.logger.Infof("Update the User '%d' with his new password '%s'\n", id, pass)
	return user, nil
}