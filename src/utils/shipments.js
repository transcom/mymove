// The PPM shipment creation is a multi-step flow so it's possible to get in a state with missing
// information and get to the review screen in an incomplete state from creating another shipment
// on the move.  Advance requested is the last required field that would mean they're finished.
export function isPPMShipmentComplete(mtoShipment) {
  if (mtoShipment?.ppmShipment?.advanceRequested != null) {
    return false;
  }
  return true;
}

export default isPPMShipmentComplete;
