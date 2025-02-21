package database

import (
	"github.com/Mopsgamer/draqun/server/controller/model_database"
)

// Create new DB record.
func (db Database) UserCreate(user model_database.User) *uint64 {
	query :=
		`INSERT INTO app_users (
			nickname,
			username,
			email,
			phone,
			password,
			avatar,
			created_at,
			last_seen
		)
    	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Sql.Exec(query,
		user.Nick,
		user.Name,
		user.Email,
		user.Phone,
		user.Password,
		user.Avatar,
		user.CreatedAt,
		user.LastSeen,
	)

	if err != nil {
		logSqlError(err)
		return nil
	}

	newId := &db.Context().LastInsertId
	return newId
}

// Change the existing DB record.
func (db Database) UserUpdate(user model_database.User) bool {
	query :=
		`UPDATE app_users
    	SET
		nickname = ?,
		username = ?,
		email = ?,
		phone = ?,
		password = ?,
		avatar = ?,
		created_at = ?,
		last_seen = ?

        WHERE id = ?`
	_, err := db.Sql.Exec(query,
		user.Nick,
		user.Name,
		user.Email,
		user.Phone,
		user.Password,
		user.Avatar,
		user.CreatedAt,
		user.LastSeen,
		user.Id,
	)

	if err != nil {
		logSqlError(err)
		return false
	}
	return true
}

// Delete the existing DB record.
func (db Database) UserDelete(userId uint64) bool {
	query := `DELETE FROM app_users WHERE id = ?`
	_, err := db.Sql.Exec(query, userId)

	if err != nil {
		logSqlError(err)
		return false
	}
	return true
}

// Get the user by his email.
func (db Database) UserByEmail(email string) *model_database.User {
	user := new(model_database.User)
	query := `SELECT * FROM app_users WHERE email = ?`
	err := db.Sql.Get(user, query, email)

	if err != nil {
		logSqlError(err)
		return nil
	}
	return user
}

// Get the user by his identificator.
func (db Database) UserById(userId uint64) *model_database.User {
	user := new(model_database.User)
	query := `SELECT * FROM app_users WHERE id = ?`
	err := db.Sql.Get(user, query, userId)

	if err != nil {
		logSqlError(err)
		return nil
	}
	return user
}

// Get the user by his username.
func (db Database) UserByUsername(username string) *model_database.User {
	user := new(model_database.User)
	query := `SELECT * FROM app_users WHERE username = ?`
	err := db.Sql.Get(user, query, username)

	if err != nil {
		logSqlError(err)
		return nil
	}
	return user
}

func (db Database) UserOwnGroupList(userId uint64) []model_database.Group {
	groupList := &[]model_database.Group{}
	query := `SELECT
		app_groups.*
	FROM app_groups
	LEFT JOIN app_group_members ON app_groups.id = app_group_members.group_id
	WHERE (user_id = ? AND is_owner = 1)`
	err := db.Sql.Select(groupList, query, userId)

	if err != nil {
		logSqlError(err)
		return *groupList
	}
	return *groupList
}

func (db Database) UserGroupList(userId uint64) []model_database.Group {
	groupList := &[]model_database.Group{}
	query := `SELECT
		app_groups.*
	FROM app_groups
	LEFT JOIN app_group_members ON app_groups.id = app_group_members.group_id
	WHERE user_id = ?`
	err := db.Sql.Select(groupList, query, userId)

	if err != nil {
		logSqlError(err)
		return *groupList
	}
	return *groupList
}
