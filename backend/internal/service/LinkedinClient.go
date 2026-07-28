package service

import(
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)
const BaseURl="https://api.linkedin.com"

type LinkedinClient struct{
	HTTPClient *http.Client
}

func NewLinkedinClient() *LinkedinClient{
	return &LinkedinClient{
		HTTPClient: &http.Client{
			Timeout: 15*time.Second,
		},

	}
}

func(c *LinkedinClient) UserInfo(ctx context.Context,accessToken string)(*UserInfo,error){
	req,err:=http.NewRequestWithContext(ctx,"GET",BaseURl+"/v2/userinfo",nil)
	if err!=nil{
		return nil,err
	}
	req.Header.Set("Authorization","Bearer "+accessToken) //space after bearer token
	resp,err:=c.HTTPClient.Do(req)
	if err!=nil{
		return nil,err
	}
	defer resp.Body.Close()
	if resp.StatusCode!=http.StatusOK{
		var errBody map[string]any
		 _=json.NewDecoder(resp.Body).Decode(&errBody)
		return nil,fmt.Errorf("linkedin userinfo failed (%d): %v",resp.StatusCode,errBody)

	}
	var ui UserInfo
	if err:=json.NewDecoder(resp.Body).Decode(&ui); err!=nil{
		return nil,err
	}
	return &ui,nil
}


type UserInfo struct{
	Sub string `json:"sub"` //linkedin member id
	Name string `json:"name"` //name
	Email string `json:"email"` //email
	Picture string `json:"picture"` //picture url
}

