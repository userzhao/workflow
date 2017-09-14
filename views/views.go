package views

import (
	"github.com/gin-gonic/gin"
	"github.com/go/src/sort"
	"project/workflow/singals"
	"net/http"
	"strings"
	"syscall"
	"strconv"
	"fmt"
	"time"

	auth "project/workflow/middleware"
	"project/workflow/models"
	"project/workflow/utils"
)

//func LoginApi(c *gin.Context) {
//	var user models.User
//	tokenId := c.Query("ticket")
//	if tokenId == "" {
//		c.JSON(http.StatusForbidden, "ticket is not exists")
//		return
//	}
//	resp, err := http.Get("http://192.168.162.108:8005/validate?ticket=" + tokenId)
//	defer resp.Body.Close()
//	if err != nil {
//		c.JSON(http.StatusBadRequest, "validate 500")
//		return
//	}
//	body, _ := ioutil.ReadAll(resp.Body)
//	json.Unmarshal(body, &user)
//
//	if tokenId != "" && user.UserName != "" {
//		ses := auth.GlobalSessions.SessionStart(c.Writer, c.Request)
//		ses.Set("authenticated", true)
//		ses.Set("user", user)
//		c.Redirect(http.StatusFound, "index")
//	} else {
//		c.Redirect(http.StatusFound, "http://192.168.162.108:8005/index")
//	}
//}


// test login
func LoginApi(c *gin.Context) {
	var user = models.User{Id: 255, UserName: "yu.zhao", Email: "yu.zhao@100credit.com", FirstName: "赵宇"}

	ses := auth.GlobalSessions.SessionStart(c.Writer, c.Request)
	ses.Set("authenticated", true)
	ses.Set("user", user)

	c.Redirect(http.StatusFound, "index")
}

func LogoutApi(c *gin.Context) {
	auth.GlobalSessions.SessionDestroy(c.Writer, c.Request)
	c.Redirect(http.StatusFound, "http://192.168.162.108:8005/index")
}

func Index(c *gin.Context) {
	ses := auth.GlobalSessions.SessionStart(c.Writer, c.Request)
	user, _ := ses.Get("user").(models.User)

	show := c.DefaultQuery("type", "apply")

	types, _ := models.GetTypes()

	var ins []*models.Instance

	switch show {
	case "apply":
		query := fmt.Sprintf("SELECT id, type_id, workflow, current_state_id," +
			" create_time, description, create_user_id, complete_time FROM instance" +
			" where create_user_id = %d order by create_time desc", user.Id)
		ins, _ = models.GetInstancesByQuery(query)
	case "approval":
		queryDft := fmt.Sprintf("select id, type_id, workflow, current_state_id," +
			" create_time, description, create_user_id, complete_time from instance where complete_time is null and %d " +
			"in (select users from state where id=current_state_id) order by create_time desc", user.Id)
		insDft, _ := models.GetInstancesByQuery(queryDft)

		// 下级状态
		suborsId := models.GetSuborsIdById(user.Id)
		var insPro, insDep []*models.Instance
		if len(suborsId) >= 1 {
			queryPro := fmt.Sprintf("select id, type_id, workflow, current_state_id, create_time, " +
				"description, create_user_id, complete_time from instance where complete_time is null and" +
				" current_state_id = %d and create_user_id in (" + strings.Join(suborsId, ",") + ") " +
				"order by create_time desc", 2)
			insPro, _ = models.GetInstancesByQuery(queryPro)
		}

		// 部门状态
		depUsersId := models.GetSubDepartUsersId(user.Id)
		if utils.Exist(user.Id, models.GetDepartUsersId()) {
			queryDep := fmt.Sprintf("select id, type_id, workflow, current_state_id, create_time, " +
				"description, create_user_id, complete_time from instance where complete_time is null and" +
				" current_state_id = %d and create_user_id in (" + strings.Join(depUsersId, ",") + ") " +
				"order by create_time desc", 3)
			insDep, _ = models.GetInstancesByQuery(queryDep)
		}


		ins = utils.Merge(insDft, insPro, insDep)

	case "all":
		ins, _ = models.GetInstances()
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"user": user,
		"type": types,
		"objs": ins,
		"show": show,
	})
}

// 初始化工作流实例, 根据表单数据确认具体流程
// typeCms 中存储具体的流程, 每个流程的state的id ","分割, 多个流程用";"分割
func CreateTask(c *gin.Context) {
	note := c.PostForm("cms-note")
	typeId, _ := strconv.Atoi(c.PostForm("cms-type"))

	stmt, _ := models.MySqlDB.Prepare(`INSERT into instance (type_id, workflow,
	 	current_state_id, description, create_user_id) values (?,?,?,?,?)`)

	typeObj, _ := models.GetTypeById(typeId)
	typeCms := strings.Split(typeObj.Cms, ";")
	var workflow string
	if len(typeCms) == 1 {
		workflow = strings.Join(typeCms,",")
	} else {
		// 根据数据确定具体的流程
	}
	statesId := strings.Split(workflow, ",")
	currentStateId, _ := strconv.Atoi(statesId[0])

	ses := auth.GlobalSessions.SessionStart(c.Writer, c.Request)
	user, _ := ses.Get("user").(models.User)

	stmt.Exec(typeId, workflow, currentStateId, note, user.Id)
	c.JSON(http.StatusOK, gin.H{"msg": 1})
}

func TaskDetail(c *gin.Context) {
	objId, _ := strconv.Atoi(c.Param("objId"))
	ses := auth.GlobalSessions.SessionStart(c.Writer, c.Request)
	user, _ := ses.Get("user").(models.User)
	ins, _ := models.GetInstanceById(objId)

	statesId := strings.Split(ins.Workflow, ",")
	stateIndex := sort.SearchStrings(statesId, strconv.Itoa(ins.CurrentStateId))

	history, _ := models.GetHistoryByInstanceId(objId)

	c.HTML(http.StatusOK, "detail.html", gin.H{
		"user": user,
		"obj": ins,
		"states": statesId,
		"stateIndex": stateIndex,
		"history": history,
	})
}

func TranAction(c *gin.Context) {
	objId, _ := strconv.Atoi(c.PostForm("objId"))
	result, _ := strconv.Atoi(c.PostForm("result"))
	note := c.PostForm("note")

	ses := auth.GlobalSessions.SessionStart(c.Writer, c.Request)
	user, _ := ses.Get("user").(models.User)
	ins, _ := models.GetInstanceById(objId)
	ins.Record(note, user.Id, result)

	if result == 1 {
		ins.Next()
	} else {
		completeTime := time.Now()
		stmt, _ := models.MySqlDB.Prepare(`UPDATE instance SET complete_time=? WHERE id=?`)
		stmt.Exec(completeTime, ins.Id)
	}

	// 发送信号，处理一些事情，发送通知邮件
	singals.S <- syscall.SIGUSR2

	c.JSON(http.StatusOK, gin.H{"msg": 1})
}
