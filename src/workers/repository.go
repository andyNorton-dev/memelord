package workers

import (
	"database/sql"
)

type WorkerRepository struct {
	db *sql.DB
}

func NewWorkerRepository(db *sql.DB) *WorkerRepository {
	return &WorkerRepository{db: db}
}

func (r *WorkerRepository) GetWorkers() ([]WorkerRepo, error) {
	rows, err := r.db.Query("SELECT id, name, description, url_image, type FROM workers WHERE type = 'worker'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	workers := []WorkerRepo{}
	for rows.Next() {
		var worker WorkerRepo
		err := rows.Scan(&worker.ID, &worker.Name, &worker.Description, &worker.UrlImage, &worker.Type)
		if err != nil {
			return nil, err
		}
		workers = append(workers, worker)
	}
	return workers, nil
}

func (r *WorkerRepository) GetArmy() ([]WorkerRepo, error) {
	rows, err := r.db.Query("SELECT id, name, description, url_image, type FROM workers WHERE type = 'army'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	army := []WorkerRepo{}
	for rows.Next() {
		var worker WorkerRepo
		err := rows.Scan(&worker.ID, &worker.Name, &worker.Description, &worker.UrlImage, &worker.Type)
		if err != nil {
			return nil, err
		}
		army = append(army, worker)
	}
	return army, nil
}

func (r *WorkerRepository) GetWorkerById(id int) (WorkerRepo, error) {
	var worker WorkerRepo
	err := r.db.QueryRow("SELECT id, name, description, url_image, type FROM workers WHERE id = $1", id).Scan(&worker.ID, &worker.Name, &worker.Description, &worker.UrlImage, &worker.Type)
	if err != nil {
		return WorkerRepo{}, err
	}
	return worker, nil
}
func (r *WorkerRepository) GetUserWorker(worker_id int, user_id int) (UserWorkerRepo, error) {
	var userWorker UserWorkerRepo
	err := r.db.QueryRow("SELECT id, id_worker, id_upgrade FROM user_workers WHERE id_worker = $1 AND id_user = $2", worker_id, user_id).Scan(&userWorker.ID, &userWorker.IdWorker, &userWorker.IdUpgrade)
	if err != nil {
		return UserWorkerRepo{}, err
	}
	return userWorker, nil
}

func (r *WorkerRepository) GetWorkerUpgrades(idWorker int, level int) (WorkerUpgradeRepo, error) {
	var workerUpgrade WorkerUpgradeRepo
	err := r.db.QueryRow("SELECT id, id_worker, level, cost, profit FROM workers_upgrade WHERE id_worker = $1 AND level = $2", idWorker, level).Scan(&workerUpgrade.ID, &workerUpgrade.IdWorker, &workerUpgrade.Level, &workerUpgrade.Cost, &workerUpgrade.Profit)
	if err != nil {
		return WorkerUpgradeRepo{}, err
	}
	return workerUpgrade, nil
}

func (r *WorkerRepository) GetUpgradeById(id int) (WorkerUpgradeRepo, error) {
	var workerUpgrade WorkerUpgradeRepo
	err := r.db.QueryRow("SELECT id, id_worker, level, cost, profit FROM workers_upgrade WHERE id = $1", id).Scan(&workerUpgrade.ID, &workerUpgrade.IdWorker, &workerUpgrade.Level, &workerUpgrade.Cost, &workerUpgrade.Profit)
	if err != nil {
		return WorkerUpgradeRepo{}, err
	}
	return workerUpgrade, nil
}

func (r *WorkerRepository) UpdateUserBalance(tg_id int, balance int) error {
	_, err := r.db.Exec("UPDATE users SET balance = $1 WHERE tg_id = $2", balance, tg_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *WorkerRepository) CreateUserWorker(tg_id int, worker_id int, upgrade_id int) error {
	_, err := r.db.Exec("INSERT INTO user_workers (id_user, id_worker, id_upgrade) VALUES ($1, $2, $3)", tg_id, worker_id, upgrade_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *WorkerRepository) UpdateUserWorker(upgrade_id int, user_worker_id int) error {
	_, err := r.db.Exec("UPDATE user_workers SET id_upgrade = $1 WHERE id = $2", upgrade_id, user_worker_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *WorkerRepository) UpdateUserWorkerProfit(id int, profit int, balance int) error {
	_, err := r.db.Exec("UPDATE users SET profit_per_hour = profit_per_hour + $1, balance = balance - $2 WHERE id = $3", profit, balance, id)
	if err != nil {
		return err
	}
	return nil
}
