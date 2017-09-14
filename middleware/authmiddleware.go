package middleware

import (
	"github.com/astaxie/session"
	_ "github.com/astaxie/session/providers/memory"
	"github.com/gin-gonic/gin"
	"net/http"
)

var GlobalSessions *session.Manager

func init() {
	GlobalSessions, _ = session.NewManager("memory", "authenticate", 3600*24)
	go GlobalSessions.GC()
}

func Authenticate() gin.HandlerFunc {
	return AuthWithWriter(GlobalSessions)
}

func AuthWithWriter(globalSessions *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		ses := globalSessions.SessionStart(c.Writer, c.Request)
		if auth, ok := ses.Get("authenticated").(bool); !ok || !auth {
			c.Redirect(http.StatusFound, "http://192.168.162.108:8005")
			return
		} else {
			//fmt.Println("wwwwwwwwww")    //process_request  处理前打印
			c.Next()                         //process_request   加不加都会处理请求，只是执行顺序
			//fmt.Println("eeeeeeeeeeeeeeee")    //process_request  处理后打印
		}
	}
}