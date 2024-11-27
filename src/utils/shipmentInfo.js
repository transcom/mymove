import MOVE_STATUSES from 'constants/moves';

const determineShipmentInfo = (move, mtoShipments) => {
  const isMoveDraft = move.status === MOVE_STATUSES.DRAFT;

  const ppmCount = move.personally_procured_moves?.length || 0;

  const mtoCount = mtoShipments?.length || 0;

  const existingShipmentCount = ppmCount + mtoCount;

  return {
    hasShipment: existingShipmentCount > 0,
    isHHGSelectable: isMoveDraft,
    isNTSSelectable: isMoveDraft,
    isNTSRSelectable: isMoveDraft,
    isPPMSelectable: ppmCount === 0,
    isBoatSelectable: isMoveDraft,
    isMobileHomeSelectable: isMoveDraft,
    isUBSelectable: isMoveDraft,
    shipmentNumber: existingShipmentCount + 1,
  };
};

export default determineShipmentInfo;
