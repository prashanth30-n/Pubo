package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultPDS = "https://bsky.social"

type BskyClient struct {
	HTTPClient *http.Client
	PDSURL     string
}

func NewBskyClient() *BskyClient {
	return &BskyClient{
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		PDSURL: DefaultPDS,
	}
}

func (c *BskyClient) CreateSession(ctx context.Context, handle, password string) (*Session, error) {
	body, _ := json.Marshal(map[string]string{
		"identifier": handle,
		"password":   password,
	})
	req, err := http.NewRequestWithContext(ctx, "POST", c.PDSURL+"/xrpc/com.atproto.server.createSession", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("bluesky auth failed (%d): %v", resp.StatusCode, errResp)
	}
	var sess Session
	if err := json.NewDecoder(resp.Body).Decode(&sess); err != nil {
		return nil, err

	}
	return &sess, nil
}

// func (c *BskyClient) RefreshSession(ctx context.Context, refreshJwt string) (*Session, error) {
//     req, _ := http.NewRequestWithContext(ctx, "POST",
//         c.PDSURL+"/xrpc/com.atproto.server.refreshSession", nil)
//     req.Header.Set("Authorization", "Bearer "+refreshJwt)

//     resp, err := c.HTTPClient.Do(req)
//     if err==nil{
// 		return resp,nil
// 	}
// }

type Session struct {
	AccessJwt  *string `json:"accessJwt"`
	RefreshJwt *string `json:"refreshJwt"`
	Handle     *string `json:"handle"`
	DID        string  `json:"did"`
	Email      string  `json:"email,omitemtpy"`
}
type BlueskyImage struct {
	Alt   string         `json:"alt"`
	Image map[string]any `json:"image"`
}

func (c *BskyClient) UploadBlob(ctx context.Context, token string, data []byte, contentType string) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.PDSURL+"/xrpc/com.atproto.repo.uploadBlob", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", contentType)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("bluesky image upload failed (%d): %s", resp.StatusCode, string(body))
	}
	var out struct {
		Blob map[string]any `json:"blob"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Blob, nil
}
func (c *BskyClient) CreatePost(ctx context.Context, token, did, text string, images []BlueskyImage) (string, error) {
	record := map[string]any{"$type": "app.bsky.feed.post", "text": text, "createdAt": time.Now().UTC().Format(time.RFC3339)}
	if len(images) > 0 {
		record["embed"] = map[string]any{"$type": "app.bsky.embed.images", "images": images}
	}
	body, _ := json.Marshal(map[string]any{"repo": did, "collection": "app.bsky.feed.post", "record": record})
	req, err := http.NewRequestWithContext(ctx, "POST", c.PDSURL+"/xrpc/com.atproto.repo.createRecord", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return "", fmt.Errorf("bluesky post failed (%d): %s", resp.StatusCode, string(body))
	}
	var out struct {
		URI string `json:"uri"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.URI, nil
}
