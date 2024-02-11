package mysql

import (
	"alexedwards.net/snippetbox/pkg/models"
	"database/sql"
)

type MealModel struct {
	DB *sql.DB
}

func (m *MealModel) Insert(mealName, weekday string, quantity int) (int, error) {
	stmt := "INSERT INTO canteen_menu (meal_name, weekday, quantity) VALUES (?, ?, ?)"

	result, err := m.DB.Exec(stmt, mealName, weekday, quantity)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *MealModel) Latest() ([]*models.Meal, error) {
	stmt := `SELECT id, meal_name, weekday, quantity FROM canteen_menu
ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meals := []*models.Meal{}
	for rows.Next() {
		meal := &models.Meal{}
		err = rows.Scan(&meal.ID, &meal.MealName, &meal.Weekday, &meal.Quantity)
		if err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return meals, nil
}
