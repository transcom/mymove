package paperwork

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

// Generator encapsulates the prerequisites for PDF generation.
type Generator struct {
	db     *pop.Connection
	logger *zap.Logger
}

// NewGenerator creates a new Generator.
func NewGenerator(db *pop.Connection, logger *zap.Logger) *Generator {
	return &Generator{
		db:     db,
		logger: logger,
	}
}

// GenerateAdvancePaperwork generates the advance paperwork for a move.
func (g *Generator) GenerateAdvancePaperwork(moveID uuid.UUID) error {
	move, err := models.FetchMoveForAdvancePaperwork(g.db, moveID)
	if err != nil {
		return err
	}
	fmt.Println(move)
	return nil
}
