package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v85/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type PRHandler struct {
	githubapp.ClientCreator
	cfg *Config
}

func (h *PRHandler) Handles() []string {
	return []string{"pull_request"}
}

func (h *PRHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.PullRequestEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	action := event.GetAction()
	if action != "opened" && action != "synchronize" && action != "reopened" {
		return nil
	}

	repo := event.GetRepo()
	pr := event.GetPullRequest()
	owner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	prNum := pr.GetNumber()

	logger := log.With().
		Str("repo", fmt.Sprintf("%s/%s", owner, repoName)).
		Int("pr", prNum).
		Logger()

	client, err := h.NewInstallationClient(event.GetInstallation().GetID())
	if err != nil {
		return err
	}

	aiAuthor, err := h.findaiCoAuthor(ctx, client, owner, repoName, prNum, logger)
	if err != nil {
		return err
	}
	if aiAuthor == "" {
		return nil
	}

	logger.Info().Str("co_author", aiAuthor).Msg("ai co-author detected, declining pr")
	return h.declinePR(ctx, client, owner, repoName, prNum, aiAuthor)
}

func (h *PRHandler) findaiCoAuthor(
	ctx context.Context,
	client *github.Client,
	owner, repo string,
	prNum int,
	logger zerolog.Logger,
) (string, error) {
	opts := &github.ListOptions{PerPage: 100}
	for {
		commits, resp, err := client.PullRequests.ListCommits(ctx, owner, repo, prNum, opts)
		if err != nil {
			return "", err
		}

		for _, c := range commits {
			msg := c.GetCommit().GetMessage()
			for _, line := range strings.Split(msg, "\n") {
				line = strings.TrimSpace(line)
				if !strings.HasPrefix(strings.ToLower(line), "co-authored-by:") {
					continue
				}
				coAuthor := strings.TrimSpace(line[len("co-authored-by:"):])
				if matchesAI(coAuthor, h.cfg.ExtraPatterns) {
					logger.Debug().
						Str("commit", c.GetSHA()).
						Str("co_author", coAuthor).
						Msg("matched ai co-author")
					return coAuthor, nil
				}
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return "", nil
}

func (h *PRHandler) declinePR(
	ctx context.Context,
	client *github.Client,
	owner, repo string,
	prNum int,
	SlopGeneratorAuthor string,
) error {
	body := h.cfg.DeclineMessage + fmt.Sprintf("\n\n**detected:** `%s`", SlopGeneratorAuthor)
	_, _, err := client.Issues.CreateComment(ctx, owner, repo, prNum, &github.IssueComment{
		Body: &body,
	})
	if err != nil {
		return fmt.Errorf("post comment: %w", err)
	}

	state := "closed"
	_, _, err = client.PullRequests.Edit(ctx, owner, repo, prNum, &github.PullRequest{
		State: &state,
	})
	if err != nil {
		return fmt.Errorf("close pr: %w", err)
	}
	return nil
}
