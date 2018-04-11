import { GetServiceMember, UpdateServiceMember } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
export const SET_PENDING_SM_NAME_DATA = 'SET_PENDING_SM_NAME_DATA';
export const GET_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes(
  'GET_SERVICE_MEMBER',
);
export const UPDATE_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes(
  'UPDATE_SERVICE_MEMBER',
);

// Action creation
export function setPendingSMNameData(value) {
  return { type: SET_PENDING_SM_NAME_DATA, payload: value };
}

export function updateServiceMember(serviceMember) {
  const action = ReduxHelpers.generateAsyncActions('UPDATE_SERVICE_MEMBER');
  return function(dispatch, getState) {
    dispatch(action.start());
    const state = getState();
    const currentServiceMember = state.serviceMember.currentServiceMember;
    if (currentServiceMember) {
      UpdateServiceMember(currentServiceMember.id, serviceMember)
        .then(item =>
          dispatch(
            action.success(Object.assign({}, currentServiceMember, item)),
          ),
        )
        .catch(error => dispatch(action.error(error)));
    }
  };
}

export function loadServiceMember(serviceMemberId) {
  const action = ReduxHelpers.generateAsyncActions('GET_SERVICE_MEMBER');
  return function(dispatch, getState) {
    dispatch(action.start);
    const state = getState();
    const currentServiceMember = state.serviceMember.currentServiceMember;
    if (!currentServiceMember) {
      GetServiceMember(serviceMemberId)
        .then(item => dispatch(action.success(item)))
        .catch(error => dispatch(action.error(error)));
    }
  };
}

// Reducer
const initialState = {
  pendingSMNameData: null,
  currentServiceMember: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
};
export function serviceMemberReducer(state = initialState, action) {
  switch (action.type) {
    case SET_PENDING_SM_NAME_DATA:
      return Object.assign({}, state, {
        pendingSMNameData: action.payload,
      });
    case UPDATE_SERVICE_MEMBER.start:
      return Object.assign({}, state, {
        hasSubmitSuccess: true,
      });
    case UPDATE_SERVICE_MEMBER.success:
      return Object.assign({}, state, {
        currentServiceMember: action.payload,
        pendingSMNameData: null,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case UPDATE_SERVICE_MEMBER.failure:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case GET_SERVICE_MEMBER.start:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
      });
    case GET_SERVICE_MEMBER.success:
      return Object.assign({}, state, {
        currentServiceMember: action.payload,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case GET_SERVICE_MEMBER.failure:
      return Object.assign({}, state, {
        currentServiceMember: null,
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
