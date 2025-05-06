package user

import (
	"github.com/labstack/echo/v4"
	"api/src/middleware"
	"github.com/sirupsen/logrus"
	"time"
)

var log = logrus.New()

type UserService interface {
	GetUser(tg_id int64) (*UserRepo, error)
	CreateUser(tg_id int64, username string) (*UserRepo, error)
	SelectProfitForTap(tg_id int64) (int, error)
	UpdateBalanceForTap(tg_id int64, balance int) error
	UpdateEnergy(tg_id int64, energy int) error
	UpdateBalanceForProfitPerHour(tg_id int64, balance int) error
	AddProfitPerHour(u *UserRepo) (int64, error)
	EnergyRestoration(u *UserRepo) (int64, error)
	AddProfitPerHourAndEnergyRestoration(u *UserRepo) (profit int64, energy int64, err error)
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

func (s *Service) GetUser(tg_id int64) (*UserRepo, error) {
	return s.repo.GetUser(tg_id)
}

func (s *Service) CreateUser(tg_id int64, username string) (*UserRepo, error) {
	return s.repo.CreateUser(tg_id, username)
}

func (s *Service) SelectProfitForTap(tg_id int64) (int, error) {
	return s.repo.SelectProfitForTap(tg_id)
}

func (s *Service) UpdateBalanceForTap(tg_id int64, balance int) error {
	return s.repo.UpdateBalanceForTap(tg_id, balance)
}

func (s *Service) UpdateEnergy(tg_id int64, energy int) error {
	return s.repo.UpdateEnergy(tg_id, energy)
}

func (s *Service) UpdateBalanceForProfitPerHour(tg_id int64, balance int) error {
	return s.repo.UpdateBalanceForProfitPerHour(tg_id, balance)
}

func (s *Service) AddProfitPerHour(u *UserRepo) (int64, error) {
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
	err := s.repo.UpdateBalanceForProfitPerHour(u.TgID, int(newBalance))
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

func (s *Service) EnergyRestoration(u *UserRepo) (int64, error) {
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
	err := s.repo.UpdateEnergy(u.TgID, int(intervals))
	if err != nil {
		log.WithError(err).Error("Ошибка при обновлении энергии")
		return 0, err
	}
	
	// Обновляем время последнего восстановления в объекте
	u.LastRestoration = time.Now()
	u.Energy += int(intervals)
	
	return int64(u.Energy), nil
}

func (s *Service) AddProfitPerHourAndEnergyRestoration(u *UserRepo) (profit int64, energy int64, err error) {
	type result struct {
		value int64
		err   error
	}

	profitChan := make(chan result)
	energyChan := make(chan result)

	go func() {
		profit, err := s.AddProfitPerHour(u)
		profitChan <- result{profit, err}
	}()

	go func() {
		energy, err := s.EnergyRestoration(u)
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
	log.Info("Начало выполнения GetUser")
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	log.WithFields(logrus.Fields{
		"user_id": telegramUser.ID,
	}).Info("Получен telegramUser")
	
	type result struct {
		user  *UserRepo
		err   error
	}

	userChan := make(chan result)
	go func() {
		user, err := s.repo.GetUser(telegramUser.ID)
		userChan <- result{user, err}
	}()

	userResult := <-userChan
	if userResult.err != nil {
		log.WithError(userResult.err).Error("Ошибка при получении пользователя")
		return c.JSON(500, userResult.err.Error())
	}

	profit, energy, err := s.AddProfitPerHourAndEnergyRestoration(userResult.user)
	if err != nil {
		log.WithError(err).Error("Ошибка при восстановлении энергии")
		return c.JSON(500, err.Error())
	}

	userResponse, err := s.ReturnUser(userResult.user)
	if err != nil {
		log.WithError(err).Error("Ошибка при преобразовании пользователя")
		return c.JSON(500, err.Error())
	}
	
	userResponse.Energy = int(energy)
	userResponse.Balance = profit
	return c.JSON(200, userResponse)
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
	log.Info("Начало выполнения TapUser")
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	log.WithFields(logrus.Fields{
		"user_id": telegramUser.ID,
	}).Info("Получен telegramUser")
	
	profit, err := s.repo.SelectProfitForTap(telegramUser.ID)
	if err != nil {
		log.WithError(err).Error("Ошибка при получении прибыли")
		return c.JSON(500, err.Error())
	}
	log.WithFields(logrus.Fields{
		"profit": profit,
	}).Info("Получена прибыль")
	
	err = s.repo.UpdateBalanceForTap(telegramUser.ID, profit)
	if err != nil {
		log.WithError(err).Fatal("Критическая ошибка при обновлении баланса")
		return c.JSON(500, err.Error())
	}
	log.WithFields(logrus.Fields{
		"user_id": telegramUser.ID,
		"profit": profit,
	}).Info("Баланс успешно обновлен")
	return s.GetUserHandler(c)
}
