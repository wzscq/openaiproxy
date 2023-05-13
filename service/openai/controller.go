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
	"openaiproxy/crv"
	"errors"
	"strings"
	"strconv"
	//"time"
)

type OpenAIRequestBody struct {
	Messages []openai.ChatCompletionMessage `json:"messages"`
	MaxTokens int `json:"maxTokens"`
	Account string `json:"account"`
	Password string `json:"password"`
	Model string `json:"model"`
}

type OpenAIProxyController struct {
	Key string
	CRVClient *crv.CRVClient
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
	model:=openai.GPT3Dot5Turbo
	if req.Model=="GPT4" {
		model=openai.GPT4
	}
	
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
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

func (opc *OpenAIProxyController)updateAcountInfo(account,password string)(bool){
	accountFields:=[]map[string]interface{}{
		{"field": "id"},
		{"field": "count"},
		{"field": "version"},
	}

	commonRep:=crv.CommonReq{
		ModelID:"gpt_account",
		Fields:&accountFields,
		Filter:&map[string]interface{}{
			"id":account,
			"password":password,
		},
	}

	rsp,commonErr:=opc.CRVClient.Query(&commonRep,"")
	if commonErr!=common.ResultSuccess {
		return false
	}

	if rsp.Error == true {
		log.Println("query account info error:",rsp.ErrorCode,rsp.Message)
		return false
	}

	resLst,ok:=rsp.Result["list"].([]interface{})
	if !ok {
		log.Println("query account info error: no list data.")
		return false
	}

	if len(resLst)==0 {
		log.Println("query account info error: no account data.")
		return false
	}

	accountInfo,ok:=resLst[0].(map[string]interface{})
	if !ok {
		log.Println("query account info error: no account data.")
		return false
	}

	//更新账号次数
	strCount:=accountInfo["count"].(string)
	count,err:=strconv.Atoi(strCount)
	if err!=nil {
		log.Println("query account info error: count is not int.")
		return false
	}

	count=count+1
	strCount=strconv.Itoa(count)

	accountInfo["count"]=strCount
	accountInfo["_save_type"]="update"

	recList:=[]map[string]interface{}{accountInfo}
	saveReq:=&crv.CommonReq{
		ModelID:"gpt_account",
		List:&recList,
	}
	opc.CRVClient.Save(saveReq,"")
	
	return true
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

	//检查用户账号并更新次数信息
	if opc.updateAcountInfo(req.Account,req.Password)==false {
		fmt.Fprintf(w, "data: 您的账号或密码可能不正确，请确认是否输入了正确的账号密码，如有疑问请与管理员联系处理\n\n")
		fmt.Fprintf(w, "event: close")
		return
	}

	//fmt.Fprintf(w, "data: "+"test"+"\n\n");

	/*sleepTime := 2
	for i := 0; i < 10; i++ {
		//在log中输出变量i的值
		log.Println("i: ", i)
		log.Println("data: "+"test")
		fmt.Fprintf(w, "data: "+"test\n\n");
		flusher.Flush()
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
	fmt.Fprintf(w, "event: close")

	return*/

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
		fmt.Fprintf(w, "data: 访问GPT4接口出错，请稍后重试，或联系管理员处理\n\n")
		fmt.Fprintf(w, "event: close")
		return
	}
	defer stream.Close()

	log.Println("openAIChatStreamGPT4 receive response")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			break;
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			break;
		}

		contentStr:=response.Choices[0].Delta.Content
		//将内容中的\n替换为<br/>
		contentStr=strings.Replace(contentStr,"\n","<br/>",-1)
		log.Println(contentStr)
		fmt.Fprintf(w, "data: "+contentStr+"\n\n")
		flusher.Flush()
	}

	fmt.Fprintf(w, "event: close")
}

func (opc *OpenAIProxyController) Bind(router *gin.Engine) {
	router.POST("/openai/v1/chat/completions/GPT3Dot5Turbo", opc.openAIV1ChatCompletions)
	router.POST("/openai/v1/chat/completions/GPT4", opc.openAIV2ChatCompletions)
	router.POST("/openai/chat/stream/GPT4", opc.openAIChatStreamGPT4)
}
