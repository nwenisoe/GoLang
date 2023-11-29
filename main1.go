package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Teacher struct {
	ID       bson.ObjectId   `bson:"_id,omitempty" json:"id"`
	Name     string          `json:"name"`
	Students []bson.ObjectId `bson:"students,omitempty" json:"students"`
}

type Student struct {
	ID       bson.ObjectId   `bson:"_id,omitempty" json:"id"`
	Name     string          `json:"name"`
	Teachers []bson.ObjectId `bson:"teachers,omitempty" json:"teachers"`
}

func main() {
	info := mgo.DialInfo{
		Addrs:   []string{"localhost:27017"},
		Timeout: 5 * time.Second,
	}

	client, err := mgo.DialWithInfo(&info)
	if err != nil {
		panic("Error occurred")
	}
	defer client.Close()

	db := client.DB("school")
	teachersColl := db.C("teachers")
	studentsColl := db.C("students")

	student := []interface{}{
		Student{
			ID:   bson.NewObjectId(),
			Name: "Theint Theint Ko",
		},
		Student{
			ID:   bson.NewObjectId(),
			Name: "Poe Phyu Thae",
		},
		Student{
			ID:   bson.NewObjectId(),
			Name: "Yin Phyu Aung",
		},
		Student{
			ID:   bson.NewObjectId(),
			Name: "Yu Ya Kyaw",
		},
		Student{
			ID:   bson.NewObjectId(),
			Name: "Nweni Soe",
		},
	}

	teacher := []interface{}{
		Teacher{
			ID:       bson.NewObjectId(),
			Name:     "Mt. Kyaw Kyaw Lwin Thant",
			Students: []bson.ObjectId{student[0].(Student).ID, student[1].(Student).ID},
		},
		Teacher{
			ID:       bson.NewObjectId(),
			Name:     "Mr.Aye Min",
			Students: []bson.ObjectId{student[0].(Student).ID, student[1].(Student).ID},
		},
		Teacher{
			ID:       bson.NewObjectId(),
			Name:     "Mt.Min Lwin",
			Students: []bson.ObjectId{student[0].(Student).ID, student[1].(Student).ID, student[2].(Student).ID, student[3].(Student).ID},
		},
		Teacher{
			ID:       bson.NewObjectId(),
			Name:     "Mt. Pyae Phyo",
			Students: []bson.ObjectId{student[0].(Student).ID, student[1].(Student).ID, student[2].(Student).ID, student[3].(Student).ID},
		},
		Teacher{
			ID:       bson.NewObjectId(),
			Name:     "Mt. Min Khant Thu",
			Students: []bson.ObjectId{student[1].(Student).ID, student[3].(Student).ID, student[4].(Student).ID},
		},
		Teacher{
			ID:   bson.NewObjectId(),
			Name: "Mr.Wyut Yee Thant",
		},
	}

	err = teachersColl.Insert(teacher...)
	if err != nil {
		panic(err)
	}

	err = studentsColl.Insert(student...)
	if err != nil {
		panic(err)
	}

	// Update the relationships
	err = updateTeacherStudents(teachersColl, teacher[0].(Teacher).ID, student[4].(Student).ID)
	if err != nil {
		panic(err)
	}

	err = updateStudentTeachers(studentsColl, student[1].(Student).ID, teacher[0].(Teacher).ID)
	if err != nil {
		panic(err)
	}

	// Query and print the data
	printData(teachersColl, studentsColl)

	// Retrieve a specific teacher by ID
	teacherID := bson.ObjectIdHex("6565b780bd5c962118fa8a4a")
	retrievedTeacher, err := getTeacherByID(teachersColl, teacherID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\nRetrieved Teacher:\nID: %s, Name: %s, Students: %v\n",
		retrievedTeacher.ID, retrievedTeacher.Name, retrievedTeacher.Students)

	// Retrieve a specific student by ID
	studentID := bson.ObjectIdHex("6565b780bd5c962118fa8a45")
	retrievedStudent, err := getStudentByID(studentsColl, studentID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\nRetrieved Student:\nID: %s, Name: %s, Teachers: %v\n",
		retrievedStudent.ID, retrievedStudent.Name, retrievedStudent.Teachers)

	// Delete the relationship
	err = deleteRelationship(teachersColl, studentsColl, teacherID, studentID)
	if err != nil {
		fmt.Println(err)
	}

	// Query and print the data after deletion
	fmt.Println("\nAfter Deletion:")
	printData(teachersColl, studentsColl)

	// Remove teacher from student's list
	err = studentsColl.Update(
		bson.M{"_id": studentID},
		bson.M{"$pull": bson.M{"teachers": teacherID}},
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("\nAfter updating..")
	printData(teachersColl, studentsColl)

}

func getTeacherByID(coll *mgo.Collection, teacherID bson.ObjectId) (Teacher, error) {
	var teacher Teacher
	err := coll.Find(bson.M{"_id": teacherID}).One(&teacher)
	return teacher, err
}

func getStudentByID(coll *mgo.Collection, studentID bson.ObjectId) (Student, error) {
	var student Student
	err := coll.Find(bson.M{"_id": studentID}).One(&student)
	return student, err
}

func updateTeacherStudents(coll *mgo.Collection, teacherID, studentID bson.ObjectId) error {
	selector := bson.M{"_id": teacherID}
	update := bson.M{"$push": bson.M{"students": studentID}}
	return coll.Update(selector, update)
}

func updateStudentTeachers(coll *mgo.Collection, studentID, teacherID bson.ObjectId) error {
	selector := bson.M{"_id": studentID}
	update := bson.M{"$push": bson.M{"teachers": teacherID}}
	return coll.Update(selector, update)
}
func deleteRelationship(teachersColl, student *mgo.Collection, teacherID, studentID bson.ObjectId) error {
	// Remove student from teacher's list
	err := teachersColl.Update(
		bson.M{"_id": teacherID},
		bson.M{"$pull": bson.M{"students": studentID}},
	)
	if err != nil {
		return err
	}

	return nil
}

func printData(teachersColl, studentsColl *mgo.Collection) {
	// Query and print teachers
	var teachers []Teacher
	err := teachersColl.Find(nil).All(&teachers)
	if err != nil {
		panic(err)
	}

	fmt.Println("Teachers:")
	for _, t := range teachers {
		fmt.Printf("ID: %s, Name: %s, Students: %v\n", t.ID, t.Name, t.Students)
	}

	// Query and print students
	var students []Student
	err = studentsColl.Find(nil).All(&students)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nStudents:")
	for _, s := range students {
		fmt.Printf("ID: %s, Name: %s, Teachers: %v\n", s.ID, s.Name, s.Teachers)
	}
}
