package main

//const database = "../../test_mock.db"
//const database = "file::memory:?cache=shared"

//
//func TestMain(m *testing.M) {
//	setup()
//	code := m.Run()
//	//dropTable()
//	os.Exit(code)
//}

//func setup() {
//	log.Println("Running TestMain...")
//	seedDatabase(database)
//}
//
//func seedDatabase(database string) {
//	log.Println("seeding test database")
//	var mockLinks = []adapters.Link{
//		{OriginalURL: "https://test.com", Hash: "12345678", Data: []adapters.DataPoints{}},
//		{OriginalURL: "mudmap.io", Hash: "abcdefgh", Data: []adapters.DataPoints{}},
//	}
//	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{})
//	//if err := db.Migrator().DropTable(adapters.Link{}, adapters.DataPoints{}); err != nil {
//	//	log.Fatalln("could not drop tables in `seedDatabase` func")
//	//}
//	adapters.InitialMigrations(database)
//	if err != nil {
//		log.Fatalln("failed to connect to mock database during seed")
//	}
//	if err := db.Debug().Create(&mockLinks).Error; err != nil {
//		log.Fatalln("failed to create `mockLinks` during seed")
//	}
//}
//
//func dropTable() {
//	log.Println("cleaning up test database")
//	db, _ := gorm.Open(sqlite.Open(database), &gorm.Config{})
//	if err := db.Migrator().DropTable(adapters.Link{}, adapters.DataPoints{}); err != nil {
//		log.Fatalln("could not drop tables in `seedDatabase` func")
//	}
//}
