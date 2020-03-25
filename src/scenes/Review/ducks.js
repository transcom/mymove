import { ValidateEntitlement } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
const editBeginType = 'EDIT_BEGIN';
const entitlementChangeBeginType = 'ENTITLEMENT_CHANGE_BEGIN';
const validateEntitlement = 'VALIDATE_ENTITLEMENT';
export const VALIDATE_ENTITLEMENT = ReduxHelpers.generateAsyncActionTypes(validateEntitlement);

const editSuccessfulType = 'EDIT_SUCCESS';
const entitlementChangeType = 'ENTITLEMENT_CHANGED';

// Actions
export function editBegin() {
  return function(dispatch, getState) {
    dispatch({ type: editBeginType });
  };
}

export function editSuccessful() {
  return function(dispatch, getState) {
    dispatch({ type: editSuccessfulType });
  };
}

export function entitlementChangeBegin() {
  return function(dispatch, getState) {
    dispatch({ type: entitlementChangeBeginType });
  };
}

export function entitlementChanged() {
  return function(dispatch, getState) {
    dispatch({ type: entitlementChangeType });
  };
}

export const checkEntitlement = ReduxHelpers.generateAsyncActionCreator(validateEntitlement, ValidateEntitlement);

// Reducer
const initialState = {
  editSuccess: null,
  entitlementChange: null,
  error: null,
};
export function reviewReducer(state = initialState, action) {
  switch (action.type) {
    case editBeginType:
      return Object.assign({}, state, {
        editSuccess: false,
      });
    case editSuccessfulType:
      return Object.assign({}, state, {
        editSuccess: true,
      });
    case entitlementChangeBeginType:
      return Object.assign({}, state, {
        entitlementChange: false,
      });
    case entitlementChangeType:
      return Object.assign({}, state, {
        entitlementChange: true,
      });
    case VALIDATE_ENTITLEMENT.success:
      return Object.assign({}, state, {
        error: null,
      });
    case VALIDATE_ENTITLEMENT.failure:
      return Object.assign({}, state, {
        error: action.error,
      });
    default:
      return state;
  }
}
