package logsService

import "hezzl/internal/core"

type LogsRepository interface {
	ProcessBatch(batch []core.Item) error
}

type LogsService struct {
	LogsRepository
}

func NewLogsService(l LogsRepository) *LogsService {
	return &LogsService{
		LogsRepository: l,
	}
}

func (l LogsService) ProcessBatch(batch []core.Item) error {
	return l.LogsRepository.ProcessBatch(batch)
}
