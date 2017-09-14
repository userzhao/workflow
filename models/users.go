package models

import (
	"strings"
	"strconv"
	"fmt"
)

type NullInt struct {
	Val     int
	Valid   bool
}

func (nt *NullInt) Scan(value interface{}) (err error) {
	if value == nil {
		nt.Val, nt.Valid = 0, false
		return
	}

	switch v := value.(type) {
	case int:
		nt.Val, nt.Valid = v, true
		return
	}

	nt.Valid = false
	return
}

type User struct {
	Id            int        `json:"id"`
	UserName      string     `json:"username" form:"username"`
	FirstName     string     `json:"name" form:"name"`
	Email         string     `json:"email" form:"email"`
	SuperiorId    int        `json:"superior_id"`
	DepartmentUserID    NullInt  `json:"department_user_id"`
}

func GetUserById(id int) *User {
	user := new(User)
	err := CasSqlDB.QueryRow("select A.id, A.first_name, A.username, A.email, B.superior_id, " +
		"C.user_id as department_user_id  from (auth_user A , cas_pro B) join " +
		"cas_department C on B.department_id = C.id  where A.id = B.user_id and A.id = ?",
		id).Scan(&user.Id, &user.FirstName, &user.UserName, &user.Email, &user.SuperiorId, &user.DepartmentUserID)
	if err != nil {
		panic(err)
	}
	return user
}

func GetTranUsersById(usersId string) string {
	users := strings.Split(usersId, ",")
	firstName  := make([]string, 0)
	for  _, id := range users {
		userId, _ := strconv.Atoi(id)
		user := GetUserById(userId)
		firstName = append(firstName, user.FirstName)
	}
	return strings.Join(firstName, ",")
}

func GetSuborsIdById(userId int) []string {
	rows, _ := CasSqlDB.Query("select user_id from cas_pro where superior_id = " + strconv.Itoa(userId))
	defer rows.Close()

	var suborsId = make([]string, 0)
	for rows.Next() {
		u := new(User)
		rows.Scan(&u.Id)
		suborsId = append(suborsId, strconv.Itoa(u.Id))
	}
	if err := rows.Err(); err != nil {
		return suborsId
	}
	return suborsId
}

// 获取多有部门负责人的id
func GetDepartUsersId() []int {
	rows, _ := CasSqlDB.Query("select user_id  from cas_department")
	defer rows.Close()

	var suborsId []int
	for rows.Next() {
		u := new(User)
		rows.Scan(&u.DepartmentUserID)
		suborsId = append(suborsId, u.DepartmentUserID.Val)
	}
	return suborsId
}


// 根据userId 找出部门下所有人
func GetSubDepartUsersId(userId int) []string {
	rows, _ := CasSqlDB.Query(fmt.Sprintf("select user_id from cas_pro where department_id = (select " +
		"id from cas_department where user_id=%d)", userId))
	defer rows.Close()

	var suborsId []string
	for rows.Next() {
		u := new(User)
		rows.Scan(&u.Id)
		suborsId = append(suborsId, strconv.Itoa(u.Id))
	}
	return suborsId
}

