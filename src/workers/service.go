package workers

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"api/src/middleware"
	"api/src/user"
	"github.com/sirupsen/logrus"
	"strconv"
)

var log = logrus.New()

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
	user, err := s.userService.GetUser(telegramUser.ID)
	if err != nil {
		log.WithError(err).Error("Ошибка при получении пользователя", telegramUser.ID)
		return err
	}
	workers, err := s.repo.GetWorkers()
	if err != nil {
		log.WithError(err).Error("Ошибка при получении работников")
		return err
	}
	result, err := s.GetWorkersOrArmy(workers, user)
	if err != nil {
		log.WithError(err).Error("Ошибка при получении работников")
		return err
	}
	return c.JSON(200, result)
}

func (s *WorkerService) GetArmy(c echo.Context) error {
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	user, err := s.userService.GetUser(telegramUser.ID)
	if err != nil {
		log.WithError(err).Error("Ошибка при получении пользователя", telegramUser.ID)
		return err
	}
	army, err := s.repo.GetArmy()
	if err != nil {
		log.WithError(err).Error("Ошибка при получении работников")
		return err
	}
	result, err := s.GetWorkersOrArmy(army, user)
	if err != nil {
		log.WithError(err).Error("Ошибка при получении работников")
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
				log.WithError(err).Error("Работник не найден")
				level = 0
			} else {
				log.WithError(err).Error("Ошибка при получении работников пользователя")
				return nil, err
			}
		} else {
			workerUpgrade, err := s.repo.GetUpgradeById(userWorker.IdUpgrade)
			if err != nil {
				log.WithError(err).Error("Ошибка при получении уровня работника")
				return nil, err
			}
			level = workerUpgrade.Level
		}
		workerUpgradeUserLevel, err := s.repo.GetWorkerUpgrades(worker.ID, level)
		if err != nil {
			if err == sql.ErrNoRows {
				log.WithError(err).Error("Уровень работника не найден")
				workerUpgradeUserLevel = WorkerUpgradeRepo{Profit: 0}
			} else {
				log.WithError(err).Error("Ошибка при получении работников пользователя")
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
				log.WithError(err).Error("Ошибка при получении работников пользователя")
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
	user, err := s.userService.GetUser(telegramUser.ID)
	if err != nil {
		log.WithError(err).Error("Ошибка при получении пользователя", telegramUser.ID)
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
				log.WithError(err).Error("Ошибка при получении работника пользователя", workerID, user.ID)
				return err
			}
		} else {
			log.WithError(err).Error("Ошибка при получении работника пользователя", workerID, user.ID)
			return err
		}
	} else {
		level = userWorker.IdUpgrade + 1
		upgradeLevel, err = s.repo.GetWorkerUpgrades(workerID, level)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(400, map[string]string{"error": "Уровень работника максимальный"})
			} else {
				log.WithError(err).Error("Ошибка при получении работника пользователя", workerID, user.ID)
				return err
			}
		}
		nowLevel, err := s.repo.GetWorkerUpgrades(workerID, userWorker.IdUpgrade)
		nowProfit = nowLevel.Profit
		if err != nil {
			log.WithError(err).Error("Ошибка при получении работника пользователя", workerID, user.ID)
			return err
		}
	}

	if int64(upgradeLevel.Cost) > user.Balance {
		return c.JSON(400, map[string]string{"error": "Недостаточно средств"})
	}

	err = s.repo.UpdateUserBalance(int(user.ID), int(user.Balance - int64(upgradeLevel.Cost)))
	if err != nil {
		log.WithError(err).Error("Ошибка при обновлении баланса пользователя", user.ID)
		return err
	}

	if level == 1 {
		err = s.repo.CreateUserWorker(int(user.ID), workerID, upgradeLevel.ID)
		if err != nil {
			log.WithError(err).Error("Ошибка при создании работника пользователя", user.ID, workerID)
			return err
		}
		if err != nil {
			log.WithError(err).Error("Ошибка при обновлении прибыли пользователя", user.ID)
			return err
		}
	} else {
		err = s.repo.UpdateUserWorker(upgradeLevel.ID, userWorker.ID)
		if err != nil {
			log.WithError(err).Error("Ошибка при обновлении работника пользователя", user.ID, workerID)
			return err
		}
		
		if err != nil {
			log.WithError(err).Error("Ошибка при обновлении прибыли пользователя", user.ID)
			return err
		}
	}
	
	err = s.repo.UpdateUserWorkerProfit(int(user.ID), int(upgradeLevel.Profit - nowProfit), int(upgradeLevel.Cost))
	if err != nil {
		log.WithError(err).Error("Ошибка при обновлении прибыли пользователя", user.ID)
		return err
	}

	worker, err := s.repo.GetWorkerById(workerID)
	if err != nil {
		log.WithError(err).Error("Ошибка при получении работника", workerID)
		return err
	}
	if worker.Type == "worker" {
		return s.GetWorkers(c)
	} else {
		return s.GetArmy(c)
	}
}