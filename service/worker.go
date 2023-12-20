package lemmylinks_service

import (
	"fmt"
	glaw "lemmy_links_bot/lib"
	"regexp"
	"strconv"
	"strings"
)

func (a *App) Work(client glaw.Client, baseURL string) error {
	a.logger.Info("Starting streaming comments...")
	postComments := client.StreamNewComments(5, a.done)
	for comment := range postComments {
		urls := extractURLs(comment.Content)

		for _, v := range urls {
			if strings.Contains(v, baseURL) {
				if strings.Contains(v, "comment") {
					id, err := a.extractPath(v)
					if err != nil {
						a.logger.Sugar().Infoln("Link found: %s", v)
						a.logger.Error(err.Error())
						continue
					}
					c, err := client.GetComment(id)
					if err != nil {
						a.logger.Error(err.Error())
						continue
					}
					if c.CreatorID == comment.CreatorID {
						a.logger.Info("Author of comment and linked comment was the same, not messaging")
						continue
					}
					a.logger.Sugar().Infof("Found a comment to message: %s", comment.ApID)
					// message the comments author
					err = client.SendPrivateMessage(fmt.Sprintf("One of your comments was linked by this comment: %s", comment.ApID), c.CreatorID)
					if err != nil {
						a.logger.Error(err.Error())
						continue
					}
				}
				if strings.Contains(v, "post") {
					id, err := a.extractPath(v)
					if err != nil {
						a.logger.Sugar().Infoln("Link found: %s", v)
						a.logger.Error(err.Error())
						continue
					}
					post, err := client.GetPost(id)
					if err != nil {
						a.logger.Error(err.Error())
						continue
					}
					if post.CreatorID == comment.CreatorID {
						a.logger.Info("Author of comment and post was the same, not messaging")
						continue
					}
					// message the comments author
					err = client.SendPrivateMessage(fmt.Sprintf("One of your posts was linked by this comment: %s", comment.ApID), post.CreatorID)
					if err != nil {
						a.logger.Error(err.Error())
						continue
					}
					a.logger.Sugar().Infof("Found a comment to message: %s", comment.ApID)
				}
			}
		}
	}
	a.logger.Info("Post channel closed.")
	return nil
}

func extractURLs(input string) []string {
	// Define a regular expression pattern for URLs
	urlPattern := regexp.MustCompile(`https?://[^\s]+`)

	// Find all matches in the input string
	matches := urlPattern.FindAllString(input, -1)

	return matches
}

func (a *App) extractPath(input string) (int, error) {
	// Define a regular expression pattern to match the integer part
	pattern := regexp.MustCompile(`\d+`)

	// Use FindStringSubmatch to find the pattern in the input
	matches := pattern.FindStringSubmatch(input)

	// Check if a match is found
	if len(matches) > 0 {
		// Extract the matched integer
		extractedInteger := matches[0]

		// Convert the string to an integer
		intValue, err := strconv.Atoi(extractedInteger)
		if err != nil {
			return 0, err
		}

		return intValue, nil
	}

	return 0, fmt.Errorf("no match found")
}
