package glaw

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// calls an endpoint of a Lemmy instance API
func (lc *LemmyClient) callLemmyAPI(method string, endpoint string, body io.Reader) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Prepare the request
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", lc.baseURL, endpoint), body)
	if err != nil {
		return nil, err
	}

	// Set the API token for authentication (if required)
	if lc.APIToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.APIToken))
	}
	if lc.jwtCookie != "" {
		req.Header.Add("cookie", lc.jwtCookie)
	}

	// Send the request
	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
