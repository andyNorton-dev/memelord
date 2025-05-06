package indexer

import (
	"database/sql"
)

type IndexRepository struct {
	db *sql.DB
}

func NewIndexRepository(db *sql.DB) *IndexRepository {
	return &IndexRepository{db: db}
}

func (r *IndexRepository) GetCurrentIndex() (int, error) {
	var index int
	err := r.db.QueryRow("SELECT idx FROM my_table LIMIT 1").Scan(&index)
	if err == sql.ErrNoRows {
		// Если запись не существует, создаем её
		_, err = r.db.Exec("INSERT INTO my_table (idx) VALUES (0)")
		if err != nil {
			return 0, err
		}
		return 0, nil
	}
	return index, err
}

func (r *IndexRepository) UpdateIndex(newIndex int) error {
	_, err := r.db.Exec("UPDATE my_table SET idx = $1", newIndex)
	return err
}