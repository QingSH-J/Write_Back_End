package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jinxinyu/go_backend/internal/api"
	"github.com/jinxinyu/go_backend/internal/models"
	"github.com/jinxinyu/go_backend/internal/storage"
	"github.com/jinxinyu/go_backend/internal/utils"
)

type Service struct {
	userRepo     storage.UserRepository
	timeout      time.Duration
	tokenmaker   utils.ToKenGenerator
	hashPassword utils.HashedPassword
	//To be soon added: emailservice
}

func NewService(userRepo storage.UserRepository, tokenmaker utils.ToKenGenerator, hashPassword utils.HashedPassword) *Service {
	return &Service{
		userRepo:     userRepo,
		tokenmaker:   tokenmaker,
		hashPassword: hashPassword,
		timeout:      time.Second * 60, // 设置默认超时时间为60秒
	}
}

func (s *Service) RegisterUser(ctx context.Context, req *api.RegisterRequest) (*models.User, error) {

	//check if the email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	} else if err != nil && !errors.Is(err, storage.ErrRecordNotFound) {
		log.Printf("查询用户时发生错误: %v", err)
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	hashedPassword, err := s.hashPassword.Hash(req.Password)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// add retry logic
	var createErr error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("创建用户重试尝试 %d/%d", attempt+1, maxRetries)
			time.Sleep(time.Second * time.Duration(attempt)) // 递增重试等待时间
		}

		startDbOp := time.Now()
		createErr = s.userRepo.Create(ctx, user)
		log.Printf("数据库创建操作耗时: %v", time.Since(startDbOp))

		if createErr == nil {
			break
		}

		log.Printf("创建用户尝试 %d 失败: %v", attempt+1, createErr)
	}

	if createErr != nil {
		return nil, fmt.Errorf("failed to create user: %w", createErr)
	}

	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, req *api.LoginRequest) (string, *models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	log.Printf("尝试登录用户: %s", req.Email)
	startTime := time.Now()

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
		return "", nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	log.Printf("成功查询到用户: %s, 用户ID: %s", user.Email, user.ID)

	log.Printf("开始验证密码")
	match, err := s.hashPassword.Compare(req.Password, user.Password)
	log.Printf("密码验证结果: 匹配=%v, 错误=%v", match, err)

	if err != nil {
		log.Printf("密码比较失败: %v", err)
		return "", nil, fmt.Errorf("failed to compare password: %w", err)
	}

	if !match {
		log.Printf("密码不匹配")
		return "", nil, fmt.Errorf("invalid credentials")
	}

	token, err := s.tokenmaker.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Printf("生成令牌失败: %v", err)
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	log.Printf("登录成功，用时: %v", time.Since(startTime))
	return token, user, nil
}
