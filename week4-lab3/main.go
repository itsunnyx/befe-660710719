package main 

import (
	"fmt"
	"errors"
)
type Student struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Year int `json:"year"`
	GPA float64 `json:"gpa"`
}

func (s *Student) IsHonor() bool {
	return s.GPA >= 3.50
}

func (s *Student) Validate() error {
	if (s.Name == "") {
		return errors.New("name is required")
	}
	if (s.Year < 1 || s.Year >4) {
		return errors.New("year must be between 1-4")
	}
	if (s.GPA < 0 || s.GPA >4 ) {
		return errors.New("gpa must be between 0-4")
	}
	return nil
}

func main() {
	// var st Student = Student{
	// 	ID : "1234",
	// 	Name : "",
	// 	Email : "Sittipornchaiya_n@su.ac.th",
	// 	Year : 3,
	// 	GPA : 3.5,
	// }

	// st := Student({
	// 	ID : "1234",
	// 	Name : "Nanthawat",
	// 	Email : "Sittipornchaiya_n@su.ac.th",
	// 	Year : 3,
	// 	Gpa : 3.5
	// })


	students := []Student{
		{
			ID : "1",
			Name : "Nanthawat",
			Email : "Sittipornchaiya_n@su.ac.th",
			Year : 3,
			GPA : 3.5,
		},
		{
			ID : "2",
			Name : "ball",
			Email : "ball123@su.ac.th",
			Year : 3,
			GPA : 1,
		},
	}

	newStudent := Student{
		ID : "3",
		Name : "ronaldo",
		Email : "ronaldo7@su.ac.th",
		Year : 3,
		GPA : 2,
	}

	students = append(students, newStudent)

	for i, student := range students {
		fmt.Printf("%v Honor = %v\n", i+1, student.IsHonor())
		fmt.Printf("%vValidated = %v\n", i+1, student.Validate())
	}

	
	

}