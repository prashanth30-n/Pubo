package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const BaseURl = "https://api.linkedin.com"
const LinkedinVersion = "202606"

type LinkedinClient struct{ HTTPClient *http.Client }

func NewLinkedinClient() *LinkedinClient {
	return &LinkedinClient{HTTPClient: &http.Client{Timeout: 15 * time.Second}}
}

type UserInfo struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func (c *LinkedinClient) UserInfo(ctx context.Context, token string) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", BaseURl+"/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("linkedin userinfo failed (%d)", resp.StatusCode)
	}
	var user UserInfo
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
func (c *LinkedinClient) CreatePost(ctx context.Context, token, author, commentary string, images [][]byte, types []string) (string, error) {
	urns := []string{}
	for i, data := range images {
		b, _ := json.Marshal(map[string]any{"initializeUploadRequest": map[string]any{"owner": author}})
		req, err := http.NewRequestWithContext(ctx, "POST", BaseURl+"/rest/images?action=initializeUpload", bytes.NewReader(b))
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Linkedin-Version", LinkedinVersion)
		req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return "", err
		}
		var init struct {
			Value struct {
				UploadURL string `json:"uploadUrl"`
				Image     string `json:"image"`
			} `json:"value"`
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			resp.Body.Close()
			return "", fmt.Errorf("linkedin image initialization failed (%d)", resp.StatusCode)
		}
		err = json.NewDecoder(resp.Body).Decode(&init)
		resp.Body.Close()
		if err != nil {
			return "", err
		}
		put, err := http.NewRequestWithContext(ctx, "PUT", init.Value.UploadURL, bytes.NewReader(data))
		if err != nil {
			return "", err
		}
		put.Header.Set("Authorization", "Bearer "+token)
		put.Header.Set("Content-Type", types[i])
		putResp, err := c.HTTPClient.Do(put)
		if err != nil {
			return "", err
		}
		putResp.Body.Close()
		if putResp.StatusCode < 200 || putResp.StatusCode > 299 {
			return "", fmt.Errorf("linkedin image upload failed (%d)", putResp.StatusCode)
		}
		urns = append(urns, init.Value.Image)
	}
	payload := map[string]any{"author": author, "commentary": commentary, "visibility": "PUBLIC", "lifecycleState": "PUBLISHED", "isReshareDisabledByAuthor": false, "distribution": map[string]any{"feedDistribution": "MAIN_FEED", "targetEntities": []any{}, "thirdPartyDistributionChannels": []any{}}}
	if len(urns) > 0 {
		payload["content"] = map[string]any{"media": urns}
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", BaseURl+"/rest/posts", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Linkedin-Version", LinkedinVersion)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("linkedin post failed (%d)", resp.StatusCode)
	}
	return resp.Header.Get("x-restli-id"), nil
}
