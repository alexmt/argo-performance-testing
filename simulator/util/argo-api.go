package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

type ArgoApi struct {
	host string
}

func NewArgoApi(host string) *ArgoApi {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return &ArgoApi{host}
}

func (api *ArgoApi) sendPost(uri string, token string, payload []byte) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", api.host+uri, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := make(map[string]interface{})
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (api *ArgoApi) sendGet(uri string, token string, result interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", api.host+uri, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	return nil
}

func (api *ArgoApi) GetToken(username, password string) (string, error) {
	payload := make(map[string]interface{})
	payload["username"] = username
	payload["password"] = password
	postBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	result, err := api.sendPost("/api/v1/session", "", postBody)
	if err != nil {
		return "", err
	}
	return result["token"].(string), nil
}

func (api *ArgoApi) ListApplications(token string) ([]v1alpha1.Application, error) {
	var apps v1alpha1.ApplicationList
	err := api.sendGet("/api/v1/applications", token, &apps)
	if err != nil {
		return nil, err
	}
	return apps.Items, nil
}

func (api *ArgoApi) UpdateReplicas(token string, replicasCount int, applicationName string) (map[string]interface{}, error) {
	data := fmt.Sprintf("\"{\\\"spec\\\":{\\\"replicas\\\":%d}}\"", replicasCount)
	uri := strings.ReplaceAll("/api/v1/applications/%s/resource?name=%s-helm-guestbook&namespace=argocd&resourceName=%s-helm-guestbook&version=v1&kind=Deployment&group=apps&patchType=application%2Fmerge-patch%2Bjson", "%s", applicationName)
	result, err := api.sendPost(uri, token, []byte(data))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (api *ArgoApi) Sync(token string, applicationName string) (map[string]interface{}, error) {
	data := `{"revision":"master","prune":false,"dryRun":false,"strategy":{"hook":{"force":false}},"resources":null,"syncOptions":{"items":["CreateNamespace=true"]}}`
	uri := fmt.Sprintf("/api/v1/applications/%s/sync", applicationName)
	result, err := api.sendPost(uri, token, []byte(data))
	if err != nil {
		return nil, err
	}
	return result, nil
}
