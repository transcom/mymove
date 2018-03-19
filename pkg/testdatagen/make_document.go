package testdatagen

import (
	"log"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeDocument creates a single Document.
func MakeDocument(db *pop.Connection, move *models.Move) (models.Document, error) {
	if move == nil {
		newMove, err := MakeMove(db)
		if err != nil {
			log.Panic(err)
		}
		move = &newMove
	}

	document := models.Document{
		UploaderID: move.UserID,
		MoveID:     move.ID,
	}

	verrs, err := db.ValidateAndSave(&document)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return document, err
}

func MakeDocumentData(db *pop.Connection) {
	moveList := []models.Move{}
	err := db.All(&moveList)
	if err != nil {
		log.Panic(err)
	}

	for _, move := range moveList {
		document, err := MakeDocument(db, &move)
		if err != nil {
			log.Panic(err)
		}
		_, err = MakeUpload(db, &document)
		if err != nil {
			log.Panic(err)
		}
	}
}
