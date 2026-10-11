package jenkins

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Client struct {
	BaseURL  string
	Username string
	APIToken string
	HTTP     *http.Client
}

func NewClientFromEnv(baseURL string) (*Client, error) {
	user := os.Getenv("JENKINS_USER")
	token := os.Getenv("JENKINS_API_TOKEN")

	if user == "" || token == "" {
		return nil, fmt.Errorf("JENKINS_USER or JENKINS_API_TOKEN is not set")
	}

	// リダイレクトを自動追従せず、その場でステータスを返すように設定
	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Client{
		BaseURL:  baseURL,
		Username: user,
		APIToken: token,
		HTTP:     httpClient,
	}, nil
}

// TriggerBuild 指定したジョブのビルドを実行する
func (c *Client) TriggerBuild(jobName string) error {
	endpoint := fmt.Sprintf("%s/job/%s/buildWithParameters", c.BaseURL, jobName)

	// 固定パラメータ BRANCH_NAME=main を設定
	val := url.Values{}
	val.Set("BRANCH_NAME", "main")
	reqBody := strings.NewReader(val.Encode())

	req, err := http.NewRequest("POST", endpoint, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Username, c.APIToken)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("failed to build '%s': status code %d", jobName, resp.StatusCode)
}