package dto

import (
	"errors"
	"strconv"

	"github.com/google/uuid"

	"task-management-platform/backend/internal/repository"
)

var ErrInvalidQuery = errors.New("invalid query params")

func ParseTaskFilters(
	projectID uuid.UUID,
	status string,
	assignee string,
	limitStr string,
	offsetStr string,
) (repository.TaskFilters, error) {

	filters := repository.TaskFilters{
		ProjectID: &projectID,
		Limit:     20,
		Offset:    0,
	}

	if status != "" {
		filters.Status = &status
	}

	if assignee != "" {
		id, err := uuid.Parse(assignee)
		if err != nil {
			return filters, ErrInvalidQuery
		}
		filters.AssigneeID = &id
	}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			return filters, ErrInvalidQuery
		}
		if limit > 100 {
			limit = 100
		}
		filters.Limit = limit
	}

	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return filters, ErrInvalidQuery
		}
		filters.Offset = offset
	}

	return filters, nil
}
