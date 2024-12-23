package model

import (
	"log"
	"strconv"
)

func CreateTodo(p ProposedTodo) (bool, error) {
	log.Println("INFO: Todo creation requested: " + p.Description)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	// get the id for the "new" status
	statusId, err := GetStatusByName("new")
	if err != nil {
		log.Println("ERROR: Could not retrieve status Id for status 'new':" + string(err.Error()))
		return false, err
	}

	q, err := t.Prepare("INSERT INTO Todo (Description, Status) VALUES (?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(p.Description, statusId)
	if err != nil {
		log.Println("ERROR: Cannot create todo with description '" + p.Description + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Todo with description '" + p.Description + "' created")
	return true, nil
}

func DeleteTodo(id int) (bool, error) {
	idString := strconv.Itoa(id)
	log.Println("INFO: Todo deletion requested: " + idString)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := DB.Prepare("DELETE FROM Todos WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(id)
	if err != nil {
		log.Println("ERROR: Cannot delete user '" + idString + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Todo with Id '" + idString + "' has been deleted")
	return true, nil
}

func GetTodos() ([]Todo, error) {
	return nil, nil
}

func GetTodoById(id int) (Todo, error) {
	return Todo{}, nil
}

func UpdateTodo(id int, statusId int) (Todo, error) {
	return Todo{}, nil
}
