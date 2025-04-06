package redis_service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ruziba3vich/itv_test_project/internal/models"
	"github.com/ruziba3vich/itv_test_project/pkg/logger"
)

// RedisService manages movie caching in Redis
type RedisService struct {
	client *redis.Client  // Redis client for interacting with the database
	log    *logger.Logger // Logger for tracking operations
	ttl    time.Duration  // Time-to-live for cached entries
}

// NewRedisService creates a new RedisService instance
func NewRedisService(client *redis.Client, log *logger.Logger, ttl time.Duration) *RedisService {
	return &RedisService{
		client: client,
		log:    log,
		ttl:    ttl,
	}
}

// SetMovie stores a movie in Redis with a TTL
func (s *RedisService) SetMovie(ctx context.Context, movie *models.Movie) error {
	// Generate the Redis key (e.g., "movie:123")
	key := fmt.Sprintf("movie:%d", movie.ID)

	// Serialize the movie to JSON
	data, err := json.Marshal(movie)
	if err != nil {
		s.log.Error("Failed to marshal movie for Redis", map[string]any{
			"error":    err.Error(),
			"movie_id": movie.ID,
		})
		return fmt.Errorf("failed to marshal movie: %s", err.Error())
	}

	// Set the movie in Redis with TTL
	err = s.client.Set(ctx, key, data, s.ttl).Err()
	if err != nil {
		s.log.Error("Failed to set movie in Redis", map[string]any{
			"error":    err.Error(),
			"movie_id": movie.ID,
		})
		return fmt.Errorf("failed to set movie in Redis: %s", err.Error())
	}

	s.log.Info("Movie cached in Redis", map[string]any{
		"movie_id": movie.ID,
		"ttl":      s.ttl.String(),
	})
	return nil
}

// RemoveMovie deletes a movie from Redis by ID
func (s *RedisService) RemoveMovie(ctx context.Context, id uint) error {
	// Generate the Redis key
	key := fmt.Sprintf("movie:%d", id)

	// Delete the movie from Redis
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		s.log.Error("Failed to remove movie from Redis", map[string]any{
			"error":    err.Error(),
			"movie_id": id,
		})
		return fmt.Errorf("failed to remove movie from Redis: %s", err.Error())
	}

	s.log.Info("Movie removed from Redis", map[string]any{
		"movie_id": id,
	})
	return nil
}

// GetMovie retrieves a movie from Redis by ID (optional, for completeness)
func (s *RedisService) GetMovie(ctx context.Context, id uint) (*models.Movie, error) {
	// Generate the Redis key
	key := fmt.Sprintf("movie:%d", id)

	// Get the movie data from Redis
	data, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		s.log.Info("Movie not found in Redis", map[string]any{
			"movie_id": id,
		})
		return nil, nil // Cache miss
	}
	if err != nil {
		s.log.Error("Failed to get movie from Redis", map[string]any{
			"error":    err.Error(),
			"movie_id": id,
		})
		return nil, fmt.Errorf("failed to get movie from Redis: %s", err.Error())
	}

	// Deserialize the movie
	var movie models.Movie
	if err := json.Unmarshal(data, &movie); err != nil {
		s.log.Error("Failed to unmarshal movie from Redis", map[string]any{
			"error":    err.Error(),
			"movie_id": id,
		})
		return nil, fmt.Errorf("failed to unmarshal movie: %s", err.Error())
	}

	s.log.Info("Movie retrieved from Redis", map[string]any{
		"movie_id": id,
	})
	return &movie, nil
}
