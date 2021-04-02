import { SHIPMENT_OPTIONS } from 'shared/constants';
import MOVE_STATUSES from 'constants/moves';

const determineShipmentInfo = (move, mtoShipments) => {
  const isMoveDraft = move.status === MOVE_STATUSES.DRAFT;

  const ppmCount = move.personally_procured_moves?.length || 0;

  const mtoCount = mtoShipments?.length || 0;

  const hasNTS = mtoShipments.some((shipment) => shipment.shipmentType === SHIPMENT_OPTIONS.NTS);

  const hasNTSR = mtoShipments.some((shipment) => shipment.shipmentType === SHIPMENT_OPTIONS.NTSR);

  const existingShipmentCount = ppmCount + mtoCount;

  return {
    hasShipment: existingShipmentCount > 0,
    isHHGSelectable: isMoveDraft,
    isNTSSelectable: isMoveDraft && !hasNTS,
    isNTSRSelectable: isMoveDraft && !hasNTSR,
    isPPMSelectable: ppmCount === 0,
    shipmentNumber: existingShipmentCount + 1,
  };
};

export default determineShipmentInfo;
