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

		if newWeightTicket.AdjustedNetWeight == nil || *newWeightTicket.AdjustedNetWeight < 0 {
			verrs.Add("AdjustedNetWeight", "Adjusted Net Weight must have a value of at least 0")
		}

		if newWeightTicket.FullWeight != nil && newWeightTicket.AdjustedNetWeight != nil && *newWeightTicket.AdjustedNetWeight >= *newWeightTicket.FullWeight {
			verrs.Add("AdjustedNetWeight", "Adjusted Net Weight cannot be greater than the full weight")
		}

		if newWeightTicket.NetWeightRemarks == nil || *newWeightTicket.NetWeightRemarks == "" {
			verrs.Add("NetWeightRemarks", "Net Weight Remarks must exist")
		}

		return verrs
	})
}

func verifyProofOfTrailerOwnershipDocument() weightTicketValidator {
	return weightTicketValidatorFunc(func(_ appcontext.AppContext, newWeightTicket *models.WeightTicket, originalWeightTicket *models.WeightTicket) error {
		verrs := validate.NewErrors()

		if newWeightTicket == nil || originalWeightTicket == nil {
			return verrs
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

		if (originalWeightTicket.Status == nil && newWeightTicket.Status != nil) ||
			(originalWeightTicket.Status != nil && newWeightTicket.Status == nil) ||
			(originalWeightTicket.Status != nil && newWeightTicket.Status != nil && *originalWeightTicket.Status != *newWeightTicket.Status) {
			verrs.Add("Status", "status cannot be modified")
		}

		if (originalWeightTicket.Reason == nil && newWeightTicket.Reason != nil) ||
			(originalWeightTicket.Reason != nil && newWeightTicket.Reason == nil) ||
			(originalWeightTicket.Reason != nil && newWeightTicket.Reason != nil && *originalWeightTicket.Reason != *newWeightTicket.Reason) {
			verrs.Add("Reason", "reason cannot be modified")
		}

		return verrs
	})
}

func verifyReasonAndStatusAreValid() weightTicketValidator {
	return weightTicketValidatorFunc(func(_ appcontext.AppContext, newWeightTicket *models.WeightTicket, originalWeightTicket *models.WeightTicket) error {
		verrs := validate.NewErrors()

		if newWeightTicket.Status != nil {
			if *newWeightTicket.Status == models.PPMDocumentStatusApproved && newWeightTicket.Reason != nil {
				verrs.Add("Reason", "reason must not be set if the status is Approved")
			} else if (*newWeightTicket.Status == models.PPMDocumentStatusExcluded || *newWeightTicket.Status == models.PPMDocumentStatusRejected) &&
				(newWeightTicket.Reason == nil || *newWeightTicket.Reason == "") {
				verrs.Add("Reason", "reason is mandatory if the status is Excluded or Rejected")
			}
		} else if newWeightTicket.Reason != nil {
			verrs.Add("Reason", "reason should not be set if the status is not set")
		}

		return verrs
	})
}

func basicChecksForCreate() []weightTicketValidator {
	return []weightTicketValidator{
		checkID(),
		checkRequiredFields(),
	}
}

func basicChecksForCustomer() []weightTicketValidator {
	return []weightTicketValidator{
		checkID(),
		checkRequiredFields(),
		verifyProofOfTrailerOwnershipDocument(),
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
