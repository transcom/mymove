import * as ReduxHelpers from 'shared/ReduxHelpers';

const updateShipmentType = 'UPDATE_SHIPMENT';

const UpdateShipmentPlaceholder = (moveId, shipmentId) => {}; // placeholder only until update is implemented

export function updateShipment(moveId, shipmentId) {
  return function(dispatch) {
    const action = ReduxHelpers.generateAsyncActions(updateShipmentType);
    dispatch(action.start());
    return UpdateShipmentPlaceholder(moveId, shipmentId)
      .then(item => dispatch(action.success(item)))
      .catch(error => dispatch(action.failure(error)));
  };
}
