package crv

import (
	"log"
	"encoding/json"
	"bytes"
	"net/http"
	"openaiproxy/common"
)

type loginRep struct {
    UserID     string  `json:"userID"`
    Password  string   `json:"password"`
	AppID     string   `json:"appID"`
} 

type loginResult struct {
    UserID     *string  `json:"userID"`
    UserName  *string  `json:"userName"`
	Token     *string  `json:"token"`
	AppID     *string  `json:"appID"`
}

type loginRsp struct {
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Error bool `json:"error"`
	Result loginResult `json:"result"`
}

type CommonHeader struct {
	Token     string  `json:"token"`
	UserID    string  `json:"userID"`
	AppDB     string  `json:"appDB"`
	UserRoles string  `json:"userRoles"`
}

type CommonRsp struct {
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Error bool `json:"error"`
	Result map[string]interface{} `json:"result"`
}

type Pagination struct {
    Current int `json:"current"` 
    PageSize int `json:"pageSize"` 
}

type Sorter struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

type CommonReq struct {
	ModelID string `json:"modelID"`
	ViewID *string `json:"viewID"`
	Filter *map[string]interface{} `json:"filter"`
	List *[]map[string]interface{} `json:"list"`
	Fields *[]map[string]interface{} `json:"fields"`
	UserID *string `json:"userID"`
	AppDB *string `json:"appDB"`
	UserRoles *string `json:"userRoles"`
	Sorter *[]Sorter `json:"sorter"`
	SelectedRowKeys *[]string `json:"selectedRowKeys"`
	Pagination *Pagination `json:"pagination"`
}

type CRVClient struct {
	Server string `json:"server"`
    User string `json:"user"`
    Password string `json:"password"`
    AppID string `json:"appID"`
	Token string `json:"token"`
}

const (
	URL_LOGIN = "/user/login"
	URL_SAVE = "/data/save"
	URL_QUERY = "/data/query"
)

const (
	SAVE_TYPE_COLUMN = "_save_type"
	SAVE_CREATE = "create"
	SAVE_UPDATE = "update"
	SAVE_DELETE = "delete"

	CC_CREATE_TIME = "create_time"
	CC_CREATE_USER = "create_user"
	CC_UPDATE_TIME = "update_time"
	CC_UPDATE_USER = "update_user"
	CC_VERSION = "version"
	CC_ID = "id"
)

func (crv *CRVClient) Login()(int) {
	log.Println("start login")
	loginRep:=loginRep{
		UserID:crv.User,
		Password:crv.Password,
		AppID:crv.AppID,
	}

	postJson,_:=json.Marshal(loginRep)
	postBody:=bytes.NewBuffer(postJson)
	resp,err:=http.Post(crv.Server+URL_LOGIN,"application/json",postBody)

	if err != nil {
		log.Println("login error",err)
		return -1
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != 200 { 
		log.Println("login error",resp)
		return -1
	}

	decoder := json.NewDecoder(resp.Body)
	loginRsp:=loginRsp{}
	err = decoder.Decode(&loginRsp)
	if err != nil {
		log.Println("result decode failed [Err:%s]", err.Error())
		return -1
	}

	if loginRsp.Result.Token == nil {
		log.Println("login error ",loginRsp)
		return -1
	}

	crv.Token=*loginRsp.Result.Token
	log.Println("login success")
	return 0
}

func (crv *CRVClient)Save(commonReq *CommonReq,token string)(*CommonRsp,int){
	log.Println("start CRVClient save ...")
	postJson,_:=json.Marshal(*commonReq)
	postBody:=bytes.NewBuffer(postJson)
	req,err:=http.NewRequest("POST",crv.Server+URL_SAVE,postBody)
	if err != nil {
		log.Println("CRVClient save NewRequest error",err)
		return nil,common.ResultSaveDataError
	}

	if len(token)==0 {
		req.Header.Set("token", crv.Token)
	} else {
		req.Header.Set("token", token)
	}
	
	req.Header.Set("Content-Type","application/json")
	
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("CRVClient save Do request error",err)
		return nil,common.ResultSaveDataError
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 { 
		log.Println("CRVClient save StatusCode error",resp)
		return nil,common.ResultSaveDataError
	}

	decoder := json.NewDecoder(resp.Body)
	commonRsp:=CommonRsp{}
	err = decoder.Decode(&commonRsp)
	if err != nil {
		log.Println("CRVClient save result decode failed [Err:%s]", err.Error())
		return nil,common.ResultSaveDataError
	}

	if commonRsp.Error == true {
		log.Printf("errorcode:%d,message:%s \n",commonRsp.ErrorCode,commonRsp.Message)
	} else {
		resultJson,_:=json.Marshal(&commonRsp.Result)
		log.Println(string(resultJson))
	}

	log.Println("end CRVClient save success")
	return &commonRsp,common.ResultSuccess
}

func (crv *CRVClient)Query(commonReq *CommonReq,token string)(*CommonRsp,int){
	log.Println("start CRVClient query ...")
	postJson,_:=json.Marshal(*commonReq)
	postBody:=bytes.NewBuffer(postJson)
	req,err:=http.NewRequest("POST",crv.Server+URL_QUERY,postBody)
	if err != nil {
		log.Println("CRVClient query NewRequest error",err)
		return nil,common.ResultQueryRequestError
	}

	if len(token)==0 {
		req.Header.Set("token", crv.Token)
	} else {
		req.Header.Set("token", token)
	}

	req.Header.Set("Content-Type","application/json")
	
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("CRVClient query Do request error",err)
		return nil,common.ResultQueryRequestError
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 { 
		log.Println("CRVClient query StatusCode error",resp)
		return nil,common.ResultQueryRequestError
	}

	decoder := json.NewDecoder(resp.Body)
	commonRsp:=CommonRsp{}
	err = decoder.Decode(&commonRsp)
	if err != nil {
		log.Println("CRVClient query result decode failed [Err:%s]", err.Error())
		return nil,common.ResultQueryRequestError
	}

	resultJson,_:=json.Marshal(&commonRsp.Result)
	log.Println(string(resultJson))

	log.Println("end CRVClient query success")
	return &commonRsp,common.ResultSuccess
}



