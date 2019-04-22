import { isNull } from 'lodash';
import { getEntitlements } from 'shared/entitlements.js';
import { selectShipment } from 'shared/Entities/modules/shipments';

// SINGLE-RESOURCE ACTION CREATORS

export function loadEntitlements(state, shipmentId) {
  const shipment = selectShipment(state, shipmentId);
  const move = shipment.move || {};
  const serviceMember = shipment.service_member || {};
  const hasDependents = move.has_dependents;
  const spouseHasProGear = move.spouse_has_progear;
  const rank = serviceMember.rank;

  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(rank)) {
    return null;
  }
  return getEntitlements(rank, hasDependents, spouseHasProGear);
}
// Reducer
const initialState = {
  storageInTransitIsCreating: false,
  storageInTransitHasCreatedSuccess: false,
  storageInTransitHasCreatedError: null,
  storageInTransitsAreLoading: false,
  storageInTransitsHasLoadSuccess: false,
  storageInTransitsHasLoadError: null,
  flashMessage: false,
};

export function tspReducer(state = initialState, action) {
  switch (action.type) {
    default:
      return state;
  }
}
