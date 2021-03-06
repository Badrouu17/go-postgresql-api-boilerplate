package queries

import "fmt"

func InsertUser(name, email, password string) string {
	return fmt.Sprintf("INSERT INTO users_user (name,email,password) VALUES ('%s','%s','%s') RETURNING id,name,email", name, email, password)
}

func GetUserWithEmail(email string) string {
	return fmt.Sprintf("SELECT id,name,email,password FROM users_user WHERE users_user.email = '%s'", email)
}

func GetUserWithId(id int) string {
	return fmt.Sprintf("SELECT id,name,email,password,password_changed_at FROM users_user WHERE id = %d", id)
}

func UpdateUserPassResetData(id int, prt string, pre int64) string {
	return fmt.Sprintf("UPDATE users_user SET password_reset_token = '%s' ,password_reset_expires = %d WHERE id = %d", prt, pre, id)
}

func GetUserByResetToken(prt string, now int64) string {
	return fmt.Sprintf("SELECT id,name,email FROM users_user WHERE users_user.password_reset_token = '%s' AND users_user.password_reset_expires > %d", prt, now)
}

func ResetPassword(id int, ps string, pca int64) string {
	return fmt.Sprintf("UPDATE users_user SET password = '%s', password_reset_token = NULL, password_reset_expires = 0, password_changed_at = %d WHERE id = %d", ps, pca, id)
}

func UpdateUserEmail(id int, email string) string {
	return fmt.Sprintf("UPDATE users_user SET email = '%s' WHERE id = %d", email, id)
}

func UpdateUserName(id int, name string) string {
	return fmt.Sprintf("UPDATE users_user SET name = '%s' WHERE id = %d", name, id)
}

func DeleteUser(id int) string {
	return fmt.Sprintf("DELETE FROM users_user WHERE id = %d", id)
}
