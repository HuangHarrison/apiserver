package user

import (
	"fmt"

	. "apiserver/handler"
	"apiserver/pkg/errno"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// Create creates a new user account.
func Create(c *gin.Context) {
	var r CreateRequest
	// 检查 Content-Type 类型，将消息体作为指定的格式解析到 Go struct 变量中
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	// 返回 URL 的参数值
	admin2 := c.Param("username")
	log.Infof("URL username: %s", admin2)

	// 读取 URL 中的地址参数
	desc := c.Query("desc")
	log.Infof("URL key param desc: %s", desc)

	// 获取 HTTP 头
	contentType := c.GetHeader("Content-Type")
	log.Infof("Header Content-Type: %s", contentType)

	log.Debugf("username is: [%s], password is [%s]", r.Username, r.Password)
	if r.Username == "" {
		SendResponse(c, errno.New(errno.ErrUserNotFound, fmt.Errorf("username can not found in db: xx.xx.xx.xx")), nil)
		return
	}

	if r.Password == "" {
		SendResponse(c, fmt.Errorf("password is empty"), nil)
	}

	rsp := CreateResponse{
		Username: r.Username,
	}

	// Show the user information.
	SendResponse(c, nil, rsp)
}
