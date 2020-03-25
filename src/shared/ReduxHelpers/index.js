/* The 4 exported functions in this file help eliminate boilerplate in our redux code when we need actions/action creators/reducers
 * for asychronous api calls.
 */
function actionName(resource, actionType) {
  return `${resource}_${actionType.toUpperCase()}`;
}

function generateActionTypes(resourceName, actionTypes) {
  resourceName = (resourceName || '').trim();
  if (!resourceName) throw new Error('No resource name provided');
  let actions = {};
  actionTypes.forEach(actionType => {
    actions[actionType] = actionName(resourceName, actionType); // eslint-disable-line security/detect-object-injection
  });
  return actions;
}
/**
 *  For `resourceName` FOO, this returns an object with 3 action types:
 *   {
 *     start: 'FOO_START',
 *     success: 'FOO_SUCCESS',
 *     failure: 'FOO_FAILURE'
 *   }
 * @param {*} resourceName
 */
export function generateAsyncActionTypes(resourceName) {
  return generateActionTypes(resourceName, ['start', 'success', 'failure']);
}

/**
 * This creates basic functions for starting an async action and handling the success and failure states
 * @param {*} resourceName
 */
export function generateAsyncActions(resourceName) {
  const actionTypes = generateAsyncActionTypes(resourceName);
  const actions = {
    start: () => ({
      type: actionTypes.start,
    }),
    success: payload => ({
      type: actionTypes.success,
      payload,
    }),
    error: error => ({
      type: actionTypes.failure,
      error,
    }),
  };
  return actions;
}
/**
 *  This returns an action creator that calls the asyncAction.
 * @param {*} resourceName the name of the resource included in the action type descriptions (for start, success, failure)
 * @param {*} asyncAction  the async action that is wrapped by this action creator
 */
export function generateAsyncActionCreator(resourceName, asyncAction) {
  const actions = generateAsyncActions(resourceName);
  return function actionCreator(...args) {
    return function(dispatch) {
      dispatch(actions.start());
      return asyncAction(...args)
        .then(item => dispatch(actions.success(item)))
        .catch(error => {
          dispatch(actions.error(error));
          return Promise.reject(error);
        });
    };
  };
}

/**
 * This produces a reducer for the provided resourceName
 * @param {*} resourceName the name of the resource the actions should be filtered by
 * @param {*} onSuccess a transform that is called on the payload of a successful action
 * @param {*} initialState an optional initialState for the reducer
 */
export function generateAsyncReducer(resourceName, onSuccess, initialState) {
  const actions = generateAsyncActionTypes(resourceName);
  const combinedInitialState = Object.assign(
    {
      isLoading: false,
      hasErrored: false,
      hasSucceeded: false,
    },
    initialState,
  );
  return function(state = combinedInitialState, action) {
    switch (action.type) {
      case actions.start:
        return Object.assign({}, state, {
          isLoading: true,
          hasErrored: false,
          hasSucceeded: false,
        });
      case actions.success: {
        const result = Object.assign({}, state, {
          isLoading: false,
          hasErrored: false,
          hasSucceeded: true,
          ...onSuccess(action.payload),
        });
        return result;
      }
      case actions.failure: {
        const result = Object.assign({}, state, {
          isLoading: false,
          hasErrored: true,
          hasSucceeded: false,
          error: action.error,
        });
        return result;
      }
      default: {
        return state;
      }
    }
  };
}
