package weightticket

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

func checkID() weightTicketValidator {
	return weightTicketValidatorFunc(func(_ appcontext.AppContext, newWeightTicket *models.WeightTicket, originalWeightTicket *models.WeightTicket) error {
		verrs := validate.NewErrors()

		if newWeightTicket == nil || originalWeightTicket == nil {
			return verrs
		}

		if newWeightTicket.ID != originalWeightTicket.ID {
			verrs.Add("ID", "new WeightTicket ID must match original WeightTicket ID")
		}

		return verrs
	})
}

func checkRequiredFields() weightTicketValidator {
	return weightTicketValidatorFunc(func(_ appcontext.AppContext, newWeightTicket *models.WeightTicket, originalWeightTicket *models.WeightTicket) error {
		verrs := validate.NewErrors()

		if newWeightTicket == nil || originalWeightTicket == nil {
			return verrs
		}

		if originalWeightTicket.PPMShipmentID.IsNil() {
			verrs.Add("PPMShipmentID", "PPMShipmentID must exist")
		}
		if originalWeightTicket.EmptyDocumentID.IsNil() {
			verrs.Add("EmptyDocumentID", "EmptyDocumentID must exist")
		}
		if originalWeightTicket.FullDocumentID.IsNil() {
			verrs.Add("FullDocumentID", "FullDocumentID must exist")
		}
		if originalWeightTicket.ProofOfTrailerOwnershipDocumentID.IsNil() {
			verrs.Add("ProofOfTrailerOwnershipDocumentID", "ProofOfTrailerOwnershipDocumentID must exist")
		}

		if newWeightTicket.VehicleDescription == nil || *newWeightTicket.VehicleDescription == "" {
			verrs.Add("VehicleDescription", "Vehicle Description must exist")
		}

		if newWeightTicket.EmptyWeight == nil || *newWeightTicket.EmptyWeight < 0 {
			verrs.Add("EmptyWeight", "Empty Weight must have a value of at least 0")
		}

		if newWeightTicket.MissingEmptyWeightTicket == nil {
			verrs.Add("MissingEmptyWeightTicket", "Missing Empty Weight Ticket is required")
		}

		if newWeightTicket.FullWeight == nil || *newWeightTicket.FullWeight < 1 {
			verrs.Add("FullWeight", "Full Weight must have a value of at least 1")
		}

		if newWeightTicket.EmptyWeight != nil && newWeightTicket.FullWeight != nil && *newWeightTicket.FullWeight <= *newWeightTicket.EmptyWeight {
			verrs.Add("FullWeight", "Full Weight must be greater than Empty Weight")
		}

		if newWeightTicket.MissingFullWeightTicket == nil {
			verrs.Add("MissingFullWeightTicket", "Missing Full Weight Ticket is required")
		}

		if len(originalWeightTicket.EmptyDocument.UserUploads) < 1 {
			verrs.Add("EmptyWeightDocument", "At least 1 empty weight file is required")
		}

		if len(originalWeightTicket.FullDocument.UserUploads) < 1 {
			verrs.Add("FullWeightDocument", "At least 1 full weight file is required")
		}

		if newWeightTicket.OwnsTrailer == nil {
			verrs.Add("OwnsTrailer", "Owns Trailer is required")
		}

		if newWeightTicket.TrailerMeetsCriteria == nil {
			verrs.Add("TrailerMeetsCriteria", "Trailer Meets Criteria is required")
		}

		if newWeightTicket.OwnsTrailer != nil && newWeightTicket.TrailerMeetsCriteria != nil {
			if *newWeightTicket.OwnsTrailer && *newWeightTicket.TrailerMeetsCriteria &&
				len(originalWeightTicket.ProofOfTrailerOwnershipDocument.UserUploads) < 1 {
				verrs.Add("ProofOfTrailerOwnershipDocument", "At least 1 proof of ownership file is required")
			}
		}

		return verrs
	})
}

// verifyReasonAndStatusAreConstant ensures that the reason and status fields do not change
func verifyReasonAndStatusAreConstant() weightTicketValidator {
	return weightTicketValidatorFunc(func(_ appcontext.AppContext, newWeightTicket *models.WeightTicket, originalWeightTicket *models.WeightTicket) error {
		verrs := validate.NewErrors()

		if newWeightTicket == nil || originalWeightTicket == nil {
			return verrs
		}

		if originalWeightTicket.Status == nil && newWeightTicket.Status != nil {
			verrs.Add("Status", "status cannot be modified")
		} else if originalWeightTicket.Status != nil && newWeightTicket.Status != nil && *originalWeightTicket.Status != *newWeightTicket.Status {
			verrs.Add("Status", "status cannot be modified")
		}

		if originalWeightTicket.Reason != newWeightTicket.Reason {
			verrs.Add("Reason", "reason cannot be modified")
		}
		return verrs
	})
}

func verifyReasonAndStatusAreValid() weightTicketValidator {
	return weightTicketValidatorFunc(func(_ appcontext.AppContext, newWeightTicket *models.WeightTicket, originalWeightTicket *models.WeightTicket) error {
		verrs := validate.NewErrors()

		if newWeightTicket == nil || originalWeightTicket == nil {
			return verrs
		}

		if newWeightTicket.Status != nil {
			if *newWeightTicket.Status == models.PPMDocumentStatusApproved && newWeightTicket.Reason != nil {
				verrs.Add("Reason", "reason must be blank if the status is Approved")
			}

			if (*newWeightTicket.Status == models.PPMDocumentStatusExcluded || *newWeightTicket.Status == models.PPMDocumentStatusRejected) && (newWeightTicket.Reason == nil || len(*newWeightTicket.Reason) <= 0) {
				verrs.Add("Reason", "reason is mandatory if the status is Excluded or Rejected")
			}
		} else {
			if newWeightTicket.Reason != nil && len(*newWeightTicket.Reason) > 0 {
				verrs.Add("Reason", "reason should be empty")
			}
		}

		return verrs
	})
}

func basicChecksForCustomer() []weightTicketValidator {
	return []weightTicketValidator{
		checkID(),
		checkRequiredFields(),
		verifyReasonAndStatusAreConstant(),
	}
}

func basicChecksForOffice() []weightTicketValidator {
	return []weightTicketValidator{
		checkID(),
		checkRequiredFields(),
		verifyReasonAndStatusAreValid(),
	}
}
