package workers

type WorkerRepo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	UrlImage string `json:"url_image"`
	Type string `json:"type"`

}

type WorkerUpgradeRepo struct {
	ID   int    `json:"id"`
	IdWorker int `json:"id_worker"`
	Level int `json:"level"`
	Cost int `json:"cost"`
	Profit int `json:"profit"`
}

type UserWorkerRepo struct {
	ID   int    `json:"id"`
	IdWorker int `json:"id_worker"`
	IdUpgrade int `json:"id_upgrade"`
}

type UserWorkerResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	UrlImage string `json:"url_image"`
	Level int `json:"level"`
	Profit int `json:"profit"`
	Cost int `json:"cost"`
	AccessToUpgrade bool `json:"access_to_upgrade"`
	
}