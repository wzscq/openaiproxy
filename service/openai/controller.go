package openai

import (
	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"context"
	"log"
	"net/http"
	"fmt"
	"io"
	"openaiproxy/common"
	"errors"
	"time"
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

	fmt.Println("call GPT3Dot5Turbo success.")
	//fmt.Println(resp.Choices[0].Message.Content)
	
	//返回结果
	rsp:=common.CreateResponse(nil,resp)
	c.IndentedJSON(http.StatusOK, rsp)
}

func (opc *OpenAIProxyController)openAIV2ChatCompletions(c *gin.Context){
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
			Model: openai.GPT4,
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

	fmt.Println("call GPT4 success.")
	//fmt.Println(resp.Choices[0].Message.Content)
	
	//返回结果
	rsp:=common.CreateResponse(nil,resp)
	c.IndentedJSON(http.StatusOK, rsp)
}

func (opc *OpenAIProxyController)openAIChatStreamGPT4(c *gin.Context){
	log.Println("openAIChatStreamGPT4 start...")
	w := c.Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

  flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("server not support") //浏览器不兼容
		return
	}

	//获取请求体中携带的消息
	var req OpenAIRequestBody
	if err := c.BindJSON(&req); err != nil {
		fmt.Fprintf(w, "data: 参数错误\n\n")
		fmt.Fprintf(w, "event: close")
		return 
  }

	fmt.Fprintf(w, "data: "+"test"+"\n\n");

	sleepTime := 2
	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "data: "+"test");
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
	fmt.Fprintf(w, "event: close")

	return

	//调用openai接口
	client := openai.NewClient(opc.Key)
	ctx := context.Background()
	gptReq := openai.ChatCompletionRequest{
		Model:     openai.GPT4,
		Messages: req.Messages,
		MaxTokens: req.MaxTokens,
		Stream: true,
	}
	log.Println("openAIChatStreamGPT4 CreateChatCompletionStream")
	stream, err := client.CreateChatCompletionStream(ctx, gptReq)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	log.Println("openAIChatStreamGPT4 receive response")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		log.Println(response.Choices[0].Delta.Content)
		fmt.Fprintf(w, "data: "+response.Choices[0].Delta.Content+"\n\n")
		flusher.Flush()
	}

	fmt.Fprintf(w, "event: close")
}

func (opc *OpenAIProxyController) Bind(router *gin.Engine) {
	router.POST("/openai/v1/chat/completions/GPT3Dot5Turbo", opc.openAIV1ChatCompletions)
	router.POST("/openai/v1/chat/completions/GPT4", opc.openAIV2ChatCompletions)
	router.POST("/openai/chat/stream/GPT4", opc.openAIChatStreamGPT4)
}
