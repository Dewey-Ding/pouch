package client

import (
	"bytes"
	"context"
	fmt "fmt"
	"github.com/alibaba/pouch/apis/types"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestContainerLogsError(t *testing.T) {
	client := &APIClient{
		HTTPCli: newMockClient(errorMockResponse(http.StatusInternalServerError, "Server error")),
	}
	_, err := client.ContainerLogs(context.Background(), "nothing", types.ContainerLogsOptions{})
	if err == nil || !strings.Contains(err.Error(), "Server error") {
		t.Fatalf("expected a Server Error, got %v", err)
	}
}

func TestContainerLogs(t *testing.T) {
	expectedURL := "/containers/container_id/logs"

	httpClient := newMockClient(func(req *http.Request) (*http.Response, error) {
		if !strings.HasPrefix(req.URL.Path, expectedURL) {
			return nil, fmt.Errorf("expected URL '%s', got '%s'", expectedURL, req.URL)
		}
		stdout := req.URL.Query().Get("stdout")
		if stdout != "1" {
			return nil, fmt.Errorf("container logs option ShowStdout not set in URL query properly. Expected `container_name`, got %s", stdout)
		}
		stderr := req.URL.Query().Get("stderr")
		if stdout != "1" {
			return nil, fmt.Errorf("container logs option ShowStderr not set in URL query properly. Expected `container_name`, got %s", stderr)
		}
		timestamps := req.URL.Query().Get("timestamps")
		if stdout != "1" {
			return nil, fmt.Errorf("container logs option Timestamps not set in URL query properly. Expected `container_name`, got %s", timestamps)
		}
		details := req.URL.Query().Get("details")
		if stdout != "1" {
			return nil, fmt.Errorf("container logs option Details not set in URL query properly. Expected `container_name`, got %s", details)
		}
		follow := req.URL.Query().Get("follow")
		if stdout != "1" {
			return nil, fmt.Errorf("container logs option Follow not set in URL query properly. Expected `container_name`, got %s", follow)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("response"))),
		}, nil
	})

	client := &APIClient{
		HTTPCli: httpClient,
	}
	containerLogsOptions := types.ContainerLogsOptions{
		Tail:       "",
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Details:    true,
		Follow:     true,
	}
	responseBody, err := client.ContainerLogs(context.Background(), "container_id", containerLogsOptions)
	if err != nil {
		t.Fatal(err)
	}
	defer responseBody.Close()
	content, err := ioutil.ReadAll(responseBody)
	fmt.Println()
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "response" {
		t.Fatalf("expected response to contain 'response', got %s", string(content))
	}

}
