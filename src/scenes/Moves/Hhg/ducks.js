import { get } from 'lodash';
import * as ReduxHelpers from 'shared/ReduxHelpers';
import { GET_LOGGED_IN_USER } from 'shared/User/ducks';
import { fetchActive } from 'shared/utils';

// Types

// Action creation

// Selectors

// Reducer
const initialState = {
  currentHhg: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
  hasLoadSuccess: false,
  hasLoadError: false,
};
export function hhgReducer(state = initialState, action) {
  switch (action.type) {
    case GET_LOGGED_IN_USER.success:
      // Initialize state when we get the logged in user
      const activeOrders = fetchActive(
        get(action.payload, 'service_member.orders'),
      );
      const activeMove = fetchActive(get(activeOrders, 'moves'));
      const activeHhg = fetchActive(get(activeMove, 'shipments'));
      return Object.assign({}, state, {
        currentHhg: activeHhg,
        hasLoadSuccess: true,
        hasLoadError: false,
      });
    default:
      return state;
  }
}
