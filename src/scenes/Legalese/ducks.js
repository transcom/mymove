import {
  GetCertifications,
  GetCertificationText,
  CreateCertification,
} from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';
import { get } from 'lodash';

// Actions

export const CREATE_SIGNED_CERT = ReduxHelpers.generateAsyncActionTypes(
  'CREATE_SIGNED_CERT',
);

export const GET_LATEST_CERT = ReduxHelpers.generateAsyncActionTypes(
  'GET_LATEST_CERT',
);

export const GET_CERT_TEXT = ReduxHelpers.generateAsyncActionTypes(
  'GET_CERT_TEXT',
);

// Action creator
export const loadCertificationText = ReduxHelpers.generateAsyncActionCreator(
  'GET_CERT_TEXT',
  GetCertificationText,
);

export const createSignedCertification = ReduxHelpers.generateAsyncActionCreator(
  'CREATE_SIGNED_CERT',
  CreateCertification,
);

export function loadLatestCertification(moveId) {
  const action = ReduxHelpers.generateAsyncActions('GET_LATEST_CERT');
  return function(dispatch, getState) {
    dispatch(action.start);
    const state = getState();
    const latestSignedCertification =
      state.signedCertification.latestSignedCertification;
    if (!latestSignedCertification) {
      GetCertifications(moveId, 1)
        .then(item => dispatch(action.success(item)))
        .catch(error => dispatch(action.error(error)));
    }
  };
}

// Reducer
const initialState = {
  hasSubmitError: false,
  hasSubmitSuccess: false,
  confirmationText: '',
  getCertificationSuccess: false,
  getCertificationError: false,
  latestSignedCertification: null,
  certificationText: null,
  error: null,
};
export function signedCertificationReducer(state = initialState, action) {
  switch (action.type) {
    case GET_CERT_TEXT.success:
      return Object.assign({}, state, {
        certificationText: action.payload,
      });
    case GET_CERT_TEXT.failure:
      return Object.assign({}, state, {
        certificationText:
          '## Error retrieving legalese. Please reload the page.',
        error: action.error,
      });
    case CREATE_SIGNED_CERT.success:
      return Object.assign({}, state, {
        hasSubmitSuccess: true,
        hasSubmitError: false,
        confirmationText: 'Feedback submitted!',
      });
    case CREATE_SIGNED_CERT.failure:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
        hasSubmitError: true,
        confirmationText: 'Submission error.',
      });
    case GET_LATEST_CERT.start:
      return Object.assign({}, state, {
        getCertificationSuccess: false,
      });
    case GET_LATEST_CERT.success:
      return Object.assign({}, state, {
        getCertificationSuccess: true,
        getCertificationError: false,
        latestSignedCertification: get(action, 'payload.0', null),
        certificationText: get(action, 'payload.0.certification_text', null),
      });
    case GET_LATEST_CERT.failure:
      return Object.assign({}, state, {
        getCertificationSuccess: false,
        getCertificationError: true,
        latestSignedCertification: null,
        certificationText: null,
        error: get(action, 'error', null),
      });
    default:
      return state;
  }
}

// export default feedbackReducer;
