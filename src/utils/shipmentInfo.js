import { MOVE_STATUSES, SHIPMENT_OPTIONS } from 'shared/constants';

const determineShipmentInfo = (move, mtoShipments) => {
  const isMoveDraft = move.status === MOVE_STATUSES.DRAFT;

  const hasPPM = Boolean(move.personally_procured_moves?.length);

  const ppmCount = hasPPM ? 1 : 0;
  const mtoCount = mtoShipments?.length || 0;

  const hasNTS = mtoShipments.some((shipment) => shipment.shipmentType === SHIPMENT_OPTIONS.NTS);

  const hasNTSR = mtoShipments.some((shipment) => shipment.shipmentType === SHIPMENT_OPTIONS.NTSR);

  return {
    isHHGSelectable: isMoveDraft,
    isNTSSelectable: isMoveDraft && !hasNTS,
    isNTSRSelectable: isMoveDraft && !hasNTSR,
    isPPMSelectable: !hasPPM,
    shipmentNumber: ppmCount + mtoCount + 1,
  };
};

export default determineShipmentInfo;
