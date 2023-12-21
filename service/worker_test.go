package lemmylinks_service

import (
	"errors"
	"testing"
	"time"

	"github.com/adadkins/glaw"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWorker(t *testing.T) {
	t.Run("Worker messages a user who has a linked comment", func(t *testing.T) {
		// given
		mockGlawClient := NewMockGlawClient()

		// set the comments we want to loop over
		inputComments := []glaw.Comment{
			{
				Content:   "Hey Check out this comment https://someURL.com/comment/1111",
				CreatorID: 1111,
				ApID:      "https://https://someURL.com/comment/1",
			},
			{
				Content:   "Hey Check out this comment https://someURL.com/comment/1111",
				CreatorID: 2222,
				ApID:      "https://https://someURL.com/comment/2",
			},
		}

		mockGlawClient.SetComments(inputComments)
		a, _ := NewApp(&mockGlawClient, zap.NewExample(), []int{})
		// when
		// do some work async
		go a.Work(&mockGlawClient, "https://someURL.com/")

		time.Sleep(2 * time.Second)

		// then
		// get the pms
		sentPMs := mockGlawClient.GetSentPMs()

		expected := []MockPM{{id: 1111, message: "One of your comments was linked by this comment: https://https://someURL.com/comment/2"}}
		assert.Equal(t, expected, sentPMs)
	})
	t.Run("Worker messages a user who has a linked comment using Markdown", func(t *testing.T) {
		// given
		mockGlawClient := NewMockGlawClient()

		// set the comments we want to loop over
		inputComments := []glaw.Comment{
			{
				Content:   "Hey Check out this comment https://someURL.com/comment/1111",
				CreatorID: 1111,
				ApID:      "https://https://someURL.com/comment/1",
			},
			{
				Content:   "[Hey Check out this comment](https://someURL.com/comment/1111)",
				CreatorID: 2222,
				ApID:      "https://https://someURL.com/comment/2",
			},
		}

		mockGlawClient.SetComments(inputComments)
		a, _ := NewApp(&mockGlawClient, zap.NewExample(), []int{})
		// when
		// do some work async
		go a.Work(&mockGlawClient, "https://someURL.com/")

		time.Sleep(2 * time.Second)

		// then
		// get the pms
		sentPMs := mockGlawClient.GetSentPMs()

		expected := []MockPM{{id: 1111, message: "One of your comments was linked by this comment: https://https://someURL.com/comment/2"}}
		assert.Equal(t, expected, sentPMs)
	})
	t.Run("Worker messages a user who has a linked a post", func(t *testing.T) {
		// given
		mockGlawClient := NewMockGlawClient()

		// set the comments we want to loop over
		inputComments := []glaw.Comment{
			{
				Content:   "Hey Check out this post https://someURL.com/post/2222",
				CreatorID: 1111,
				ApID:      "https://https://someURL.com/comment/3",
			},
			{
				Content:   "Hey Check out this post https://someURL.com/post/2222",
				CreatorID: 2222,
				ApID:      "https://https://someURL.com/comment/4",
			},
		}

		samplePosts := []glaw.Post{{ApID: "https://https://someURL.com/comment/3", CreatorID: 1111}}

		mockGlawClient.SetComments(inputComments)
		mockGlawClient.SetPosts(samplePosts)
		a, _ := NewApp(&mockGlawClient, zap.NewExample(), []int{})
		// when
		// do some work async
		go a.Work(&mockGlawClient, "https://someURL.com/")

		time.Sleep(2 * time.Second)

		// then
		// get the pms
		sentPMs := mockGlawClient.GetSentPMs()

		expected := []MockPM{{id: 1111, message: "One of your posts was linked by this comment: https://https://someURL.com/comment/4"}}
		assert.Equal(t, expected, sentPMs)
	})
	t.Run("Worker doesnt messages a user who has comment linked by a banlisted user", func(t *testing.T) {
		// given
		mockGlawClient := NewMockGlawClient()

		// set the comments we want to loop over
		inputComments := []glaw.Comment{
			{
				Content:   "Hey Check out this comment https://someURL.com/comment/1111",
				CreatorID: 1111,
				ApID:      "https://https://someURL.com/comment/1",
			},
			{
				Content:   "Hey Check out this comment https://someURL.com/comment/1111",
				CreatorID: 2222,
				ApID:      "https://https://someURL.com/comment/2",
			},
		}

		mockGlawClient.SetComments(inputComments)
		a, _ := NewApp(&mockGlawClient, zap.NewExample(), []int{1111, 2222})
		// when
		// do some work async
		go a.Work(&mockGlawClient, "https://someURL.com/")

		time.Sleep(2 * time.Second)

		// then
		// get the pms
		sentPMs := mockGlawClient.GetSentPMs()

		expected := []MockPM(nil)
		assert.Equal(t, expected, sentPMs)
	})

	t.Run("Worker doesnt message a user who has a post linked by a banlisted user", func(t *testing.T) {
		// given
		mockGlawClient := NewMockGlawClient()

		// set the comments we want to loop over
		inputComments := []glaw.Comment{
			{
				Content:   "Hey Check out this post https://someURL.com/post/2222",
				CreatorID: 1111,
				ApID:      "https://https://someURL.com/comment/3",
			},
			{
				Content:   "Hey Check out this post https://someURL.com/post/2222",
				CreatorID: 2222,
				ApID:      "https://https://someURL.com/comment/4",
			},
		}

		samplePosts := []glaw.Post{{ApID: "https://https://someURL.com/comment/3", CreatorID: 1111}}

		mockGlawClient.SetComments(inputComments)
		mockGlawClient.SetPosts(samplePosts)
		a, _ := NewApp(&mockGlawClient, zap.NewExample(), []int{1111, 2222})
		// when
		// do some work async
		go a.Work(&mockGlawClient, "https://someURL.com/")

		time.Sleep(2 * time.Second)

		// then
		// get the pms
		sentPMs := mockGlawClient.GetSentPMs()

		expected := []MockPM(nil)
		assert.Equal(t, expected, sentPMs)
	})
	// t.Run("Worker handles a shutdown of the comment channel", func(t *testing.T) {
	// 	// given
	// 	mockGlawClient := NewMockGlawClient()

	// 	// set the comments we want to loop over
	// 	inputComments := []glaw.Comment{
	// 		{
	// 			Content:   "Hey Check out this comment https://someURL.com/comment/1111",
	// 			CreatorID: 1111,
	// 			ApID:      "1",
	// 		},
	// 		{
	// 			Content:   "Hey Check out this comment https://someURL.com/comment/1111",
	// 			CreatorID: 2222,
	// 			ApID:      "2",
	// 		},
	// 	}

	// 	mockGlawClient.SetComments(inputComments)
	// 	a, _ := NewApp(&mockGlawClient)
	// 	// when
	// 	// do some work async
	// 	go a.Work(&mockGlawClient, "https://someURL.com/")

	// 	time.Sleep(2 * time.Second)

	// 	// send the shudown
	// 	a.done <- struct{}{}

	// 	// then
	// 	// get the pms
	// 	sentPMs := mockGlawClient.GetSentPMs()

	// 	expected := []MockPM{{id: 1111, message: "Your comment was linked by in this comment: https://https://someURL.com/comment/2"}}
	// 	assert.Equal(t, expected, sentPMs)
	// })
	t.Run("Worker handles a baseURL that isnt a post/comment/or is incorrectly formated", func(t *testing.T) {
		// given
		mockGlawClient := NewMockGlawClient()

		// set the comments we want to loop over
		inputComments := []glaw.Comment{
			{
				Content:   "Hey Check out this profile https://someURL.com/profile/1",
				CreatorID: 1111,
				ApID:      "3",
			},
			{
				Content:   "Hey Check out your settings https://someURL.com/settings",
				CreatorID: 2222,
				ApID:      "4",
			},
			{
				Content:   "Hey Check out this bad link https://someURL.com/comment/abcd",
				CreatorID: 1111,
				ApID:      "3",
			},
			{
				Content:   "Hey Check out your bad link https://someURL.com/post/abcd",
				CreatorID: 2222,
				ApID:      "4",
			},
			{
				Content:   "Hey Check out bad link https://someURL.com/comment//pathdoesntexist",
				CreatorID: 1111,
				ApID:      "3",
			},
			{
				Content:   "Hey Check out bad link https://someURL.com/post//pathdoesntexist",
				CreatorID: 2222,
				ApID:      "4",
			},
		}

		samplePosts := []glaw.Post{{ApID: "3", CreatorID: 1111}}

		mockGlawClient.SetComments(inputComments)
		mockGlawClient.SetPosts(samplePosts)
		a, _ := NewApp(&mockGlawClient, zap.NewExample(), []int{})
		// when
		// do some work async
		go a.Work(&mockGlawClient, "https://someURL.com/")

		time.Sleep(2 * time.Second)

		// then
		// get the pms
		sentPMs := mockGlawClient.GetSentPMs()

		expected := []MockPM(nil)
		assert.Equal(t, expected, sentPMs)
	})

	t.Run("Worker handles a linked comment that doesnt exist/errors", func(t *testing.T) {
		// given
		mockGlawClient := NewMockGlawClient()

		// set the comments we want to loop over
		inputComments := []glaw.Comment{
			{
				Content:   "Hey Check out this comment https://someURL.com/comment/1111",
				CreatorID: 1111,
				ApID:      "1",
			},
			{
				Content:   "Hey Check out this comment https://someURL.com/comment/1111",
				CreatorID: 2222,
				ApID:      "2",
			},
		}

		mockGlawClient.SetComments(inputComments)
		mockGlawClient.err = errors.New("Get Comment Error")
		a, _ := NewApp(&mockGlawClient, zap.NewExample(), []int{})

		// when
		// do some work async
		go a.Work(&mockGlawClient, "https://someURL.com/")

		time.Sleep(2 * time.Second)

		// then
		// get the pms
		sentPMs := mockGlawClient.GetSentPMs()

		expected := []MockPM(nil)
		assert.Equal(t, expected, sentPMs)
	})

	t.Run("Worker handles a linked Post that doesnt exist/errors", func(t *testing.T) {
		// given
		mockGlawClient := NewMockGlawClient()

		// set the comments we want to loop over
		inputComments := []glaw.Comment{
			{
				Content:   "Hey Check out this comment https://someURL.com/post/1111",
				CreatorID: 1111,
				ApID:      "1",
			},
			{
				Content:   "Hey Check out this comment https://someURL.com/post/1111",
				CreatorID: 2222,
				ApID:      "2",
			},
		}

		mockGlawClient.SetComments(inputComments)
		mockGlawClient.err = errors.New("Get Post Error")
		a, _ := NewApp(&mockGlawClient, zap.NewExample(), []int{})

		// when
		go a.Work(&mockGlawClient, "https://someURL.com/")

		// then
		time.Sleep(2 * time.Second)
		sentPMs := mockGlawClient.GetSentPMs()
		expected := []MockPM(nil)

		assert.Equal(t, expected, sentPMs)
	})

	t.Run("Worker handles sending a private message that errors", func(t *testing.T) {
		// given
		mockGlawClient := NewMockGlawClient()

		// set the comments we want to loop over
		inputComments := []glaw.Comment{
			{
				Content:   "Hey Check out this comment https://someURL.com/comment/1111",
				CreatorID: 1111,
				ApID:      "1",
			},
			{
				Content:   "Hey Check out this comment https://someURL.com/comment/1111",
				CreatorID: 2222,
				ApID:      "2",
			},
			{
				Content:   "Hey Check out this post https://someURL.com/post/1111",
				CreatorID: 1111,
				ApID:      "1",
			},
			{
				Content:   "Hey Check out this post https://someURL.com/post/1111",
				CreatorID: 2222,
				ApID:      "2",
			},
		}

		mockGlawClient.SetComments(inputComments)
		mockGlawClient.err = errors.New("Private Message Error")
		samplePosts := []glaw.Post{{ApID: "3", CreatorID: 1111}}
		mockGlawClient.SetPosts(samplePosts)
		a, _ := NewApp(&mockGlawClient, zap.NewExample(), []int{})

		// when
		go a.Work(&mockGlawClient, "https://someURL.com/")

		// then
		time.Sleep(2 * time.Second)
		sentPMs := mockGlawClient.GetSentPMs()
		expected := []MockPM(nil)

		assert.Equal(t, expected, sentPMs)
	})
}

func NewMockGlawClient() MockGlawClient {
	return MockGlawClient{}
}

type MockPM struct {
	id      int
	message string
}

type MockGlawClient struct {
	comments []glaw.Comment
	posts    []glaw.Post
	sentPMs  []MockPM
	err      error
}

// create mock functions
func (a *MockGlawClient) StreamNewPosts(pauseAfter int, closeChan chan struct{}) chan glaw.Post {
	return nil
}
func (a *MockGlawClient) SendPrivateMessage(messageBody string, targetUserID int) error {
	if a.err != nil && a.err.Error() == "Private Message Error" {
		return a.err
	}
	a.sentPMs = append(a.sentPMs, MockPM{targetUserID, messageBody})

	return nil
}
func (a *MockGlawClient) StreamNewComments(pauseAfter int, closeChan chan struct{}) chan glaw.Comment {
	commentsChan := make(chan glaw.Comment, 1000)

	// add some
	for _, v := range a.comments {
		commentsChan <- v
	}

	return commentsChan
}
func (a *MockGlawClient) GetPost(id int) (glaw.Post, error) {
	if a.err != nil && a.err.Error() == "Get Post Error" {
		return glaw.Post{}, a.err
	}
	return a.posts[0], nil
}
func (a *MockGlawClient) GetComment(id int) (glaw.Comment, error) {
	if a.err != nil && a.err.Error() == "Get Comment Error" {
		return glaw.Comment{}, a.err
	}
	return a.comments[0], nil
}

// mock helpers
func (a *MockGlawClient) SetError(err error) {
	a.err = err
}
func (a *MockGlawClient) SetComments(comments []glaw.Comment) {
	a.comments = comments
}
func (a *MockGlawClient) SetPosts(posts []glaw.Post) {
	a.posts = posts
}
func (a *MockGlawClient) GetSentPMs() []MockPM {
	return a.sentPMs
}
