package user

import (
	"context"
	"github.com/labstack/echo/v4"
	"api/src/middleware"
	"github.com/sirupsen/logrus"
	"time"
	"sync"
)

var log = logrus.New()

type UserService interface {
	GetUser(ctx context.Context, tg_id int64) (*UserRepo, error)
	CreateUser(ctx context.Context, tg_id int64, username string) (*UserRepo, error)
	SelectProfitForTap(ctx context.Context, tg_id int64) (int, error)
	UpdateBalanceForTap(ctx context.Context, tg_id int64, balance int) error
	UpdateEnergy(ctx context.Context, tg_id int64, energy int) error
	UpdateBalanceForProfitPerHour(ctx context.Context, tg_id int64, balance int) error
	AddProfitPerHour(ctx context.Context, u *UserRepo) (int64, error)
	EnergyRestoration(ctx context.Context, u *UserRepo) (int64, error)
	AddProfitPerHourAndEnergyRestoration(ctx context.Context, u *UserRepo) (profit int64, energy int64, err error)
	ReturnUser(u *UserRepo) (*UserResponse, error)
	GetUserHandler(c echo.Context) error
	CreateUserHandler(c echo.Context) error
	TapUserHandler(c echo.Context) error
}

type Service struct {
	repo *UserRepository
}

func NewService(repo *UserRepository) UserService {
	return &Service{repo: repo}
}

func (s *Service) GetUser(ctx context.Context, tg_id int64) (*UserRepo, error) {
	return s.repo.GetUser(ctx, tg_id)
}

func (s *Service) CreateUser(ctx context.Context, tg_id int64, username string) (*UserRepo, error) {
	return s.repo.CreateUser(ctx, tg_id, username)
}

func (s *Service) SelectProfitForTap(ctx context.Context, tg_id int64) (int, error) {
	return s.repo.SelectProfitForTap(ctx, tg_id)
}

func (s *Service) UpdateBalanceForTap(ctx context.Context, tg_id int64, balance int) error {
	return s.repo.UpdateBalanceForTap(ctx, tg_id, balance)
}

func (s *Service) UpdateEnergy(ctx context.Context, tg_id int64, energy int) error {
	return s.repo.UpdateEnergy(ctx, tg_id, energy)
}

func (s *Service) UpdateBalanceForProfitPerHour(ctx context.Context, tg_id int64, balance int) error {
	return s.repo.UpdateBalanceForProfitPerHour(ctx, tg_id, balance)
}

func (s *Service) AddProfitPerHour(ctx context.Context, u *UserRepo) (int64, error) {
	timeDifference := time.Since(u.LastProfitPerHour)
	minutes := int64(timeDifference.Minutes())
	intervals := minutes 
	if intervals == 0 {
		return u.Balance, nil
	}
	
	// Рассчитываем прибыль только за прошедшее время
	profit := ((int64(u.ProfitPerHour) / 60) * intervals)
	newBalance := u.Balance + profit
	
	// Обновляем баланс и время последнего начисления
	err := s.repo.UpdateBalanceForProfitPerHour(ctx, u.TgID, int(newBalance))
	if err != nil {
		log.WithError(err).Error("Ошибка при обновлении баланса")
		return 0, err
	}
	
	// Обновляем время последнего начисления в объекте
	u.LastProfitPerHour = time.Now()
	u.Balance = newBalance
	
	log.WithFields(logrus.Fields{
		"old_balance": u.Balance - profit,
		"profit": profit,
		"new_balance": newBalance,
		"minutes_passed": minutes,
		"profit_per_hour": u.ProfitPerHour,
	}).Info("Обновление баланса")
	
	return newBalance, nil
}

func (s *Service) EnergyRestoration(ctx context.Context, u *UserRepo) (int64, error) {
	timeDifference := time.Since(u.LastRestoration)
	minutes := int64(timeDifference.Minutes())
	intervals := minutes / 2
	needEnergy := int64(u.MaxEnergy - u.Energy)
	log.Info("Восстановление энергии: ", intervals, needEnergy, u.Energy)
	if intervals == 0 {
		return int64(u.Energy), nil
	}
	if intervals > needEnergy {
		intervals = needEnergy
	}
	err := s.repo.UpdateEnergy(ctx, u.TgID, int(intervals))
	if err != nil {
		log.WithError(err).Error("Ошибка при обновлении энергии")
		return 0, err
	}
	
	// Обновляем время последнего восстановления в объекте
	u.LastRestoration = time.Now()
	u.Energy += int(intervals)
	
	return int64(u.Energy), nil
}

func (s *Service) AddProfitPerHourAndEnergyRestoration(ctx context.Context, u *UserRepo) (profit int64, energy int64, err error) {
	type result struct {
		value int64
		err   error
	}

	profitChan := make(chan result)
	energyChan := make(chan result)

	go func() {
		profit, err := s.AddProfitPerHour(ctx, u)
		profitChan <- result{profit, err}
	}()

	go func() {
		energy, err := s.EnergyRestoration(ctx, u)
		energyChan <- result{energy, err}
	}()

	profitResult := <-profitChan
	if profitResult.err != nil {
		log.WithError(profitResult.err).Error("Ошибка при добавлении прибыли")
		return 0, 0, profitResult.err
	}

	energyResult := <-energyChan
	if energyResult.err != nil {
		log.WithError(energyResult.err).Error("Ошибка при восстановлении энергии")
		return 0, 0, energyResult.err
	}

	return profitResult.value, energyResult.value, nil
}

func (s *Service) ReturnUser(u *UserRepo) (*UserResponse, error) {
	user := &UserResponse{
		Username: u.Username,
		Balance: u.Balance,
		Level: u.Level,
		Energy: u.Energy,
		MaxEnergy: u.MaxEnergy,
		ProfitPerHour: u.ProfitPerHour,
		Head: u.Head,
		Body: u.Body,
		Legs: u.Legs,
		Foot: u.Foot,
		ProfitForTap: u.ProfitForTap,
	}
	return user, nil
}

func (s *Service) GetUserHandler(c echo.Context) error {
	ctx := c.Request().Context()
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	
	// Создаем канал для результатов
	type result struct {
		user *UserRepo
		profit int64
		energy int64
		err error
	}
	
	resultChan := make(chan result)
	
	go func() {
		// Получаем пользователя асинхронно
		user, err := s.GetUser(ctx, telegramUser.ID)
		if err != nil {
			resultChan <- result{err: err}
			return
		}
		
		// Обновляем профит и энергию асинхронно
		profit, energy, err := s.AddProfitPerHourAndEnergyRestoration(ctx, user)
		resultChan <- result{
			user: user,
			profit: profit,
			energy: energy,
			err: err,
		}
	}()
	
	// Ожидаем результат с таймаутом
	select {
	case <-ctx.Done():
		return c.JSON(500, "Request timeout")
	case r := <-resultChan:
		if r.err != nil {
			return c.JSON(500, r.err.Error())
		}
		userResponse, err := s.ReturnUser(r.user)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, userResponse)
	}
}

func (s *Service) CreateUserHandler(c echo.Context) error {
	log.Info("Начало выполнения CreateUser")
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	log.WithFields(logrus.Fields{
		"user_id": telegramUser.ID,
		"username": telegramUser.Username,
	}).Info("Получен telegramUser")
	
	user, err := s.repo.CreateUser(telegramUser.ID, telegramUser.Username)
	if err != nil {
		log.WithError(err).Error("Ошибка при создании пользователя")
		return c.JSON(500, err.Error())
	}
	log.WithFields(logrus.Fields{
		"user": user,
	}).Info("Пользователь успешно создан")
	return c.JSON(200, user)
}

func (s *Service) TapUserHandler(c echo.Context) error {
	ctx := c.Request().Context()
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	
	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	profitChan := make(chan int, 1)
	
	// Асинхронно получаем профит за тап
	wg.Add(1)
	go func() {
		defer wg.Done()
		profit, err := s.SelectProfitForTap(ctx, telegramUser.ID)
		if err != nil {
			errChan <- err
			return
		}
		profitChan <- profit
	}()
	
	// Ждем завершения горутин
	go func() {
		wg.Wait()
		close(errChan)
		close(profitChan)
	}()
	
	// Проверяем ошибки
	select {
	case err := <-errChan:
		if err != nil {
			return c.JSON(500, err.Error())
		}
	case <-ctx.Done():
		return c.JSON(500, "Request timeout")
	default:
	}
	
	// Получаем профит и обновляем баланс
	profit := <-profitChan
	if err := s.UpdateBalanceForTap(ctx, telegramUser.ID, profit); err != nil {
		return c.JSON(500, err.Error())
	}
	
	return c.JSON(200, map[string]interface{}{
		"profit": profit,
	})
}
