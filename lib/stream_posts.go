package glaw

import (
	"encoding/json"
	"time"
)

func (lc *LemmyClient) StreamNewPosts(pauseAfter int, closeChan chan struct{}) chan Post {
	// Initialize a set to track seen items
	seenItems := make(map[int]bool)
	postsChan := make(chan Post, 1000)

	go func() {
		// Initialize variables for exponential backoff
		backoff := 1 * time.Second
		maxBackoff := 16 * time.Second
		backoffReset := false
		responsesWithoutNew := 0

		for {
			// Pause mechanism
			if pauseAfter > 0 && backoffReset {
				responsesWithoutNew++
				if responsesWithoutNew > pauseAfter {
					// Reset backoff and responses count
					backoff = 1 * time.Second
					backoffReset = false
					responsesWithoutNew = 0
				}
			}

			// parse this into struct
			postsBody, err := lc.callLemmyAPI("GET", "post/list?sort=New", nil)
			if err != nil {
				lc.logger.Error(err.Error())
			}
			var postResponse PostsResponse
			err = json.Unmarshal(postsBody, &postResponse)
			if err != nil {
				lc.logger.Sugar().Infof("postsBody: %s", postsBody)
				lc.logger.Error(err.Error())
			}

			for _, postview := range postResponse.PostView {
				if !seenItems[postview.Post.ID] {
					postsChan <- postview.Post
					seenItems[postview.ID] = true
					backoffReset = true
				}
			}

			// Exponential backoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			// Wait for the posts channel to be closed or a timeout
			select {
			case <-closeChan:
				// The posts channel was closed as expected
				close(postsChan)
				lc.logger.Info("Comments channel closed.")
			case <-time.After(backoff):
			}

		}
	}()

	return postsChan
}
