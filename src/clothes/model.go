package clothes

type ClotheRepo struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	UrlImage string `json:"url_image"`
	Price int `json:"price"`
	Type string `json:"type"`
	Rarity string `json:"rarity"`
	PerForTap int `json:"per_for_tap"`
	PlusEnergy int `json:"plus_energy"`
}

type ClothesRepo struct {
	ID int `json:"id"`
	UrlImage string `json:"url_image"`
	Price int `json:"price"`
	Type string `json:"type"`
	Rarity string `json:"rarity"`
}
type ClothesUserRepo struct {
	ID int `json:"id"`
	ClothesID int `json:"clothes_id"`
	UserID int `json:"user_id"`
}

type ClotheUserResponse struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	UrlImage string `json:"url_image"`
	Price int `json:"price"`
	Rarity string `json:"rarity"`
	PerForTap int `json:"per_for_tap"`
	PlusEnergy int `json:"plus_energy"`
	Type string `json:"type"`
	IsActive bool `json:"is_active"`
	IsBought bool `json:"is_bought"`
	CanBuy bool `json:"can_buy"`
}

type ClothesUserResponse struct {
	ID int `json:"id"`
	UrlImage string `json:"url_image"`
	Price int `json:"price"`
	Type string `json:"type"`
	Rarity string `json:"rarity"`
	IsActive bool `json:"is_active"`
	IsBought bool `json:"is_bought"`
}