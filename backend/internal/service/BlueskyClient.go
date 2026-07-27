package service
import(
	"context"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)
const DefaultPDS="https://bsky.social"

type BskyClient struct{
	HTTPClient *http.Client
	PDSURL string
}
func NewBskyClient() *BskyClient{
	return &BskyClient{
		HTTPClient: &http.Client{
			Timeout: 15*time.Second,
		},
		PDSURL: DefaultPDS,
	}
}

func(c *BskyClient)CreateSession(ctx context.Context,handle,password string)(*Session,error){
body,_:=json.Marshal(map[string]string{
	"identifier":handle,
	"password":password,
})
req,err:=http.NewRequestWithContext(ctx,"POST",c.PDSURL+"/xrpc/com.atproto.server.createSession",bytes.NewReader(body))
if err!=nil{
	return nil,err
}
req.Header.Set("Content-Type","application/json")
resp,err:=c.HTTPClient.Do(req)
if err!=nil{
	return nil,err
}

defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
        var errResp map[string]any
        _ = json.NewDecoder(resp.Body).Decode(&errResp)
        return nil, fmt.Errorf("bluesky auth failed (%d): %v", resp.StatusCode, errResp)
    }
	var sess Session
	if err:=json.NewDecoder(resp.Body).Decode(&sess); err!=nil{
		return nil,err

	}
	return &sess,nil
}

type Session struct{
	AccessJwt *string `json:"accessJwt"`
	RefreshJwt *string `json:"refreshJwt"`
	Handle *string `json:"handle"`
	DID string `json:"did"`
	Email string `json:"email,omitemtpy"`
}
