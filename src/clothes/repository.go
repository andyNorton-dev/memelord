package clothes

import (
	"database/sql"
	"strconv"
)

type ClothesRepository struct {
	db *sql.DB
}

func NewClothesRepository(db *sql.DB) *ClothesRepository {
	return &ClothesRepository{db: db}
}

func (r *ClothesRepository) GetClothes() ([]ClothesRepo, error) {
	rows, err := r.db.Query("SELECT id, url_image, price, type, rarity FROM clothes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	clothes := []ClothesRepo{}
	for rows.Next() {
		var cloth ClothesRepo
		err = rows.Scan(&cloth.ID, &cloth.UrlImage, &cloth.Price, &cloth.Type, &cloth.Rarity)
		if err != nil {
			return nil, err
		}
		clothes = append(clothes, cloth)
	}
	return clothes, nil
}

func (r *ClothesRepository) GetClothe(id string) (*ClotheRepo, error) {
	clothID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var cloth ClotheRepo
	err = r.db.QueryRow("SELECT id, name, description, url_image, price, type, rarity, per_for_tap, plus_energy FROM clothes WHERE id = $1", clothID).
		Scan(&cloth.ID, &cloth.Name, &cloth.Description, &cloth.UrlImage, &cloth.Price, &cloth.Type, &cloth.Rarity, &cloth.PerForTap, &cloth.PlusEnergy)
	if err != nil {
		return nil, err
	}
	return &cloth, nil
}

func (r *ClothesRepository) GetClothesUser(userID int) ([]ClothesUserRepo, error) {
	rows, err := r.db.Query("SELECT id, clothes_id, user_id FROM clothes_user WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	clothesUser := []ClothesUserRepo{}
	for rows.Next() {
		var clothUser ClothesUserRepo
		err = rows.Scan(&clothUser.ID, &clothUser.ClothesID, &clothUser.UserID)
		if err != nil {
			return nil, err
		}
		clothesUser = append(clothesUser, clothUser)
	}
	return clothesUser, nil
}

func (r *ClothesRepository) ExistsClothesUser(userID int, clothesID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM clothes_user WHERE user_id = $1 AND clothes_id = $2)", userID, clothesID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *ClothesRepository) AddClotheUser(userID int, clothesID string) error {
	clothID, err := strconv.Atoi(clothesID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("INSERT INTO clothes_user (user_id, clothes_id) VALUES ($1, $2)", userID, clothID)
	return err
}

func (r *ClothesRepository) UpdateUserBalance(userID int, balance int64) error {
	_, err := r.db.Exec("UPDATE users SET balance = $1 WHERE id = $2", balance, userID)
	return err
}

func (r *ClothesRepository) EquipClothe(userID int, clothesID string) error {
	clothID, err := strconv.Atoi(clothesID)
	if err != nil {
		return err
	}

	// Сначала проверяем, есть ли у пользователя эта вещь
	exists, err := r.ExistsClothesUser(userID, clothID)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	// Обновляем экипированную вещь
	_, err = r.db.Exec("UPDATE users SET equipped_clothes_id = $1 WHERE id = $2", clothID, userID)
	return err
}

func (r *ClothesRepository) GetClotheByUrlImage(urlImage string) (*ClotheRepo, error) {
	var cloth ClotheRepo
	err := r.db.QueryRow("SELECT id, name, description, url_image, price, type, rarity, per_for_tap, plus_energy FROM clothes WHERE url_image = $1", urlImage).
		Scan(&cloth.ID, &cloth.Name, &cloth.Description, &cloth.UrlImage, &cloth.Price, &cloth.Type, &cloth.Rarity, &cloth.PerForTap, &cloth.PlusEnergy)
	if err != nil {
		return nil, err
	}
	return &cloth, nil
}

func (r *ClothesRepository) UpdateUserHead(userID int, urlImage string) error {
	_, err := r.db.Exec("UPDATE users SET head = $1 WHERE id = $2", urlImage, userID)
	return err
}

func (r *ClothesRepository) UpdateUserBody(userID int, urlImage string) error {
	_, err := r.db.Exec("UPDATE users SET body = $1 WHERE id = $2", urlImage, userID)
	return err
}

func (r *ClothesRepository) UpdateUserLegs(userID int, urlImage string) error {
	_, err := r.db.Exec("UPDATE users SET legs = $1 WHERE id = $2", urlImage, userID)
	return err
}

func (r *ClothesRepository) UpdateUserFoot(userID int, urlImage string) error {
	_, err := r.db.Exec("UPDATE users SET foot = $1 WHERE id = $2", urlImage, userID)
	return err
}

func (r *ClothesRepository) UpdateUserHand(userID int, urlImage string) error {
	_, err := r.db.Exec("UPDATE users SET hand = $1 WHERE id = $2", urlImage, userID)
	return err
}

func (r *ClothesRepository) UpdateUserEnergy(userID int, energy int) error {
	_, err := r.db.Exec("UPDATE users SET max_energy = $1 WHERE id = $2", energy, userID)
	return err
}

func (r *ClothesRepository) UpdateUserPerForTap(userID int, perForTap int) error {
	_, err := r.db.Exec("UPDATE users SET profit_for_tap = $1 WHERE id = $2", perForTap, userID)
	return err
}


