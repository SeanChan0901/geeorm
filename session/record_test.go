package session

import "testing"

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func testSqlite3RecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSqlite3Session().Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}

func testMysqlRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewMysqlSession().Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}

func TestSqlite3Session_Insert(t *testing.T) {
	s := testSqlite3RecordInit(t)
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSqlite3Session_Find(t *testing.T) {
	s := testSqlite3RecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}
}

func TestSqlite3Session_Limit(t *testing.T) {
	s := testSqlite3RecordInit(t)
	var users []User
	err := s.Limit(1).Find(&users)
	if err != nil || len(users) != 1 {
		t.Fatal("failed to query with limit condition")
	}
}

func TestSqlite3Session_Update(t *testing.T) {
	sqlite3Session := testSqlite3RecordInit(t)
	affected, _ := sqlite3Session.Where("Name = ?", "Tom").Update("Age", 30)
	u := &User{}
	_ = sqlite3Session.OrderBy("Age DESC").First(u)

	if affected != 1 || u.Age != 30 {
		t.Fatal("failed to update")
	}
}

func TestSqlite3Session_DeleteAndCount(t *testing.T) {
	s := testSqlite3RecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Delete()

	count, _ := s.Count()

	if affected != 1 || count != 1 {
		t.Fatal("failed to delete or count")
	}
}

func TestMysqlSession_Insert(t *testing.T) {
	s := testMysqlRecordInit(t)
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestMysqlSession_Limit(t *testing.T) {
	s := testMysqlRecordInit(t)
	var users []User
	err := s.Limit(1).Find(&users)
	if err != nil || len(users) != 1 {
		t.Fatal("failed to query with limit condition")
	}
}

func TestMysqlSession_Find(t *testing.T) {
	s := testMysqlRecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}
}

func TestMysqlSession_Update(t *testing.T) {
	s := testMysqlRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Update("Age", 30)
	u := &User{}
	_ = s.OrderBy("Age DESC").First(u)

	if affected != 1 || u.Age != 30 {
		t.Fatal("failed to update")
	}
}

func TestMysqlSession_DeleteAndCount(t *testing.T) {
	s := testMysqlRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Delete()

	count, _ := s.Count()

	if affected != 1 || count != 1 {
		t.Fatal("failed to delete or count")
	}
}
