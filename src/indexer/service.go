package indexer

import (
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type IndexService struct {
	repo *IndexRepository
	mu   sync.Mutex
}

func NewIndexService(repo *IndexRepository) *IndexService {
	return &IndexService{
		repo: repo,
		mu:   sync.Mutex{},
	}
}

func (s *IndexService) GetCurrentIndex(c echo.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	index, err := s.repo.GetCurrentIndex()
	if err != nil {
		log.Printf("Ошибка при получении индекса: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "ошибка при получении индекса",
		})
	}
	
	log.Printf("Получен запрос на получение текущего индекса. Текущее значение: %d", index)
	return c.JSON(http.StatusOK, map[string]interface{}{"index": index})
}

func (s *IndexService) IncrementIndex(c echo.Context) error {
	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		
		index, err := s.repo.GetCurrentIndex()
		if err != nil {
			log.Printf("Ошибка при получении индекса: %v", err)
			return
		}
		
		oldValue := index
		newValue := index + 1
		
		err = s.repo.UpdateIndex(newValue)
		if err != nil {
			log.Printf("Ошибка при обновлении индекса: %v", err)
			return
		}
		
		log.Printf("Индекс увеличен с %d до %d", oldValue, newValue)
	}()
	
	log.Printf("Запущена асинхронная операция увеличения индекса")
	return c.JSON(http.StatusNoContent, nil)
}

func (s *IndexService) DoubleIndex(c echo.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	index, err := s.repo.GetCurrentIndex()
	if err != nil {
		log.Printf("Ошибка при получении индекса: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "ошибка при получении индекса",
		})
	}
	
	if index%2 == 0 {
		log.Printf("Попытка удвоить четное число %d отклонена", index)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "нельзя удвоить четное число",
		})
	}
	
	oldValue := index
	newValue := index * 2
	
	err = s.repo.UpdateIndex(newValue)
	if err != nil {
		log.Printf("Ошибка при обновлении индекса: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "ошибка при обновлении индекса",
		})
	}
	
	log.Printf("Индекс удвоен с %d до %d", oldValue, newValue)
	return c.JSON(http.StatusNoContent, nil)
}

func (s *IndexService) AddToIndex(c echo.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	req, ok := c.Get("request").(*AddToIndexRequest)
	if !ok {
		log.Printf("Ошибка при получении данных из контекста")
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "внутренняя ошибка сервера",
		})
	}

	index, err := s.repo.GetCurrentIndex()
	if err != nil {
		log.Printf("Ошибка при получении индекса: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "ошибка при получении индекса",
		})
	}

	oldValue := index
	newValue := index + req.Number

	err = s.repo.UpdateIndex(newValue)
	if err != nil {
		log.Printf("Ошибка при обновлении индекса: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "ошибка при обновлении индекса",
		})
	}

	log.Printf("К индексу добавлено число %d: %d -> %d", req.Number, oldValue, newValue)
	return c.JSON(http.StatusOK, AddToIndexResponse{
		OldValue: oldValue,
		NewValue: newValue,
		Added:    req.Number,
	})
}
