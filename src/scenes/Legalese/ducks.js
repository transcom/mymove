import * as ReduxHelpers from 'shared/ReduxHelpers';
import moment from 'moment';

const signAndSubmitForApprovalType = 'SIGN_AND_SUBMIT_FOR_APPROVAL';

// Actions
export const CREATE_SIGNED_CERT = ReduxHelpers.generateAsyncActionTypes('CREATE_SIGNED_CERT');

export const GET_LATEST_CERT = ReduxHelpers.generateAsyncActionTypes('GET_LATEST_CERT');

export const GET_CERT_TEXT = ReduxHelpers.generateAsyncActionTypes('GET_CERT_TEXT');

// Action creator
const SIGN_AND_SUBMIT_FOR_APPROVAL = ReduxHelpers.generateAsyncActionTypes(signAndSubmitForApprovalType);

export function dateToTimestamp(dt) {
  return moment(dt).format();
}

// Reducer
const initialState = {
  hasSubmitError: false,
  hasSubmitSuccess: false,
  confirmationText: '',
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
        certificationText: '## Error retrieving legalese. Please reload the page.',
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
    case SIGN_AND_SUBMIT_FOR_APPROVAL.success:
      return { ...state, moveSubmitSuccess: true };
    case SIGN_AND_SUBMIT_FOR_APPROVAL.failure:
      return { ...state, error: action.error };
    default:
      return state;
  }
}

// export default feedbackReducer;
