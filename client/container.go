package client

import (
	"bufio"
	"context"
	"net"
	"net/url"

	"github.com/alibaba/pouch/apis/types"
)

// ContainerAttach attach a container
func (client *APIClient) ContainerAttach(ctx context.Context, name string, stdin bool) (net.Conn, *bufio.Reader, error) {
	q := url.Values{}
	if stdin {
		q.Set("stdin", "1")
	} else {
		q.Set("stdin", "0")
	}

	header := map[string][]string{
		"Content-Type": {"text/plain"},
	}

	return client.hijack(ctx, "/containers/"+name+"/attach", q, nil, header)
}

// ContainerCreateExec creates exec process.
func (client *APIClient) ContainerCreateExec(ctx context.Context, name string, config *types.ExecCreateConfig) (*types.ExecCreateResp, error) {
	response, err := client.post(ctx, "/containers/"+name+"/exec", url.Values{}, config, nil)
	if err != nil {
		return nil, err
	}

	body := &types.ExecCreateResp{}
	decodeBody(body, response.Body)
	ensureCloseReader(response)

	return body, nil
}

// ContainerStartExec starts exec process.
func (client *APIClient) ContainerStartExec(ctx context.Context, execid string, config *types.ExecStartConfig) (net.Conn, *bufio.Reader, error) {
	header := map[string][]string{
		"Content-Type": {"text/plain"},
	}

	return client.hijack(ctx, "/exec/"+execid+"/start", url.Values{}, config, header)
}

// ContainerUpgrade upgrade a container with new image and args.
func (client *APIClient) ContainerUpgrade(ctx context.Context, name string, config types.ContainerConfig, hostConfig *types.HostConfig) error {
	// TODO
	upgradeConfig := types.ContainerUpgradeConfig{
		ContainerConfig: config,
		HostConfig:      hostConfig,
	}
	resp, err := client.post(ctx, "/containers/"+name+"/upgrade", url.Values{}, upgradeConfig, nil)
	ensureCloseReader(resp)

	return err
}
