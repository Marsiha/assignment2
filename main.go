package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Student struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Age          uint //new column
	Email        string
	DepartmentID uint
	Department   Department `gorm:"foreignKey:DepartmentID"`
	Enrollments  []Enrollment
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Course struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	DepartmentID uint
	InstructorID uint
	Department   Department `gorm:"foreignKey:DepartmentID"`
	Instructor   Instructor `gorm:"foreignKey:InstructorID"`
	Students     []Student  `gorm:"many2many:enrollments"`
}

type Department struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type Enrollment struct {
	ID        uint `gorm:"primaryKey"`
	StudentID uint
	CourseID  uint
	Grade     string
	Student   Student `gorm:"foreignKey:StudentID"`
	Course    Course  `gorm:"foreignKey:CourseID"`
}

type Instructor struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	DepartmentID uint
	Department   Department `gorm:"foreignKey:DepartmentID"`
	Courses      []Course
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CourseWithEnrolledStudents struct {
	CourseID   uint
	CourseName string
	Enrolled   int
}

func main() {
	dsn := "root:possible04@tcp(localhost:3306)/golang"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("error")
	}

	db.AutoMigrate(&Student{}, &Department{}, &Course{}, &Enrollment{}, &Instructor{})

	db.Model(&Student{}).Where("Name = ?", "mars").Update("Name", "Mars")          //update with condition
	db.Create(&Student{Name: "adsf", Email: "asdfa@example.com", DepartmentID: 1}) //create data
	db.Delete(&Student{}, 9)                                                       //Delete data with Id

	//take all data from table
	var students []Student
	db.Find(&students)
	for _, student := range students {
		fmt.Printf("ID: %d, Name: %s, Email: %s, DepartmentID: %d\n", student.ID, student.Name, student.Email, student.DepartmentID)
	}

	//add new student
	newStudent(db, "mars", "asdf", 1)

	//retrieve course
	fmt.Println(retrieveCourse(db, 2))
	var courses []Course
	db.Find(&courses)
	for _, course := range courses {
		fmt.Printf("ID: %d, Name:%s, DepartmentID: %d, InstructorID: %d\n", course.ID, course.Name, course.DepartmentID, course.InstructorID)
	}

	//update instructor name
	updateInstructor(db, 1, "mars")
	var instructors []Instructor
	db.Find(&instructors)
	for _, instructor := range instructors {
		fmt.Printf("ID: %d, Name: %s, DepartmentID: %d\n", instructor.ID, instructor.Name, instructor.DepartmentID)
	}

	//delete department with id
	deleteDepartment(db, 5).Error()
	var departments []Department
	db.Find(&departments)
	for _, department := range departments {
		fmt.Printf("ID: %d, Name: %s\n", department.ID, department.Name)
	}

	studentFromDepartment(db, 1)
	courseWithInstructor(db, 2)
	studentEnrollment(db, 2)
	enrollStudentInCourse(db, 1, 1, "A")

	student := Student{Name: "Alice", Email: "alice@example.com", DepartmentID: 1}
	db.Create(&student)

	//Soft delete the student
	err = deletedStudents(db, 14)
	if err != nil {
		panic(err)
	}

	activeStudents(db)

	//Custom query

	studentInDepartment(db, 6)
	enrolled(db)
}

func newStudent(db *gorm.DB, name, email string, departmentID uint) error {
	student := Student{Name: name, Email: email, DepartmentID: departmentID}
	return db.Create(&student).Error
}

func retrieveCourse(db *gorm.DB, courseID uint) (uint, error) {
	var course Course
	result := db.First(&course, courseID)
	return course.ID, result.Error
}

func updateInstructor(db *gorm.DB, instructorID uint, newName string) error {
	var instructor Instructor
	result := db.First(&instructor, instructorID)
	if result.Error != nil {
		return result.Error
	}
	instructor.Name = newName
	return db.Save(&instructor).Error
}

func deleteDepartment(db *gorm.DB, departmentID uint) error {
	var department Department
	result := db.First(&department, departmentID)
	if result.Error != nil {
		return result.Error
	}
	return db.Delete(&department).Error
}

// Querying
func studentFromDepartment(db *gorm.DB, departmentId uint) {
	var students []Student
	db.Find(&students, "department_id = ?", departmentId)
	for _, student := range students {
		fmt.Printf("Name: %s, DepartmentID: %d\n", student.Name, student.DepartmentID)
	}
}

func courseWithInstructor(db *gorm.DB, instructorId uint) {
	var courses []Course
	db.Find(&courses, "instructor_id = ?", instructorId)
	for _, course := range courses {
		fmt.Printf("Course name: %s, InstructorID: %d", course.Name, course.InstructorID)
	}
}

func studentEnrollment(db *gorm.DB, studentID uint) {
	var enrolls []Enrollment
	db.Find(&enrolls, "student_id = ?", studentID)
	for _, enroll := range enrolls {
		fmt.Printf("StudentID: %d, CourseID: %d, Grade: %s\n", enroll.StudentID, enroll.CourseID, enroll.Grade)
	}
}

// Transactions
func enrollStudentInCourse(db *gorm.DB, studentID, courseID uint, grade string) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	enrollment := Enrollment{StudentID: studentID, CourseID: courseID, Grade: grade}
	if err := tx.Create(&enrollment).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Hooks
func (student *Student) BeforeCreate(tx *gorm.DB) (err error) {
	student.CreatedAt = time.Now()
	student.UpdatedAt = time.Now()
	return nil
}

func (instructor *Instructor) BeforeUpdate(tx *gorm.DB) (err error) {
	instructor.UpdatedAt = time.Now()
	return nil
}

// Self delete
func deletedStudents(db *gorm.DB, studentID uint) error {
	result := db.Model(&Student{}).Where("id = ?", studentID).Update("deleted_at", time.Now())
	return result.Error
}

func activeStudents(db *gorm.DB) {
	var students []Student
	db.Where("deleted_at IS NULL").Find(&students)
	// Print active students
	fmt.Println("Active Students:")
	for _, s := range students {
		fmt.Printf("ID: %d, Name: %s, Email: %s, DepartmentID: %d\n", s.ID, s.Name, s.Email, s.DepartmentID)
	}
}

func studentInDepartment(db *gorm.DB, departmentID uint) {
	var count int64
	db.Model(&Student{}).Where("department_id = ?", departmentID).Count(&count)
	fmt.Printf("DepartmentID: %d, Students count: %d", departmentID, count)
}

// enrolled
func enrolled(db *gorm.DB) ([]CourseWithEnrolledStudents, error) {
	var courses []CourseWithEnrolledStudents
	query := db.Table("courses c").
		Select("c.id AS course_id, c.name AS course_name, COUNT(e.student_id) AS enrolled").
		Joins("LEFT JOIN enrollments e ON c.id = e.course_id").
		Group("c.id, c.name").
		Scan(&courses)
	//query.Find(&courses)
	//for _, course := range courses {
	//	fmt.Printf("Course name: %s, Enrolled: %d\n", course.CourseName, course.Enrolled)
	//}
	return courses, query.Error
}

func GetStudentCountByDepartmentID(db *gorm.DB, departmentId uint) int64 {
	var count int64
	db.Model(&Student{}).Where("department_id = ?", departmentId).Count(&count)
	return count
}
