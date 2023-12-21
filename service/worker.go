package lemmylinks_service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/adadkins/glaw"
)

func (a *App) Work() error {
	a.logger.Info("Starting streaming comments...")
	postComments := a.lemmyClient.StreamNewComments(5, a.done)
	for comment := range postComments {
		urls := extractURLs(comment.Content)
	commentLoop:
		for _, url := range urls {
			if strings.Contains(url, a.baseURL) {
				if strings.Contains(url, "comment") {
					err := a.handleLinkingToComment(comment, url)
					if err != nil {
						continue commentLoop
					}
				}
				if strings.Contains(url, "post") {
					err := a.handleLinkingToPost(comment, url)
					if err != nil {
						continue commentLoop
					}
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

func extractPath(input string) (int, error) {
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

func (a *App) handleLinkingToComment(comment glaw.Comment, link string) error {
	id, err := extractPath(link)
	if err != nil {
		a.logger.Sugar().Infoln("Link found: %s", link)
		a.logger.Error(err.Error())
		return err
	}
	c, err := a.lemmyClient.GetComment(id)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}
	if c.CreatorID == comment.CreatorID {
		msg := "Author of comment and linked comment was the same, not messaging"
		a.logger.Info(msg)
		return fmt.Errorf(msg)
	}
	// check if linker is a ban listed account
	for _, v := range a.banListedAccounts {
		if c.CreatorID == v || comment.CreatorID == v {
			msg := fmt.Sprintf("BANLISTED! Comment: %s, ApiID: %s, postCreatorID: %v, commentCreatorID: %v, Blacklisted: %v", comment.Content, comment.ApID, c.CreatorID, comment.CreatorID, v)
			a.logger.Sugar().Infof(msg)
			return fmt.Errorf(msg)
		}
	}
	a.logger.Sugar().Infof("Found a comment to message: %s", comment.ApID)
	// message the comments author
	err = a.lemmyClient.SendPrivateMessage(fmt.Sprintf("One of your comments was linked by this comment: %s", comment.ApID), c.CreatorID)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}
	return nil
}

func (a *App) handleLinkingToPost(comment glaw.Comment, link string) error {
	id, err := extractPath(link)
	if err != nil {
		a.logger.Sugar().Infoln("Link found: %s", link)
		a.logger.Error(err.Error())
		return err
	}
	post, err := a.lemmyClient.GetPost(id)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}
	if post.CreatorID == comment.CreatorID {
		msg := "Author of comment and post was the same, not messaging"
		a.logger.Info(msg)
		return fmt.Errorf(msg)
	}
	for _, v := range a.banListedAccounts {
		if post.CreatorID == v || comment.CreatorID == v {
			msg := fmt.Sprintf("BANLISTED! Comment: %s, ApiID: %s, postCreatorID: %v, commentCreatorID: %v, Banlisted: %v", comment.Content, comment.ApID, post.CreatorID, comment.CreatorID, v)
			a.logger.Sugar().Infof(msg)
			return fmt.Errorf(msg)
		}
	}
	// message the comments author
	err = a.lemmyClient.SendPrivateMessage(fmt.Sprintf("One of your posts was linked by this comment: %s", comment.ApID), post.CreatorID)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}
	a.logger.Sugar().Infof("Found a comment to message: %s", comment.ApID)
	return nil
}
