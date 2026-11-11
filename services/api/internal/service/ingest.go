package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/myntdeveloper/agromir_boiler_server/internal/model"
	"github.com/myntdeveloper/agromir_boiler_server/internal/repo"
)

type UnknownBoilersError struct {
	Numbers []int
}

func (e *UnknownBoilersError) Error() string {
	return fmt.Sprintf("unknown boiler_number: %v", e.Numbers)
}

type IngestService struct {
	repo *repo.IngestRepo
}

func NewIngestService(r *repo.IngestRepo) *IngestService {
	return &IngestService{repo: r}
}

func (s *IngestService) Ingest(ctx context.Context, req model.IngestRequest) (int, error) {
	fillMissingTimestamps(&req)

	numbers := uniqueBoilerNumbers(req.Readings)
	ids, err := s.repo.MapBoilerNumbers(ctx, numbers)
	if err != nil {
		return 0, err
	}

	missing := missingNumbers(numbers, ids)
	if len(missing) > 0 {
		return 0, &UnknownBoilersError{Numbers: missing}
	}

	if err := s.repo.SavePacket(ctx, req, ids); err != nil {
		return 0, err
	}
	return len(req.Readings), nil
}

func (s *IngestService) PlantStatus(ctx context.Context, plantID string) (*model.PlantStatus, error) {
	status, err := s.repo.GetPlantStatus(ctx, plantID)
	if err != nil {
		return nil, err
	}
	if status == nil {
		return nil, ErrPlantNotFound
	}
	return status, nil
}

// PLC clocks are often unsynced, so readings without ts get sent_at or server time.
func fillMissingTimestamps(req *model.IngestRequest) {
	fallback := req.SentAt
	if fallback.IsZero() {
		fallback = time.Now().UTC().Truncate(time.Second)
	}
	for i := range req.Readings {
		if req.Readings[i].TS.IsZero() {
			req.Readings[i].TS = fallback
		}
	}
}

func uniqueBoilerNumbers(readings []model.IngestReading) []int {
	seen := make(map[int]struct{}, len(readings))
	out := make([]int, 0, 24)
	for _, reading := range readings {
		if _, ok := seen[reading.BoilerNumber]; ok {
			continue
		}
		seen[reading.BoilerNumber] = struct{}{}
		out = append(out, reading.BoilerNumber)
	}
	return out
}

func missingNumbers(wanted []int, found map[int]int) []int {
	missing := make([]int, 0)
	for _, number := range wanted {
		if _, ok := found[number]; !ok {
			missing = append(missing, number)
		}
	}
	sort.Ints(missing)
	return missing
}
