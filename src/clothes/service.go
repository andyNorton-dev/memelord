package clothes

import (
	"github.com/labstack/echo/v4"
	"api/src/core/loger"
	"api/src/middleware"
	"api/src/user"
	"net/http"
	"fmt"
	"go.uber.org/zap"
)

type ClothesService struct {
	repo *ClothesRepository
	userService user.UserService
}

func NewClothesService(repo *ClothesRepository, userService user.UserService) *ClothesService {
	return &ClothesService{repo: repo, userService: userService}
}

func (s *ClothesService) GetClothes(c echo.Context) error {
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	user, err := s.userService.GetUser(c.Request().Context(), telegramUser.ID)
	if err != nil {
		loger.Logger.Error("Ошибка при получении пользователя",
			zap.Error(err),
			zap.Int64("user_id", telegramUser.ID))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении пользователя"})
	}
	clothes, err := s.repo.GetClothes()
	if err != nil {
		loger.Logger.Error("Ошибка при получении одежды",
			zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении одежды"})
	}

	clothesUserResponse := []ClothesUserResponse{}
	for _, cloth := range clothes {
		exists, err := s.repo.ExistsClothesUser(user.ID, cloth.ID)
		if err != nil {
			loger.Logger.Error("Ошибка при проверке существования одежды пользователя",
				zap.Error(err),
				zap.Int("user_id", int(user.ID)),
				zap.Int("cloth_id", cloth.ID))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при проверке существования одежды пользователя"})
		}
		clothesUserResponse = append(clothesUserResponse, ClothesUserResponse{
			ID: cloth.ID,
			UrlImage: cloth.UrlImage,
			Price: cloth.Price,
			Type: cloth.Type,
			Rarity: cloth.Rarity,
			IsBought: exists,
		})
	}

	return c.JSON(http.StatusOK, clothesUserResponse)
}

func (s *ClothesService) GetClothe(c echo.Context) error {
	id := c.Param("id")
	clothe, err := s.repo.GetClothe(id)
	if err != nil {
		loger.Logger.Error("Ошибка при получении одежды",
			zap.Error(err),
			zap.String("cloth_id", id))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении одежды"})
	}
	
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	user, err := s.userService.GetUser(c.Request().Context(), telegramUser.ID)
	if err != nil {
		loger.Logger.Error("Ошибка при получении пользователя",
			zap.Error(err),
			zap.Int64("user_id", telegramUser.ID))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении пользователя"})
	}
	exists, err := s.repo.ExistsClothesUser(user.ID, clothe.ID)
	if err != nil {
		loger.Logger.Error("Ошибка при проверке существования одежды пользователя",
			zap.Error(err),
			zap.Int("user_id", int(user.ID)),
			zap.Int("cloth_id", clothe.ID))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при проверке существования одежды пользователя"})
	}

	var active bool
	var canBuy bool

	if clothe.Type == "head" {
		if user.Head != nil && *user.Head == clothe.UrlImage {
			active = true
		}
	} else if clothe.Type == "body" {
		if user.Body != nil && *user.Body == clothe.UrlImage {
			active = true
		}
	} else if clothe.Type == "legs" {
		if user.Legs != nil && *user.Legs == clothe.UrlImage {
			active = true
		}
	} else if clothe.Type == "foot" {
		if user.Foot != nil && *user.Foot == clothe.UrlImage {
			active = true
		}
	} else if clothe.Type == "hand" {
		if user.Hand != nil && *user.Hand == clothe.UrlImage {
			active = true
		}
	} else {
		active = false
	}

	if active {
		canBuy = false
	} else {
		if user.Balance >= int64(clothe.Price) {
			canBuy = true
		} else {
			canBuy = false
		}
	}

	clotheUserResponse := ClotheUserResponse{
		ID: clothe.ID,
		Name: clothe.Name,
		Description: clothe.Description,
		UrlImage: clothe.UrlImage,
		Price: clothe.Price,
		Type: clothe.Type,
		Rarity: clothe.Rarity,
		IsBought: exists,
		IsActive: active,
		CanBuy: canBuy,
		PerForTap: clothe.PerForTap,
		PlusEnergy: clothe.PlusEnergy,
	}

	return c.JSON(http.StatusOK, clotheUserResponse)
}

func (s *ClothesService) BuyClothe(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID не указан"})
	}

	clothe, err := s.repo.GetClothe(id)
	if err != nil {
		loger.Logger.Error("Ошибка при получении одежды",
			zap.Error(err),
			zap.String("cloth_id", id))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении одежды"})
	}
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	user, err := s.userService.GetUser(c.Request().Context(), telegramUser.ID)
	if err != nil {
		loger.Logger.Error("Ошибка при получении пользователя",
			zap.Error(err),
			zap.Int64("user_id", telegramUser.ID))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении пользователя"})
	}
    if user.Balance < int64(clothe.Price) {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Недостаточно средств"})
    }
	err = s.repo.AddClotheUser(user.ID, id)
	if err != nil {
		loger.Logger.Error("Ошибка при покупке одежды",
			zap.Error(err),
			zap.Int("user_id", int(user.ID)),
			zap.String("cloth_id", id))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при покупке одежды"})
	}
	err = s.repo.UpdateUserBalance(user.ID, user.Balance - int64(clothe.Price))
	if err != nil {
		loger.Logger.Error("Ошибка при обновлении баланса",
			zap.Error(err),
			zap.Int("user_id", int(user.ID)))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении баланса"})
	}
	return s.GetClothe(c)
}

func (s *ClothesService) EquipClothe(c echo.Context) error {
	id := c.Param("id")
	clothe, err := s.repo.GetClothe(id)
	if err != nil {
		loger.Logger.Error("Ошибка при получении одежды",
			zap.Error(err),
			zap.String("cloth_id", id))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении одежды"})
	}
	
	telegramUser := c.Get("telegram_user").(*middleware.TelegramUser)
	user, err := s.userService.GetUser(c.Request().Context(), telegramUser.ID)
	if err != nil {
		loger.Logger.Error("Ошибка при получении пользователя",
			zap.Error(err),
			zap.Int64("user_id", telegramUser.ID))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении пользователя"})
	}

	exists, err := s.repo.ExistsClothesUser(user.ID, clothe.ID)
	if err != nil {
		loger.Logger.Error("Ошибка при проверке существования одежды пользователя",
			zap.Error(err),
			zap.Int("user_id", int(user.ID)),
			zap.Int("cloth_id", clothe.ID))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при проверке существования одежды пользователя"})
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Одежда не куплена"})
	}	
	var old_clothe *ClotheRepo
	if clothe.Type == "head" {
		if user.Head != nil {
			old_clothe, err = s.repo.GetClotheByUrlImage(*user.Head)
			if err != nil {
				loger.Logger.Error("Ошибка при получении старой одежды",
					zap.Error(err),
					zap.String("url_image", *user.Head))
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении старой одежды"})
			}
		}else{
			old_clothe = nil
		}
		err = s.repo.UpdateUserHead(user.ID, clothe.UrlImage)
	} else if clothe.Type == "body" {
		if user.Body != nil {
			old_clothe, err = s.repo.GetClotheByUrlImage(*user.Body)
			if err != nil {
				loger.Logger.Error("Ошибка при получении старой одежды",
					zap.Error(err),
					zap.String("url_image", *user.Body))
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении старой одежды"})
			}
		}else{
			old_clothe = nil
		}
		err = s.repo.UpdateUserBody(user.ID, clothe.UrlImage)
	} else if clothe.Type == "legs" {
		if user.Legs != nil {
			old_clothe, err = s.repo.GetClotheByUrlImage(*user.Legs)
			if err != nil {
				loger.Logger.Error("Ошибка при получении старой одежды",
					zap.Error(err),
					zap.String("url_image", *user.Legs))
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении старой одежды"})
			}
		}else{
			old_clothe = nil
		}
		err = s.repo.UpdateUserLegs(user.ID, clothe.UrlImage)
	} else if clothe.Type == "foot" {
		if user.Foot != nil {
			old_clothe, err = s.repo.GetClotheByUrlImage(*user.Foot)
			if err != nil {
				loger.Logger.Error("Ошибка при получении старой одежды",
					zap.Error(err),
					zap.String("url_image", *user.Foot))
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении старой одежды"})
			}
		}else{
			old_clothe = nil
		}
		err = s.repo.UpdateUserFoot(user.ID, clothe.UrlImage)
	} else if clothe.Type == "hand" {
		if user.Hand != nil {
			old_clothe, err = s.repo.GetClotheByUrlImage(*user.Hand)
			if err != nil {
				loger.Logger.Error("Ошибка при получении старой одежды",
					zap.Error(err),
					zap.String("url_image", *user.Hand))
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении старой одежды"})
			}
		}else{
			old_clothe = nil
		}
		err = s.repo.UpdateUserHand(user.ID, clothe.UrlImage)
	}

	if err != nil {
		loger.Logger.Error("Ошибка при экипировке одежды",
			zap.Error(err),
			zap.Int("user_id", int(user.ID)),
			zap.Int("cloth_id", clothe.ID))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при экипировке одежды"})
	}

	var plas_energy int
	var per_for_tap int
	if old_clothe != nil {
		plas_energy = clothe.PlusEnergy - old_clothe.PlusEnergy
		loger.Logger.Debug("Расчет изменения энергии",
			zap.Int("plas_energy", plas_energy),
			zap.Int("old_energy", old_clothe.PlusEnergy),
			zap.Int("new_energy", clothe.PlusEnergy))
		per_for_tap = clothe.PerForTap - old_clothe.PerForTap
		loger.Logger.Debug("Расчет изменения per_for_tap",
			zap.Int("per_for_tap", per_for_tap),
			zap.Int("old_per_for_tap", old_clothe.PerForTap),
			zap.Int("new_per_for_tap", clothe.PerForTap))
	} else {
		plas_energy = clothe.PlusEnergy
		per_for_tap = clothe.PerForTap
	}

	err = s.repo.UpdateUserEnergy(user.ID, user.MaxEnergy + plas_energy)
	if err != nil {
		loger.Logger.Error("Ошибка при обновлении энергии",
			zap.Error(err),
			zap.Int("user_id", int(user.ID)))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении энергии"})
	}
	
	err = s.repo.UpdateUserPerForTap(user.ID, user.ProfitForTap + per_for_tap)
	if err != nil {
		loger.Logger.Error("Ошибка при обновлении per_for_tap",
			zap.Error(err),
			zap.Int("user_id", int(user.ID)))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении per_for_tap"})
	}
	
	return s.GetClothe(c)
}


