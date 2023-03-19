package openaiproxy

import (
	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"context"
	"log"
	"net/http"
	"fmt"
	"openaiproxy/common"
)

type OpenAIRequestBody struct {
	Messages []openai.ChatCompletionMessage `json:"messages"`
	MaxTokens int `json:"maxTokens"`
}

type OpenAIProxyController struct {
	Key string
}

func (opc *OpenAIProxyController)openAIV1ChatCompletions(c *gin.Context){
	//获取请求体中携带的消息
	var req OpenAIRequestBody
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return 
  }
	
	//调用openai接口
	client := openai.NewClient(opc.Key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: req.Messages,
			MaxTokens: req.MaxTokens,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultCreateOpenAiChatCompletionError,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
	
	//返回结果
	rsp:=common.CreateResponse(nil,resp)
	c.IndentedJSON(http.StatusOK, rsp)
}

func (opc *OpenAIProxyController) Bind(router *gin.Engine) {
	log.Println("Bind UserController")
	router.POST("/openai/v1/chat/completions/GPT3Dot5Turbo", opc.openAIV1ChatCompletions)
}
