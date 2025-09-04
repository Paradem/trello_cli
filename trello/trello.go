package trello

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const baseURL = "https://api.trello.com/1"

type Client struct {
	apiKey   string
	apiToken string
	client   *http.Client
}

type Card struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	IDMembers []string `json:"idMembers"`
	ShortLink string   `json:"shortLink"`
	IDShort   int      `json:"idShort"`
	IDList    string   `json:"idList"`
}

type Board struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type List struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	ID   string `json:"id"`
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
	Date          string `json:"date"`
	MemberCreator struct {
		FullName string `json:"fullName"`
		Username string `json:"username"`
	} `json:"memberCreator"`
}

type DetailedCard struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Desc             string   `json:"desc"`
	IDMembers        []string `json:"idMembers"`
	ShortLink        string   `json:"shortLink"`
	IDShort          int      `json:"idShort"`
	IDList           string   `json:"idList"`
	Closed           bool     `json:"closed"`
	DateLastActivity string   `json:"dateLastActivity"`
	Labels           []struct {
		Name string `json:"name"`
	} `json:"labels"`
}

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"displayName"`
}

func NewClient(apiKey, apiToken string) *Client {
	return &Client{
		apiKey:   apiKey,
		apiToken: apiToken,
		client:   &http.Client{},
	}
}

func (c *Client) makeRequest(method, endpoint string, params map[string]string) (*http.Response, error) {
	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("token", c.apiToken)

	for k, v := range params {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *Client) GetCards(boardID string) ([]Card, error) {
	resp, err := c.makeRequest("GET", fmt.Sprintf("/boards/%s/cards", boardID), map[string]string{
		"members": "true",
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cards []Card
	if err := json.Unmarshal(body, &cards); err != nil {
		return nil, err
	}

	return cards, nil
}

func (c *Client) GetMemberID() (string, error) {
	resp, err := c.makeRequest("GET", "/members/me", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var member struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &member); err != nil {
		return "", err
	}

	return member.ID, nil
}

func (c *Client) GetOrganizations() ([]Organization, error) {
	resp, err := c.makeRequest("GET", "/members/me/organizations", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var organizations []Organization
	if err := json.Unmarshal(body, &organizations); err != nil {
		return nil, err
	}

	return organizations, nil
}

func (c *Client) GetBoards(organizationID string) ([]Board, error) {
	resp, err := c.makeRequest("GET", fmt.Sprintf("/organizations/%s/boards", organizationID), map[string]string{
		"filter": "open",
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var boards []Board
	if err := json.Unmarshal(body, &boards); err != nil {
		return nil, err
	}

	return boards, nil
}

func (c *Client) GetLists(boardID string) ([]List, error) {
	resp, err := c.makeRequest("GET", fmt.Sprintf("/boards/%s/lists", boardID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var lists []List
	if err := json.Unmarshal(body, &lists); err != nil {
		return nil, err
	}

	return lists, nil
}

func (c *Client) GetCardDetails(cardID string) (*DetailedCard, error) {
	resp, err := c.makeRequest("GET", fmt.Sprintf("/cards/%s", cardID), map[string]string{
		"members": "true",
		"labels":  "true",
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var card DetailedCard
	if err := json.Unmarshal(body, &card); err != nil {
		return nil, err
	}

	return &card, nil
}

func (c *Client) GetCardComments(cardID string) ([]Comment, error) {
	resp, err := c.makeRequest("GET", fmt.Sprintf("/cards/%s/actions", cardID), map[string]string{
		"filter": "commentCard",
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var comments []Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}
