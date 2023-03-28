package factory

import (
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type moveBuildType byte

const (
	moveBuildBasic moveBuildType = iota
	moveBuildWithoutMoveType
)

func buildMoveWithBuildType(db *pop.Connection, customs []Customization, traits []Trait, buildType moveBuildType) models.Move {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cMove models.Move
	if result := findValidCustomization(customs, Move); result != nil {
		cMove = result.Model.(models.Move)

		if result.LinkOnly {
			return cMove
		}
	}

	order := BuildOrder(db, customs, traits)

	var defaultReferenceID string
	var err error
	if db != nil {
		defaultReferenceID, err = models.GenerateReferenceID(db)
		if err != nil {
			log.Panic(err)
		}
	}
	var moveType *models.SelectedMoveType
	var ppmType *string

	// only set these for basic builds
	if buildType == moveBuildBasic {
		ppmMoveType := models.SelectedMoveTypePPM
		moveType = &ppmMoveType
		partialType := "PARTIAL"
		ppmType = &partialType
	}
	contractor := BuildContractor(db, customs, traits)
	defaultShow := true

	// customize here as MergeModels does not handle pointer
	// customization of booleans correctly
	if cMove.Show != nil {
		defaultShow = *cMove.Show
	}
	defaultLocator := models.GenerateLocator()

	move := models.Move{
		Orders:           order,
		OrdersID:         order.ID,
		SelectedMoveType: moveType,
		PPMType:          ppmType,
		Status:           models.MoveStatusDRAFT,
		Locator:          defaultLocator,
		Show:             &defaultShow,
		Contractor:       &contractor,
		ContractorID:     &contractor.ID,
		ReferenceID:      &defaultReferenceID,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&move, cMove)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &move)
	}

	return move
}

func BuildMove(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	return buildMoveWithBuildType(db, customs, traits, moveBuildBasic)
}

func BuildMoveWithoutMoveType(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	return buildMoveWithBuildType(db, customs, traits, moveBuildWithoutMoveType)
}

func BuildAvailableMove(db *pop.Connection) models.Move {
	now := time.Now()
	return BuildMove(db, []Customization{
		{
			Model: models.Move{
				AvailableToPrimeAt: &now,
				Status:             models.MoveStatusAPPROVED,
			},
		},
	}, nil)
}

func BuildStubbedMoveWithStatus(status models.MoveStatus) models.Move {
	return BuildMove(nil, []Customization{
		{
			Model: models.Entitlement{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
		{
			Model: models.DutyLocation{
				ID: uuid.Must(uuid.NewV4()),
			},
			Type: &DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Order{
				ID:               uuid.Must(uuid.NewV4()),
				UploadedOrdersID: uuid.Must(uuid.NewV4()),
			},
		},
		{
			Model: models.Move{
				ID:     uuid.Must(uuid.NewV4()),
				Status: status,
			},
		},
	}, nil)
}

func buildMoveWithOverrides(db *pop.Connection, customs []Customization, traits []Trait, overrides models.Move) models.Move {
	customs = setupCustomizations(customs, traits)
	// Find move assertion and apply approvals
	if result := findValidCustomization(customs, Move); result != nil {
		if result.LinkOnly {
			log.Fatal("Cannot create overrides Move with LinkOnly Move")
		}
		cMove := result.Model.(models.Move)

		// now override for approvals
		testdatagen.MergeModels(&cMove, overrides)
		result.Model = cMove
	} else {
		customs = append(customs, Customization{
			Model: overrides,
		})
	}

	return BuildMove(db, customs, traits)
}

func BuildApprovalsRequestedMove(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	now := time.Now()
	approvalsRequestedMove := models.Move{
		AvailableToPrimeAt: &now,
		Status:             models.MoveStatusAPPROVALSREQUESTED,
	}
	return buildMoveWithOverrides(db, customs, traits, approvalsRequestedMove)
}

func BuildNeedsServiceCounselingMove(db *pop.Connection) models.Move {
	return buildMoveWithOverrides(db, nil, nil, models.Move{
		Status: models.MoveStatusNeedsServiceCounseling,
	})
}

func BuildServiceCounselingCompletedMove(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	now := time.Now()
	scCompletedMove := models.Move{
		ServiceCounselingCompletedAt: &now,
		Status:                       models.MoveStatusServiceCounselingCompleted,
	}
	return buildMoveWithOverrides(db, customs, traits, scCompletedMove)
}
