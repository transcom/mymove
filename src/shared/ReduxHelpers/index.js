function actionName(resource, actionType) {
  return `${resource}_${actionType.toUpperCase()}`;
}

function generateActionTypes(resourceName, actionTypes) {
  resourceName = (resourceName || '').trim();
  if (!resourceName) throw new Error('No resource name provided');
  let actions = {};
  actionTypes.forEach(actionType => {
    actions[actionType] = actionName(resourceName, actionType);
  });
  return actions;
}

export function generateAsyncActionTypes(resourceName) {
  return generateActionTypes(resourceName, ['start', 'success', 'failure']);
}

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

export function generateAsyncActionCreator(resourceName, asyncAction) {
  const actions = generateAsyncActions(resourceName);
  return function actionCreator(...args) {
    return function(dispatch) {
      dispatch(actions.start());
      asyncAction(...args)
        .then(item => dispatch(actions.success(item)))
        .catch(error => dispatch(actions.error(error)));
    };
  };
}

export function generateAsyncReducer(resourceName, onSuccess) {
  const actions = generateAsyncActionTypes(resourceName);
  const initialState = { isLoading: false, hasErrored: false };
  return function(state = initialState, action) {
    switch (action.type) {
      case actions.start:
        return Object.assign({}, state, {
          isLoading: true,
          hasErrored: false,
        });
      case actions.success: {
        const result = Object.assign({}, state, {
          isLoading: false,
          hasErrored: false,
          ...onSuccess(action.payload),
        });
        // result[payloadKey] = onSuccess(action.payload);
        return result;
      }
      case actions.failure: {
        const result = Object.assign({}, state, {
          isLoading: false,
          hasErrored: true,
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
