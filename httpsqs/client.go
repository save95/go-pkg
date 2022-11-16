package httpsqs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/save95/xerror"
)

type client struct {
	config *Config
}

func NewClient(config *Config) IClient {
	return &client{config: config}
}

func (c *client) Put(ctx context.Context, name, data string) (int64, error) {
	values := url.Values{}
	values.Set("name", name)
	values.Set("opt", "put")
	values.Set("data", data)

	body, headers, err := c.httpGet(ctx, values, []string{"pos"})
	if err != nil {
		return 0, xerror.Wrap(err, "httpsqs put failed")
	}

	if body == "HTTPSQS_PUT_OK" {
		pos, _ := strconv.Atoi(headers["pos"])
		return int64(pos), nil
	}

	return 0, xerror.Wrapf(err, "httpsqs put failed: %s", body)
}

func (c *client) Get(ctx context.Context, name string) (string, int64, error) {
	values := url.Values{}
	values.Set("name", name)
	values.Set("opt", "get")

	body, headers, err := c.httpGet(ctx, values, []string{"pos"})
	if err != nil {
		return "", 0, xerror.Wrap(err, "httpsqs get failed")
	}

	if body == "HTTPSQS_GET_END" {
		return "", 0, nil
	}

	pos, _ := strconv.Atoi(headers["pos"])
	return body, int64(pos), nil
}

func (c *client) Status(ctx context.Context, name string) (*Status, error) {
	values := url.Values{}
	values.Set("name", name)
	values.Set("opt", "status_json")

	body, _, err := c.httpGet(ctx, values, nil)
	if err != nil {
		return nil, xerror.Wrap(err, "httpsqs status failed")
	}

	var res Status
	if err := json.Unmarshal([]byte(body), &res); nil != err {
		return nil, xerror.Wrapf(err, "httpsqs status failed: response not json: %s", body)
	}

	return &res, nil
}

func (c *client) View(ctx context.Context, name string, pos int64) (string, error) {
	if pos <= 0 || pos > 1000000000 {
		return "", xerror.New("input pos error")
	}

	values := url.Values{}
	values.Set("name", name)
	values.Set("opt", "view")
	values.Set("pos", strconv.Itoa(int(pos)))

	body, _, err := c.httpGet(ctx, values, nil)
	if err != nil {
		return "", xerror.Wrap(err, "httpsqs view failed")
	}

	return body, nil
}

func (c *client) Reset(ctx context.Context, name string) error {
	values := url.Values{}
	values.Set("name", name)
	values.Set("opt", "reset")

	body, _, err := c.httpGet(ctx, values, nil)
	if err != nil {
		return xerror.Wrap(err, "httpsqs reset failed")
	}

	if body == "HTTPSQS_RESET_OK" {
		return nil
	}

	return xerror.New(fmt.Sprintf("httpsqs reset failed: %s", body))
}

func (c *client) SetMaxQueue(ctx context.Context, name string, max int) error {
	values := url.Values{}
	values.Set("name", name)
	values.Set("opt", "maxqueue")
	values.Set("num", strconv.Itoa(max))

	body, _, err := c.httpGet(ctx, values, nil)
	if err != nil {
		return xerror.Wrap(err, "httpsqs set max_queue failed")
	}

	if body == "HTTPSQS_MAXQUEUE_OK" {
		return nil
	}

	return xerror.New(fmt.Sprintf("httpsqs set max_queue failed: %s", body))
}

func (c *client) SetSyncTime(ctx context.Context, name string, duration time.Duration) error {
	values := url.Values{}
	values.Set("name", name)
	values.Set("opt", "synctime")
	values.Set("num", strconv.Itoa(int(duration.Seconds())))

	body, _, err := c.httpGet(ctx, values, nil)
	if err != nil {
		return xerror.Wrap(err, "httpsqs set sync_time failed")
	}

	if body == "HTTPSQS_MAXQUEUE_OK" {
		return nil
	}

	return xerror.New(fmt.Sprintf("httpsqs set max_queue failed: %s", body))
}

func (c *client) httpGet(ctx context.Context, values url.Values, parseHeaders []string) (string, map[string]string, error) {
	furl := fmt.Sprintf("http://%s/", c.config.Addr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, furl, nil)
	if nil != err {
		return "", nil, err
	}

	if len(c.config.Password) > 0 {
		values.Set("auth", c.config.Password)
	}
	req.URL.RawQuery = values.Encode()

	client := http.DefaultClient
	if c.config.Timeout > 0 {
		client.Timeout = c.config.Timeout
	} else {
		client.Timeout = 5 * time.Second
	}
	resp, err := client.Do(req)
	if nil != err {
		return "", nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	headerResult := make(map[string]string, 0)
	for _, header := range parseHeaders {
		headerResult[header] = resp.Header.Get(header)
	}

	return string(bodyBytes), headerResult, err
}
