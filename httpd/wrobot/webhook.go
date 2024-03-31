package wrobot

import (
	"github.com/gin-gonic/gin"
	"github.com/opentdp/wrest-chat/dbase/webhook"
	"github.com/opentdp/wrest-chat/wcferry"
	"github.com/opentdp/wrest-chat/wclient/webhookApp"
	"io"
)

type Webhook struct {
	*wcferry.Client
}

// @Summary webhook列表
// @Produce json
// @Tags BOT::webhook
// @Success 200 {array} tables.Webhook
// @Router /bot/webhook/list [post]
func (*Webhook) list(c *gin.Context) {

	if lst, err := webhook.FetchAll(); err == nil {
		c.Set("Payload", lst)
	} else {
		c.Set("Error", err)
	}

}

// @Summary 获取webhook
// @Produce json
// @Tags BOT::webhook
// @Param body body webhook.FetchWebHookParam true "获取webhook参数"
// @Success 200 {object} tables.Webhook
// @Router /bot/webhook/detail [post]
func (*Webhook) detail(c *gin.Context) {

	var rq *webhook.FetchWebHookParam

	if err := c.ShouldBind(&rq); err != nil {
		c.Set("Error", err)
		return
	}

	if res, err := webhook.Fetch(rq); err == nil {
		c.Set("Payload", res)
	} else {
		c.Set("Error", err)
	}

}

// @Summary 添加webhook
// @Produce json
// @Tags BOT::webhook
// @Param body body webhook.CreateWebHookParam true "添加webhook参数"
// @Success 200
// @Router /bot/webhook/create [post]
func (*Webhook) create(c *gin.Context) {

	var rq *webhook.CreateWebHookParam

	if err := c.ShouldBind(&rq); err != nil {
		c.Set("Error", err)
		return
	}

	if id, _, err := webhook.Create(rq); err == nil {
		c.Set("Message", "添加成功")
		c.Set("Payload", id)
	} else {
		c.Set("Error", err)
	}

}

// @Summary 删除webhook
// @Produce json
// @Tags BOT::webhook
// @Param body body webhook.DeleteWebHookParam true "删除webhook"
// @Success 200
// @Router /bot/webhook/delete [post]
func (*Webhook) delete(c *gin.Context) {

	var rq *webhook.DeleteWebHookParam

	if err := c.ShouldBind(&rq); err != nil {
		c.Set("Error", err)
		return
	}

	if err := webhook.Delete(rq); err == nil {
		c.Set("Message", "删除成功")
	} else {
		c.Set("Error", err)
	}

}

// @Summary 接收webhook消息
// @Produce json
// @Tags BOT::webhook
// @Param token path string true "webhook token"
// @Param app path string true "webhook类型(例如： github, gitea)"
// @Param body body interface{} true "event报文"
// @Success 200
// @Router /bot/webhook/{token}/{app} [post]
func (w *Webhook) receive(c *gin.Context) {
	token := c.Param("token")
	app := c.Param("app")

	hook, err := webhook.FindByToken(token)

	if err != nil {
		c.Set("Error", err)
		return
	}

	request, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.Set("Error", err)
		return
	}

	msg := string(request)

	// 根据app类型不同，调用不同的处理方式，参照handler的注册
	sendMsg := webhookApp.Handler(c.Request.Header, app, msg)
	res := w.CmdClient.SendTxt(sendMsg, hook.TargetId, "")

	if res == 0 {
		c.Set("Message", "OK")
	} else {
		c.Set("Error", "消息处理失败")
	}
}
