package training

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/pinlan/app-admin/internal/app/app-admin/schema"
	bll "github.com/pinlan/app-admin/internal/app/app-admin/training"
)

// TrainingServer
type TrainingServer struct {
	Server *bll.Server
}

// NewTrainingServer
func NewTrainingServer(server *bll.Server) *TrainingServer {
	return &TrainingServer{Server: server}
}

// Create
func (cli *TrainingServer) Create(c *gin.Context) {
	user_id := c.GetHeader("UserID")
	training_schema := schema.TrainingSchema{}

	if err := c.Bind(&training_schema); err != nil {
		c.JSON(200, schema.GenerateApiData("Body 解析失败", 400, nil))
		return
	}

	fmt.Println(training_schema)
	training, err := cli.Server.Create(user_id, training_schema)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			c.JSON(200, schema.GenerateApiData("任务已经存在", 400, nil))
			return

		} else {
			c.JSON(200, schema.GenerateApiData("任务创建失败", 400, nil))
			return
		}
	}
	c.JSON(200, schema.GenerateApiData("任务创建成功", 200, training))
}

// Delete
func (cli *TrainingServer) Delete(c *gin.Context) {
	id := c.Params.ByName("id")
	err := cli.Server.Delete(id)
	if err != nil {
		c.JSON(200, schema.GenerateApiData("删除任务失败", 400, nil))
		return
	}
	c.JSON(200, schema.GenerateApiData("删除任务成功", 200, nil))
}

// Get
func (cli *TrainingServer) Get(c *gin.Context) {
	id := c.Params.ByName("id")
	training, err := cli.Server.Get(id)
	if err != nil {
		c.JSON(200, schema.GenerateApiData("获取任务详情失败", 400, nil))
		return
	}
	c.JSON(200, schema.GenerateApiData("获取任务成功", 200, training))
}

// List
func (cli *TrainingServer) List(c *gin.Context) {
	user_id := c.GetHeader("UserID")
	training_list, err := cli.Server.List(user_id)
	if err != nil {
		c.JSON(200, schema.GenerateApiData("获取任务列表失败", 400, nil))
		return
	}
	c.JSON(200, schema.GenerateApiData("", 200, training_list))
}

// GetLogs
func (cli *TrainingServer) GetLogs(c *gin.Context) {
	id := c.Params.ByName("id")
	log_details, err := cli.Server.GetLogs(id)
	if err != nil {
		c.JSON(200, schema.GenerateApiData("获取任务日志失败", 400, nil))
		return
	}
	c.JSON(200, schema.GenerateApiData("", 200, log_details))
}
