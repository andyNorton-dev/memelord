package workers

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"api/src/middleware"
	"api/src/user"
	"api/src/core/loger"
	"strconv"
	"go.uber.org/zap"
)

type WorkerService struct {
	repo *WorkerRepository
	userService user.UserService
}

func NewWorkerService(repo *WorkerRepository, userService user.UserService) *WorkerService {
	return &WorkerService{
		repo: repo,
		userService: userService,
	}
}

func (s *WorkerService) GetWorkers(c echo.Context) error {
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	user, err := s.userService.GetUser(c.Request().Context(), telegramUser.ID)
	if err != nil {
		loger.Logger.Error("Ошибка при получении пользователя",
			zap.Error(err),
			zap.Int64("user_id", telegramUser.ID))
		return err
	}
	workers, err := s.repo.GetWorkers()
	if err != nil {
		loger.Logger.Error("Ошибка при получении работников",
			zap.Error(err))
		return err
	}
	result, err := s.GetWorkersOrArmy(workers, user)
	if err != nil {
		loger.Logger.Error("Ошибка при получении работников",
			zap.Error(err))
		return err
	}
	return c.JSON(200, result)
}

func (s *WorkerService) GetArmy(c echo.Context) error {
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	user, err := s.userService.GetUser(c.Request().Context(), telegramUser.ID)
	if err != nil {
		loger.Logger.Error("Ошибка при получении пользователя",
			zap.Error(err),
			zap.Int64("user_id", telegramUser.ID))
		return err
	}
	army, err := s.repo.GetArmy()
	if err != nil {
		loger.Logger.Error("Ошибка при получении работников",
			zap.Error(err))
		return err
	}
	result, err := s.GetWorkersOrArmy(army, user)
	if err != nil {
		loger.Logger.Error("Ошибка при получении работников",
			zap.Error(err))
		return err
	}
	return c.JSON(200, result)
}

func (s *WorkerService) GetWorkersOrArmy(workers []WorkerRepo, user *user.UserRepo) ([]UserWorkerResponse, error) {
	result := []UserWorkerResponse{}
	
	for _, worker := range workers {
		userWorker, err := s.repo.GetUserWorker(worker.ID, int(user.ID))
		var level int
		if err != nil {
			if err == sql.ErrNoRows {
				loger.Logger.Info("Работник не найден",
					zap.Int("worker_id", worker.ID),
					zap.Int("user_id", int(user.ID)))
				level = 0
			} else {
				loger.Logger.Error("Ошибка при получении работников пользователя",
					zap.Error(err))
				return nil, err
			}
		} else {
			workerUpgrade, err := s.repo.GetUpgradeById(userWorker.IdUpgrade)
			if err != nil {
				loger.Logger.Error("Ошибка при получении уровня работника",
					zap.Error(err))
				return nil, err
			}
			level = workerUpgrade.Level
		}
		workerUpgradeUserLevel, err := s.repo.GetWorkerUpgrades(worker.ID, level)
		if err != nil {
			if err == sql.ErrNoRows {
				loger.Logger.Info("Уровень работника не найден",
					zap.Int("worker_id", worker.ID),
					zap.Int("level", level))
				workerUpgradeUserLevel = WorkerUpgradeRepo{Profit: 0}
			} else {
				loger.Logger.Error("Ошибка при получении работников пользователя",
					zap.Error(err))
				return nil, err
			}
		}

		workerUpgradeNextLevel, err := s.repo.GetWorkerUpgrades(worker.ID, level+1)
		var accessToUpgrade bool
		var cost int
		if err != nil {
			if err == sql.ErrNoRows {
				accessToUpgrade = false
				cost = 0
			} else {
				loger.Logger.Error("Ошибка при получении работников пользователя",
					zap.Error(err))
				return nil, err
			}
		} else {
			if user.Balance >= int64(workerUpgradeNextLevel.Cost) {
				accessToUpgrade = true
			} else {
				accessToUpgrade = false
			}
			cost = workerUpgradeNextLevel.Cost
		}

		result = append(result, UserWorkerResponse{
			Id:              worker.ID,
			Name:            worker.Name,
			Description:     worker.Description,
			UrlImage:        worker.UrlImage,
			Level:           level,
			Profit:          workerUpgradeUserLevel.Profit,
			Cost:            cost,
			AccessToUpgrade: accessToUpgrade,
		})
	}
	return result, nil
}

func (s *WorkerService) BuyWorker(c echo.Context) error {
	workerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Неверный ID работника"})
	}

	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	user, err := s.userService.GetUser(c.Request().Context(), telegramUser.ID)
	if err != nil {
		loger.Logger.Error("Ошибка при получении пользователя",
			zap.Error(err),
			zap.Int64("user_id", telegramUser.ID))
		return err
	}

	userWorker, err := s.repo.GetUserWorker(workerID, int(user.ID))
	var level int
	var upgradeLevel WorkerUpgradeRepo
	nowProfit := 0
	if err != nil {
		if err == sql.ErrNoRows {
			level = 1
			upgradeLevel, err = s.repo.GetWorkerUpgrades(workerID, 1)
			if err != nil {
				loger.Logger.Error("Ошибка при получении работника пользователя",
					zap.Error(err),
					zap.Int("worker_id", workerID),
					zap.Int("user_id", int(user.ID)))
				return err
			}
		} else {
			loger.Logger.Error("Ошибка при получении работника пользователя",
				zap.Error(err),
				zap.Int("worker_id", workerID),
				zap.Int("user_id", int(user.ID)))
			return err
		}
	} else {
		level = userWorker.IdUpgrade + 1
		upgradeLevel, err = s.repo.GetWorkerUpgrades(workerID, level)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(400, map[string]string{"error": "Уровень работника максимальный"})
			} else {
				loger.Logger.Error("Ошибка при получении работника пользователя",
					zap.Error(err),
					zap.Int("worker_id", workerID),
					zap.Int("user_id", int(user.ID)))
				return err
			}
		}
		nowLevel, err := s.repo.GetWorkerUpgrades(workerID, userWorker.IdUpgrade)
		nowProfit = nowLevel.Profit
		if err != nil {
			loger.Logger.Error("Ошибка при получении работника пользователя",
				zap.Error(err),
				zap.Int("worker_id", workerID),
				zap.Int("user_id", int(user.ID)))
			return err
		}
	}

	if int64(upgradeLevel.Cost) > user.Balance {
		return c.JSON(400, map[string]string{"error": "Недостаточно средств"})
	}

	err = s.repo.UpdateUserBalance(int(user.ID), int(user.Balance - int64(upgradeLevel.Cost)))
	if err != nil {
		loger.Logger.Error("Ошибка при обновлении баланса пользователя",
			zap.Error(err),
			zap.Int("user_id", int(user.ID)))
		return err
	}

	if level == 1 {
		err = s.repo.CreateUserWorker(int(user.ID), workerID, upgradeLevel.ID)
		if err != nil {
			loger.Logger.Error("Ошибка при создании работника пользователя",
				zap.Error(err),
				zap.Int("user_id", int(user.ID)),
				zap.Int("worker_id", workerID))
			return err
		}
	} else {
		err = s.repo.UpdateUserWorker(upgradeLevel.ID, userWorker.ID)
		if err != nil {
			loger.Logger.Error("Ошибка при обновлении работника пользователя",
				zap.Error(err),
				zap.Int("user_id", int(user.ID)),
				zap.Int("worker_id", workerID))
			return err
		}
	}
	
	err = s.repo.UpdateUserWorkerProfit(int(user.ID), int(upgradeLevel.Profit - nowProfit), int(upgradeLevel.Cost))
	if err != nil {
		loger.Logger.Error("Ошибка при обновлении прибыли пользователя",
			zap.Error(err),
			zap.Int("user_id", int(user.ID)))
		return err
	}

	worker, err := s.repo.GetWorkerById(workerID)
	if err != nil {
		loger.Logger.Error("Ошибка при получении работника",
			zap.Error(err),
			zap.Int("worker_id", workerID))
		return err
	}
	if worker.Type == "worker" {
		return s.GetWorkers(c)
	} else {
		return s.GetArmy(c)
	}
}