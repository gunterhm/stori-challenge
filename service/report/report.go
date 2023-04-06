package report

import (
	"awesomeProject/model"
	"context"
)

// IReport is a service interface for retrieving report information
type IReport interface {
	GetSummaryEmailInfo(ctx context.Context, accountID string) (*model.SummaryEmail, error)
}
