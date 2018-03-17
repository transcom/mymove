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
