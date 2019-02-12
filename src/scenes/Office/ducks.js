import { isNull, get } from 'lodash';
import { LoadMove } from './api.js';
import { getEntitlements } from 'shared/entitlements.js';
import { loadPPMs } from 'shared/Entities/modules/ppms';
import { loadServiceMember, loadBackupContacts } from 'shared/Entities/modules/serviceMembers';
import { loadOrders, selectOrdersForMove } from 'shared/Entities/modules/orders';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// SINGLE RESOURCE ACTION TYPES
const loadMoveType = 'LOAD_MOVE';
const REMOVE_BANNER = 'REMOVE_BANNER';
const SHOW_BANNER = 'SHOW_BANNER';
const RESET_MOVE = 'RESET_MOVE';

// MULTIPLE-RESOURCE ACTION TYPES
const loadDependenciesType = 'LOAD_DEPENDENCIES';

// SINGLE RESOURCE ACTION TYPES

export const resetMove = () => ({
  type: RESET_MOVE,
});

const LOAD_MOVE = ReduxHelpers.generateAsyncActionTypes(loadMoveType);

// MULTIPLE-RESOURCE ACTION TYPES

const LOAD_DEPENDENCIES = ReduxHelpers.generateAsyncActionTypes(loadDependenciesType);

// SINGLE-RESOURCE ACTION CREATORS

export const loadMove = ReduxHelpers.generateAsyncActionCreator(loadMoveType, LoadMove);

export const removeBanner = () => {
  return {
    type: REMOVE_BANNER,
  };
};

export const showBanner = () => {
  return {
    type: SHOW_BANNER,
  };
};
// MULTIPLE-RESOURCE ACTION CREATORS
//
// These action types typically dispatch to other actions above to
// perform their work and exist to encapsulate when multiple requests
// need to be made in response to a user action.

export function loadMoveDependencies(moveId) {
  const actions = ReduxHelpers.generateAsyncActions(loadDependenciesType);
  return async function(dispatch, getState) {
    dispatch(actions.start());
    try {
      await dispatch(loadMove(moveId));
      const move = getState().office.officeMove;
      const ordersId = move.orders_id;
      await dispatch(loadOrders(ordersId));
      const serviceMemberId = get(getState(), `entities.orders.${ordersId}.service_member_id`);
      await dispatch(loadServiceMember(serviceMemberId));
      await dispatch(loadBackupContacts(serviceMemberId));
      // TODO: load PPMs in parallel to move using moveId
      await dispatch(loadPPMs(moveId));
      return dispatch(actions.success());
    } catch (ex) {
      return dispatch(actions.error(ex));
    }
  };
}

// Selectors
export function loadEntitlements(state, moveId) {
  const orders = selectOrdersForMove(state, moveId);
  const hasDependents = orders.has_dependents;
  const spouseHasProGear = orders.spouse_has_pro_gear;
  const rank = get(state, 'office.officeServiceMember.rank', null);
  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(rank)) {
    return null;
  }
  return getEntitlements(rank, hasDependents, spouseHasProGear);
}

// Reducer
const initialState = {
  moveIsLoading: false,
  moveHasLoadError: null,
  moveHasLoadSuccess: false,
  officeMove: {},
  loadDependenciesHasError: null,
  loadDependenciesHasSuccess: false,
  flashMessage: false,
};

export function officeReducer(state = initialState, action) {
  switch (action.type) {
    // SINGLE-RESOURCE ACTION TYPES

    // MOVES
    case LOAD_MOVE.start:
      return Object.assign({}, state, {
        moveIsLoading: true,
        moveHasLoadSuccess: false,
      });
    case LOAD_MOVE.success:
      return Object.assign({}, state, {
        moveIsLoading: false,
        officeMove: action.payload,
        officeShipment: get(action.payload, 'shipments.0', null),
        moveHasLoadSuccess: true,
        moveHasLoadError: false,
      });
    case LOAD_MOVE.failure:
      return Object.assign({}, state, {
        moveIsLoading: false,
        officeMove: null,
        officeShipment: null,
        moveHasLoadSuccess: false,
        moveHasLoadError: true,
        error: action.error.message,
      });

    case SHOW_BANNER:
      return Object.assign({}, state, {
        flashMessage: true,
      });
    case REMOVE_BANNER:
      return Object.assign({}, state, {
        flashMessage: false,
      });

    // MULTIPLE-RESOURCE ACTION TYPES
    //
    // These action types typically dispatch to other actions above to
    // perform their work and exist to encapsulate when multiple requests
    // need to be made in response to a user action.

    // ALL DEPENDENCIES
    case LOAD_DEPENDENCIES.start:
      return Object.assign({}, state, {
        loadDependenciesHasSuccess: false,
        loadDependenciesHasError: false,
      });
    case LOAD_DEPENDENCIES.success:
      return Object.assign({}, state, {
        loadDependenciesHasSuccess: true,
        loadDependenciesHasError: false,
      });
    case LOAD_DEPENDENCIES.failure:
      return Object.assign({}, state, {
        loadDependenciesHasSuccess: false,
        loadDependenciesHasError: true,
        error: action.error.message,
      });

    default:
      return state;
  }
}
