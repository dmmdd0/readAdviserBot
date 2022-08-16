package telegram

// test GH
import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"readAdviserBot/lib/e"
	"strconv"
)

type Client struct {
	host     string
	basePath string // tg-bot.com/<token>
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "semdMassege"
)

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

// get update from tg
// Updates but not GetUpdates
func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	// Itoa == Integer to Ascii
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	// do request <- getUpdates
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	// stoped lesson 3 17:45
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}
	return nil
}

//
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	//const errMsg = "can't do request"
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		// Path:   c.basePath + "/" + method,
		// + is bad because we will get "" or "/" or "//" or "///"
		Path: path.Join(c.basePath, method),
	}
	//req, err := http.NewRequest("GET") equival
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		//return nil, err // err from http package
		return nil, err
	}
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	// lesson 3 15:00
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}
