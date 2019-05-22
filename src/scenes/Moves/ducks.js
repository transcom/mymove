import { get, head, pick } from 'lodash';
import { GET_LOGGED_IN_USER } from 'shared/Data/users';
import { fetchActive } from 'shared/utils';

import * as ReduxHelpers from 'shared/ReduxHelpers';
// Types
export const getMoveType = 'GET_MOVE';
export const GET_MOVE = ReduxHelpers.generateAsyncActionTypes(getMoveType);

export const createOrUpdateMoveType = 'CREATE_OR_UPDATE_MOVE';
export const CREATE_OR_UPDATE_MOVE = ReduxHelpers.generateAsyncActionTypes(createOrUpdateMoveType);

// Reducer
const initialState = {
  currentMove: null,
  latestMove: null,
  pendingMoveType: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
  error: null,
};
function reshapeMove(move) {
  if (!move) return null;
  return pick(move, ['id', 'locator', 'orders_id', 'selected_move_type', 'status']);
}
export function moveReducer(state = initialState, action) {
  switch (action.type) {
    case GET_LOGGED_IN_USER.success:
      const lastOrdersMoves = get(action.payload, 'service_member.orders.0.moves', []);
      const activeOrders = fetchActive(get(action.payload, 'service_member.orders'));

      const activeMove = fetchActive(get(activeOrders, 'moves'));

      return Object.assign({}, state, {
        latestMove: reshapeMove(head(lastOrdersMoves)),
        currentMove: reshapeMove(activeMove),
        hasLoadError: false,
        hasLoadSuccess: true,
      });
    case CREATE_OR_UPDATE_MOVE.success:
      return Object.assign({}, state, {
        currentMove: reshapeMove(action.payload),
        latestMove: null,
        pendingMoveType: null,
        hasSubmitSuccess: true,
        hasSubmitError: false,
        error: null,
      });
    case CREATE_OR_UPDATE_MOVE.failure:
      return Object.assign({}, state, {
        currentMove: {},
        latestMove: null,
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case GET_MOVE.success:
      return Object.assign({}, state, {
        currentMove: reshapeMove(action.payload),
        latestMove: null,
        hasLoadSuccess: true,
        hasLoadError: false,
        error: null,
      });
    case GET_MOVE.failure:
      return Object.assign({}, state, {
        currentMove: {},
        latestMove: null,
        hasLoadSuccess: false,
        hasLoadError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
