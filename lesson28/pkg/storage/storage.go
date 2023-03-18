package storage

import (
	"errors"
	"fmt"
	"lesson28/pkg/student"
	"log"
)

type StudentClass map[string]*student.Student

func (sc StudentClass) Put(student *student.Student) error {
	if sc.Get(student.Name) != nil {
		return errors.New(fmt.Sprintf("Student %s already is in the class", student.Name))
	}
	sc[student.Name] = student
	log.Printf("[Added] %s : %+v\n", student.Name, student)
	return nil
}

func (sc StudentClass) Get(studentName string) *student.Student {
	if sc[studentName] != nil {
		return sc[studentName]
	}
	log.Printf("[Not Found] %s\n", studentName)
	return nil
}

func (sc StudentClass) Form() {
	fmt.Println("Expected input: `<name> <age> <grade>`\nExample: Dima 38 18\nHit `Ctrl+Z` to stop input")
	for {
		var (
			name  string
			age   int
			grade int
		)

		_, err := fmt.Scanf("%s %d %d", &name, &age, &grade)
		if err != nil {
			log.Printf("[Error] %s\n", err)
		}
		err = sc.Put(student.NewStudent(name, age, grade))
		if err != nil {
			log.Printf("[Error] %s\n", err)
		}
	}
}
