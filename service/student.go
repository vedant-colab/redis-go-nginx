package service

import (
	"context"
	"fmt"
	"redisgo/db"

	"github.com/jackc/pgx/v4"
)

type Student struct {
	SID    int
	Name   string
	Age    int
	Course string
}

func FetchStudentById(ctx context.Context, id int) (*Student, error) {
	query := "SELECT sid, name, age, course FROM student WHERE sid = $1"
	row := db.DB.QueryRow(ctx, query, id)

	var student Student
	err := row.Scan(&student.SID, &student.Name, &student.Age, &student.Course)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("student with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch student: %w", err)
	}

	return &student, nil
}
