package dbHelper

import (
	"ToDo/database/migrations"
	"ToDo/models"

	"github.com/google/uuid"
)

func CreateTodo(todo *models.Todo) error {
	query := `
	INSERT INTO todo(user_id, t_name, description)
	VALUES($1,$2,$3)
	RETURNING id
	`

	return migrations.DB.QueryRow(
		query,
		todo.UserID,
		todo.Name,
		todo.Description,
	).Scan(&todo.ID)
}

func GetTodo(id uuid.UUID, userID uuid.UUID) (models.Todo, error) {
	var todo models.Todo

	err := migrations.DB.Get(
		&todo,
		`SELECT id, user_id, t_name, description FROM todo WHERE id=$1 AND user_id=$2`,
		id,
		userID,
	)

	return todo, err
}

func GetTodos(userID uuid.UUID) ([]models.Todo, error) {
	var todos []models.Todo

	err := migrations.DB.Select(
		&todos,
		`SELECT id, user_id, t_name, description FROM todo WHERE user_id=$1 ORDER BY id`,
		userID,
	)

	return todos, err
}

func DeleteTodo(id uuid.UUID, userID uuid.UUID) (int64, error) {

	result, err := migrations.DB.Exec(
		`DELETE FROM todo WHERE id=$1 AND user_id=$2`,
		id,
		userID,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func UpdateTodo(id uuid.UUID, userID uuid.UUID, todo *models.Todo) (int64, error) {

	result, err := migrations.DB.Exec(
		`UPDATE todo
				SET t_name=$1,
		    	description=$2
				WHERE id=$3 AND user_id=$4`,
		todo.Name,
		todo.Description,
		id,
		userID,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
