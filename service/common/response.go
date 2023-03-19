package common

type CommonRsp struct {
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Error bool `json:"error"`
	Result interface{} `json:"result"`
	Params map[string]interface{} `json:"params"`
}

type CommonError struct {
	ErrorCode int `json:"errorCode"`
	Params map[string]interface{} `json:"params"`
	Message string `json:"message"`
}

const (
	ResultSuccess = 10000000
	ResultWrongRequest = 10000001
	ResultCreateOpenAiChatCompletionError = 90100001
)

var errMsg = map[int]CommonRsp{
	ResultSuccess:CommonRsp{
		ErrorCode:ResultSuccess,
		Message:"操作成功",
		Error:false,
	},
	ResultWrongRequest:CommonRsp{
		ErrorCode:ResultWrongRequest,
		Message:"请求参数错误，请检查参数是否完整，参数格式是否正确",
		Error:true,
	},
	ResultCreateOpenAiChatCompletionError:CommonRsp{
		ErrorCode:ResultCreateOpenAiChatCompletionError,
		Message:"调用openai接口失败",
		Error:true,
	},
}

func CreateResponse(err *CommonError,result interface{})(*CommonRsp){
	if err==nil {
		commonRsp:=errMsg[ResultSuccess]
		commonRsp.Result=result
		return &commonRsp
	}

	commonRsp,ok:=errMsg[err.ErrorCode]
	if !ok {
		commonRsp=CommonRsp{
			ErrorCode:err.ErrorCode,
			Message:err.Message,
			Error:true,
		}
	}
	commonRsp.Result=result
	commonRsp.Params=err.Params
	return &commonRsp
}

func CreateError(errorCode int,params map[string]interface{})(*CommonError){
	return &CommonError{
		ErrorCode:errorCode,
		Params:params,
	}
}