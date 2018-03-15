import { GetCertificationText, CreateCertification } from './api.js';

// Types

export const LOAD_CERTIFICATION_TEXT = 'LOAD_CERTIFICATION_TEXT';
export const LOAD_CERTIFICATION_TEXT_SUCCESS =
  'LOAD_CERTIFICATION_TEXT_SUCCESS';
export const LOAD_CERTIFICATION_TEXT_FAILURE =
  'LOAD_CERTIFICATION_TEXT_FAILURE';

export const CREATE_CERTIFICATION = 'CREATE_CERTIFICATION';
export const CREATE_CERTIFICATION_SUCCESS = 'CREATE_CERTIFICATION_SUCCESS';
export const CREATE_CERTIFICATION_FAILURE = 'CREATE_CERTIFICATION_FAILURE';

// Actions

// loading cert text
export const createLoadCertificationTextRequest = () => ({
  type: LOAD_CERTIFICATION_TEXT,
});

export const createLoadCertificationTextSuccess = certificationText => ({
  type: LOAD_CERTIFICATION_TEXT_SUCCESS,
  certificationText,
});

export const createLoadCertificationTextFailure = error => ({
  type: LOAD_CERTIFICATION_TEXT_FAILURE,
  error,
});

// creating certification
export const createSignedCertificationRequest = () => ({
  type: CREATE_CERTIFICATION,
});

export const createSignedCertificationSuccess = item => ({
  type: CREATE_CERTIFICATION_SUCCESS,
  item,
});

export const createSignedCertificationFailure = error => ({
  type: CREATE_CERTIFICATION_FAILURE,
  error,
});

// Action creator
export function loadCertificationText() {
  return function(dispatch) {
    dispatch(createLoadCertificationTextRequest());
    GetCertificationText()
      .then(spec => dispatch(createLoadCertificationTextSuccess(spec)))
      .catch(error => dispatch(createLoadCertificationTextFailure(error)));
  };
}

export function createSignedCertification(value) {
  return function(dispatch, getState) {
    dispatch(createSignedCertificationRequest());
    CreateCertification(value)
      .then(item => dispatch(createSignedCertificationSuccess(item)))
      .catch(error => dispatch(createSignedCertificationFailure(error)));
  };
}

// Reducer
const initialState = {
  hasSubmitError: false,
  hasSubmitSuccess: false,
  confirmationText: '',
};
export function signedCertificationReducer(state = initialState, action) {
  switch (action.type) {
    case LOAD_CERTIFICATION_TEXT_SUCCESS:
      return Object.assign({}, state, {
        certificationText: action.certificationText,
      });
    case LOAD_CERTIFICATION_TEXT_FAILURE:
      return Object.assign({}, state, {
        certificationText:
          '## Error retrieving legalese. Please reload the page.',
      });
    case CREATE_CERTIFICATION_SUCCESS:
      return Object.assign({}, state, {
        hasSubmitSuccess: true,
        hasSubmitError: false,
        confirmationText: 'Feedback submitted!',
      });
    case CREATE_CERTIFICATION_FAILURE:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
        hasSubmitError: true,
        confirmationText: 'Submission error.',
      });
    default:
      return state;
  }
}

// export default feedbackReducer;
