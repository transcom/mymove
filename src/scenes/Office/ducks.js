import { get } from 'lodash';
import { loadMove, selectMove } from 'shared/Entities/modules/moves';
import { loadPPMs } from 'shared/Entities/modules/ppms';
import { loadServiceMember, loadBackupContacts } from 'shared/Entities/modules/serviceMembers';
import { loadOrders } from 'shared/Entities/modules/orders';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// SINGLE RESOURCE ACTION TYPES
const REMOVE_BANNER = 'REMOVE_BANNER';
const SHOW_BANNER = 'SHOW_BANNER';

// MULTIPLE-RESOURCE ACTION TYPES
const loadDependenciesType = 'LOAD_DEPENDENCIES';

// MULTIPLE-RESOURCE ACTION TYPES
const LOAD_DEPENDENCIES = ReduxHelpers.generateAsyncActionTypes(loadDependenciesType);

// SINGLE-RESOURCE ACTION CREATORS

export const removeBanner = () => {
  return {
    type: REMOVE_BANNER,
  };
};

export const showBanner = ({ messageLines }) => {
  return {
    type: SHOW_BANNER,
    payload: { messageLines },
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
      const move = selectMove(getState(), moveId);
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

// Reducer
const initialState = {
  loadDependenciesHasError: null,
  loadDependenciesHasSuccess: false,
  flashMessage: false,
  flashMessageLines: [],
};

export function officeReducer(state = initialState, action) {
  switch (action.type) {
    // SINGLE-RESOURCE ACTION TYPES
    case SHOW_BANNER:
      return Object.assign({}, state, {
        flashMessage: true,
        flashMessageLines: action.payload.messageLines,
      });
    case REMOVE_BANNER:
      return Object.assign({}, state, {
        flashMessage: false,
        flashMessageLines: [],
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
