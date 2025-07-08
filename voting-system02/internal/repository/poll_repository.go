package repository

import (
	"encoding/json"
	"os"
	"sync"

	"voting-system/internal/model"
	"voting-system/pkg/utils"
)

// 数据存储结构
type DataStore struct {
	Polls map[string]*model.Poll `json:"polls"`
	Users map[string]*model.User `json:"users"`
}

type PollRepository struct {
	store *DataStore
	mu    sync.RWMutex
	file  string
}

func NewPollRepository(filePath string) *PollRepository {
	repo := &PollRepository{
		store: &DataStore{
			Polls: make(map[string]*model.Poll),
			Users: make(map[string]*model.User),
		},
		file: filePath,
	}

	repo.loadFromFile()

	// 初始化管理员用户
	if _, exists := repo.store.Users["admin"]; !exists {
		repo.store.Users["admin"] = &model.User{
			ID:       utils.GenerateID(),
			Username: "admin",
			Password: "admin123", // 实际项目中应使用加密密码
		}
		repo.saveToFile()
	}

	return repo
}

func (r *PollRepository) CreatePoll(poll *model.Poll) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store.Polls[poll.ID] = poll
	return r.saveToFile()
}

func (r *PollRepository) GetPoll(id string) (*model.Poll, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	poll, exists := r.store.Polls[id]
	return poll, exists
}

func (r *PollRepository) GetAllPolls() []*model.Poll {
	r.mu.RLock()
	defer r.mu.RUnlock()

	polls := make([]*model.Poll, 0, len(r.store.Polls))
	for _, poll := range r.store.Polls {
		polls = append(polls, poll)
	}
	return polls
}

func (r *PollRepository) Vote(pollID, userID, optionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	poll, exists := r.store.Polls[pollID]
	if !exists {
		return ErrNotFound
	}

	// 检查用户是否已投票
	if _, voted := poll.Votes[userID]; voted {
		return ErrAlreadyVoted
	}

	// 找到选项索引
	optionIndex := -1
	for i, opt := range poll.Options {
		if opt.ID == optionID {
			optionIndex = i
			break
		}
	}

	if optionIndex == -1 {
		return ErrInvalidOption
	}

	// 记录投票
	poll.Votes[userID] = optionIndex
	return r.saveToFile()
}

func (r *PollRepository) Authenticate(username, password string) (*model.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.store.Users {
		if user.Username == username && user.Password == password {
			return user, true
		}
	}
	return nil, false
}

// 辅助方法
func (r *PollRepository) loadFromFile() {
	file, err := os.Open(r.file)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(r.store); err != nil {
		panic(err)
	}
}

func (r *PollRepository) saveToFile() error {
	file, err := os.Create(r.file)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r.store)
}

// 错误定义
var (
	ErrNotFound      = NewError("poll not found")
	ErrAlreadyVoted  = NewError("user already voted")
	ErrInvalidOption = NewError("invalid option")
)

func NewError(msg string) error {
	return &RepositoryError{msg}
}

type RepositoryError struct {
	msg string
}

func (e *RepositoryError) Error() string {
	return e.msg
}
