package db

import "github.com/zemld/pr-manager/pr-manager/internal/domain"

type PullRequestStorage struct {
	Config
	Transactor
	selectQuery                  string
	createQuery                  string
	mergeQuery                   string
	reassignQuery                string
	userPullRequestsReviewsQuery string
}

func NewPullRequestStorage(config Config, transactor Transactor) *PullRequestStorage {
	return &PullRequestStorage{Config: config, Transactor: transactor}
}

func (s *PullRequestStorage) SetSelectQuery(selectQuery string) {
	s.selectQuery = selectQuery
}

func (s *PullRequestStorage) SetCreateQuery(createQuery string) {
	s.createQuery = createQuery
}

func (s *PullRequestStorage) SetMergeQuery(mergeQuery string) {
	s.mergeQuery = mergeQuery
}

func (s *PullRequestStorage) SetReassignQuery(reassignQuery string) {
	s.reassignQuery = reassignQuery
}

func (s *PullRequestStorage) SetUserPullRequestsReviewsQuery(query string) {
	s.userPullRequestsReviewsQuery = query
}

func (s *PullRequestStorage) Select(pullRequestID string) (domain.PullRequest, error) {
	rows, err := s.tx.Query(s.ctx, s.selectQuery, pullRequestID)
	if err != nil {
		return domain.PullRequest{}, err
	}
	defer rows.Close()

	var pullRequest domain.PullRequest
	if rows.Next() {
		err = rows.Scan(
			&pullRequest.ID,
			&pullRequest.Name,
			&pullRequest.AuthorID,
			&pullRequest.Status,
			&pullRequest.AssignedReviewers,
			&pullRequest.CreatedAt,
			&pullRequest.MergedAt,
		)
		if err != nil {
			return domain.PullRequest{}, err
		}
	}
	return pullRequest, nil
}

func (s *PullRequestStorage) Create(pullRequest domain.PullRequest) error {
	_, err := s.tx.Exec(s.ctx, s.createQuery,
		pullRequest.ID,
		pullRequest.Name,
		pullRequest.AuthorID,
		pullRequest.AssignedReviewers,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PullRequestStorage) Merge(pullRequest domain.PullRequest) error {
	_, err := s.tx.Exec(s.ctx, s.mergeQuery,
		pullRequest.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PullRequestStorage) Reassign(pullRequest domain.PullRequest) error {
	_, err := s.tx.Exec(s.ctx, s.reassignQuery,
		pullRequest.ID,
		pullRequest.AssignedReviewers,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PullRequestStorage) SelectUserPullRequestsReviews(userID string) ([]domain.PullRequest, error) {
	rows, err := s.tx.Query(s.ctx, s.userPullRequestsReviewsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pullRequests []domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		err = rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.AuthorID,
			&pr.Status,
			&pr.AssignedReviewers,
			&pr.CreatedAt,
			&pr.MergedAt,
		)
		if err != nil {
			return nil, err
		}
		pullRequests = append(pullRequests, pr)
	}
	return pullRequests, nil
}
