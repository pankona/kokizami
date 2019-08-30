package repo

import (
	"fmt"

	"github.com/pankona/kokizami/models"
)

func CreateTables(db models.XODB) error {
	if err := models.CreateKizamiTable(db); err != nil {
		return fmt.Errorf("failed to create kizami table: %v", err)
	}

	if err := models.CreateTagTable(db); err != nil {
		return fmt.Errorf("failed to create tag table: %v", err)
	}

	if err := models.CreateRelationTable(db); err != nil {
		return fmt.Errorf("failed to create relation table: %v", err)
	}

	return nil
}
