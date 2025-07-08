package service

import (
	"time"

	"voting-system/internal/model"
	"voting-system/internal/repository"
	"voting-system/pkg/utils"
)

type PollService struct {
	repo *repository.PollRepository
}

func NewPollService(repo *repository.PollRepository) *PollService {
	return &PollService{repo: repo}
}

func (s *PollService) CreatePoll(req *model.CreatePollRequest) (*model.Poll, error) {
	// 创建投票对象
	poll := &model.Poll{
		ID:          utils.GenerateID(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now(),
		EndAt:       req.EndAt,
		Votes:       make(map[string]int),
	}

	// 添加选项
	for _, text := range req.Options {
		poll.Options = append(poll.Options, model.PollOption{
			ID:   utils.GenerateID(),
			Text: text,
		})
	}

	// 保存到存储
	if err := s.repo.CreatePoll(poll); err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *PollService) GetAllPolls() []*model.Poll {
	return s.repo.GetAllPolls()
}

func (s *PollService) GetPoll(id string) (*model.Poll, error) {
	poll, exists := s.repo.GetPoll(id)
	if !exists {
		return nil, repository.ErrNotFound
	}
	return poll, nil
}

func (s *PollService) Vote(pollID, userID, optionID string) error {
	// 检查投票是否过期
	poll, exists := s.repo.GetPoll(pollID)
	if !exists {
		return repository.ErrNotFound
	}

	if time.Now().After(poll.EndAt) {
		return NewServiceError("投票已结束")
	}

	// 执行投票
	return s.repo.Vote(pollID, userID, optionID)
}

func (s *PollService) Authenticate(username, password string) (*model.User, bool) {
	return s.repo.Authenticate(username, password)
}

// 服务层错误
type ServiceError struct {
	msg string
}

func (e *ServiceError) Error() string {
	return e.msg
}

func NewServiceError(msg string) error {
	return &ServiceError{msg}
}
