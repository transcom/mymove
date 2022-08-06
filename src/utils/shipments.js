// The PPM shipment creation is a multi-step flow so it's possible to get in a state with missing
// information and get to the review screen in an incomplete state from creating another shipment
// on the move. hasRequestedAdvance is the last required field that would mean they're finished.
export function isPPMShipmentComplete(mtoShipment) {
  if (mtoShipment?.ppmShipment?.hasRequestedAdvance != null) {
    return true;
  }
  return false;
}

// isWeightTicketComplete - checks that the required fields for a weight ticket have valid data
// to check if the weight ticket can be considered complete. For the purposes of this function,
// any data is enough to consider some fields valid.
export function isWeightTicketComplete(weightTicket) {
  const hasValidEmptyWeight = weightTicket.emptyWeight != null && weightTicket.emptyWeight >= 0;

  const hasTrailerDocUpload = weightTicket.proofOfTrailerOwnershipDocument.uploads.length > 0;
  const needsTrailerUpload = weightTicket.ownsTrailer && weightTicket.trailerMeetsCriteria;
  const trailerNeedsMet = needsTrailerUpload ? hasTrailerDocUpload : true;

  return !!(
    weightTicket.vehicleDescription &&
    hasValidEmptyWeight &&
    weightTicket.emptyDocument.uploads.length > 0 &&
    weightTicket.fullWeight > 0 &&
    weightTicket.fullDocument.uploads.length > 0 &&
    trailerNeedsMet
  );
}

export function hasCompleteWeightTickets(weightTickets) {
  if (!weightTickets?.length) {
    return false;
  }

  return !!weightTickets?.every((weightTicket) => {
    return weightTicket.vehicleDescription && weightTicket.emptyWeight && weightTicket.fullWeight;
  });
}

export default isPPMShipmentComplete;
