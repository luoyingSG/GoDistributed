package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

func (s Student) Average() float32 {
	var result float32
	for _, v := range s.Grades {
		result += v.Score
	}

	return result / float32(len(s.Grades))
}

type Students []Student

var (
	students      Students
	studentsMutex sync.Mutex
)

func (ss Students) GetByID(ID int) (*Student, error) {
	for _, s := range ss {
		if s.ID == ID {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("student with ID %d not found", ID)
}

type GradeType string

const (
	GradeQuiz = GradeType("quiz")
	GradeTest = GradeType("test")
	GradeExam = GradeType("exam")
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}
