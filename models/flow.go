package models

import (
	"time"
	"sort"
	"strings"
	"strconv"
	"github.com/go-sql-driver/mysql"
)

type Type struct {
	Id           int         `json:"id" form:"id"`
	Name         string      `json:"name" form:"name"`
	Cms          string		 `json:"cms" form:"cms"`
}

func GetTypes() ([]*Type, error) {
	types := make([]*Type, 0)
	rows, err := MySqlDB.Query("SELECT id, name, cms FROM type")
	defer rows.Close()

	if err != nil {
		return types, err
	}

	for rows.Next() {
		t := new(Type)
		rows.Scan(&t.Id, &t.Name, &t.Cms)
		types = append(types, t)
	}
	if err = rows.Err(); err != nil {
		return types, err
	}
	return types, err
}

func GetTypeById(id int) (*Type, error) {
	t := new(Type)
	err := MySqlDB.QueryRow("SELECT id, name, cms FROM type where id = ?",
		id).Scan(&t.Id, &t.Name, &t.Cms)
	if err != nil {
		panic(err)
	}
	return t, err
}

type State struct {
	Id           int         `json:"id" form:"id"`
	Name         string      `json:"name" form:"name"`
	Description  string      `json:"description"`
	TypeId       int         `json:"type_id"`
	UsersId      string      `json:"users_id"`
	IsEnd        bool        `json:"is_end"`
}

func GetStateById(id int) (*State, error) {
	s := new(State)
	err := MySqlDB.QueryRow("SELECT id, name, description, type_id, is_end, users FROM state where id = ?",
		id).Scan(&s.Id, &s.Name, &s.Description, &s.TypeId, &s.IsEnd, &s.UsersId)
	if err != nil {
		panic(err)
	}
	return s, err
}

func Perm(userId, stateId int) bool {
	state, _ := GetStateById(stateId)
	permUsers := strings.Split(state.UsersId, ",")
	for _, v := range permUsers {
		if strconv.Itoa(userId) == v {
			return true
		}
	}
	return false
}

type Instance struct {
	Id        		int    		`json:"id"`
	TypeId    		int    		`json:"type_id"`
	Workflow  		string		`json:"workflow"`
	CurrentStateId  int   		`json:"current_state_id"`
	CreateTime      time.Time 	`json:"create_time"`
	Description  	string      `json:"description"`
	CreateUserId    int       	`json:"create_user_id"`
	CompleteTime    mysql.NullTime   `json:"complete_time"`
}

func GetInstances() ([]*Instance, error) {
	ins := make([]*Instance, 0)
	rows, err := MySqlDB.Query("SELECT id, type_id, workflow, current_state_id," +
		" create_time, description, create_user_id, complete_time FROM instance order by create_time desc")
	defer rows.Close()

	if err != nil {
		return ins, err
	}

	for rows.Next() {
		i := new(Instance)
		rows.Scan(&i.Id, &i.TypeId, &i.Workflow, &i.CurrentStateId, &i.CreateTime, &i.Description,
			&i.CreateUserId, &i.CompleteTime)
		ins = append(ins, i)
	}
	if err = rows.Err(); err != nil {
		return ins, err
	}
	return ins, err
}

func GetInstancesByQuery(query string) ([]*Instance, error) {
	ins := make([]*Instance, 0)
	rows, err := MySqlDB.Query(query)
	defer rows.Close()

	if err != nil {
		return ins, err
	}

	for rows.Next() {
		i := new(Instance)
		rows.Scan(&i.Id, &i.TypeId, &i.Workflow, &i.CurrentStateId, &i.CreateTime, &i.Description,
			&i.CreateUserId, &i.CompleteTime)
		ins = append(ins, i)
	}
	if err = rows.Err(); err != nil {
		return ins, err
	}
	return ins, err
}

func GetInstanceById(id int) (*Instance, error) {
	i := new(Instance)
	err := MySqlDB.QueryRow("SELECT id, type_id, workflow, current_state_id, create_time," +
		" description, create_user_id, complete_time FROM instance where id = ? ", id).Scan(
		&i.Id, &i.TypeId, &i.Workflow, &i.CurrentStateId, &i.CreateTime, &i.Description,
		&i.CreateUserId, &i.CompleteTime)
	if err != nil {
		panic(err)
	}
	return i, err
}

// 如果执行同意则进行改处理，工单状态向下走一步
// 判断如果是结束状态则结束工单
func (obj Instance) Next() {
	workflow := strings.Split(obj.Workflow, ",")
	stateIndex := sort.SearchStrings(workflow, strconv.Itoa(obj.CurrentStateId))
	var nextStateId int
	if stateIndex == (len(workflow)-2) {
		nextStateId, _ = strconv.Atoi(workflow[stateIndex+2])
	} else {
		nextStateId, _ = strconv.Atoi(workflow[stateIndex+1])
	}
	state, _ := GetStateById(nextStateId)
	if state.IsEnd {
		completeTime := time.Now()
		stmt, _ := MySqlDB.Prepare(`UPDATE instance SET current_state_id=?,complete_time=? WHERE id=?`)
		stmt.Exec(nextStateId, completeTime, obj.Id)
	} else {
		stmt, _ := MySqlDB.Prepare(`UPDATE instance SET current_state_id=? WHERE id=?`)
		stmt.Exec(nextStateId, obj.Id)
	}
}

type History struct {
	Id        		int    		`json:"id"`
	InstanceId      int         `json:"instance_id"`
	InstanceStateId int         `json:"instance_state_id"`
	Note            string      `json:"note"`
	CreateUserId    int         `json:"create_user_id"`
	CreateTime      time.Time   `json:"create_time"`
	Result          bool        `json:"result"`
}

func GetHistoryByInstanceId (id int) ([]*History, error) {
	his := make([]*History, 0)
	rows, err := MySqlDB.Query("SELECT id, instance_id, instance_state_id, note," +
		" create_user_id, create_time, result FROM history where instance_id = " + strconv.Itoa(id))
	defer rows.Close()

	if err != nil {
		return his, err
	}

	for rows.Next() {
		h := new(History)
		rows.Scan(&h.Id, &h.InstanceId, &h.InstanceStateId, &h.Note, &h.CreateUserId,
			&h.CreateTime, &h.Result)
		his = append(his, h)
	}
	if err = rows.Err(); err != nil {
		return his, err
	}
	return his, err
}

func (obj Instance) Record(note string, userId, result int) {
	workflow := strings.Split(obj.Workflow, ",")
	stateIndex := sort.SearchStrings(workflow, strconv.Itoa(obj.CurrentStateId))
	nextStateId, _ := strconv.Atoi(workflow[stateIndex+1])
	stmt, _ := MySqlDB.Prepare("INSERT into history (instance_id, instance_state_id, " +
		"note, create_user_id, result) values (?,?,?,?,?)")
	stmt.Exec(obj.Id, nextStateId, note, userId, result)
}
