
/*import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Student struct {
	ID       bson.ObjectId `bson:"s_id,omitempty"`
	SName    string        `bson:"s_name"`
	SAddress string        `bson:"s_address"`
	Section  string        `bson:"section"`
	Teacher  []Teacher     `bson:"t_id"`
}

type Teacher struct {
	ID       bson.ObjectId `bson:"t_id,omitempty"`
	TName    string        `bson:"t_name"`
	TAddress string        `bson:"t_address"`
	Students []Student     `bson:"s_id,omitempty"`
}

func main() {
	info := mgo.DialInfo{
		Addrs:   []string{"localhost:27017"},
		Timeout: 5 * time.Second,
	}

	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		panic("Error occurred")
	}
	defer session.Close()

	scollect := session.DB("school").C("students")
	tcollect := session.DB("school").C("teachers")

	// Ensure index on teacher_id field
	//index := mgo.Index{
	//	Key:        []string{"teacher_id"},
	//	Background: true,
	//}

	//err = scollect.EnsureIndex(index)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Create a teacher with a valid ObjectId
	teacher := []interface{}{
		Teacher{
			TName:    "Mr. Aye Min",
			TAddress: "Kamaryut",
			ID:       bson.NewObjectId(),
		},
		Teacher{
			TName:    "Mr. Taylor",
			TAddress: "123 Yankim St",
			ID:       bson.NewObjectId(),
		},
		Teacher{
			TName:    "Mr. Khin",
			TAddress: "Mingalar St",
			ID:       bson.NewObjectId(),
		},
		Teacher{
			TName:    "Mr.John Doe",
			TAddress: "PyiGyiTakhon St",
			ID:       bson.NewObjectId(),
		},
	}

	err = tcollect.Insert(teacher...)
	if err != nil {
		log.Fatal(err)
	}

	// Create a student associated with the teacher
	student := []interface{}{
		Student{
			SName:    "Theint Theint Ko",
			Section:  "D",
			SAddress: "Taungoo",
			Teacher:  bson.NewObjectId(),
		},
		Student{
			SName:    "Poe Phyu Thae",
			Section:  "B",
			SAddress: "Phyu",
			Teacher:  bson.NewObjectId(),
		},
		//Student{
		//	SName:    "Su Wai Hlaing",
		//	Section:  "A",
		//	SAddress: "SanYeikNyein",
		//	Teacher:  bson.ObjectId(""),
		//},
		//Student{
		//	SName:    "Aye Myo Thant",
		//	Section:  "C",
		//	SAddress: "Hlaing",
		//	Teacher:  bson.ObjectId(""),
		//},
		//Student{
		//	SName:    "Yin Phyu Phyu Aung",
		//	Section:  "A",
		//	SAddress: "HtarWel",
		//	Teacher:  bson.ObjectId(""),
		//},
	}

	err = scollect.Insert(student...)
	if err != nil {
		log.Fatal(err)
	}

	// Query to retrieve the teacher and associated students
	/*	var queriedTeacher Teacher
		err = tcollect.Find(bson.M{"t_id": "6561c17b7986a5ff0d3924b5"}).One(&queriedTeacher)
		if err != nil {
			log.Fatal(err)
		}

		var queriedStudents []Student
		err = scollect.Find(bson.M{"t_id": teacher}).All(&queriedStudents)
		if err != nil {
			log.Fatal(err)
		}

		// Print results
		fmt.Printf("Queried Teacher: %+v\n", queriedTeacher)
		fmt.Printf("Associated Students: %+v\n", queriedStudents)}*/

