package jenkins

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
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

	// SSL/TLS証明書エラー対策を含めたClientを初期化
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Client{
		BaseURL:  baseURL,
		Username: user,
		APIToken: token,
		HTTP:     &http.Client{Transport: tr},
	}, nil
}

// TriggerBuild 指定したジョブのビルドを実行する
func (c *Client) TriggerBuild(jobName string) error {
	endpoint := fmt.Sprintf("%s/job/%s/build", c.BaseURL, jobName)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}

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