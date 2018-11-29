import { GetCertificationText, CreateCertification } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';
import { pick } from 'lodash';
import { SubmitForApproval } from '../Moves/ducks.js';
import { normalize } from 'normalizr';
import { move } from 'shared/Entities/schema';
import { addEntities } from 'shared/Entities/actions';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

const signAndSubmitForApprovalType = 'SIGN_AND_SUBMIT_FOR_APPROVAL';
const signAndSubmitPpmForApprovalType = 'SIGN_AND_SUBMIT_PPM_FOR_APPROVAL';

// Actions

export const CREATE_SIGNED_CERT = ReduxHelpers.generateAsyncActionTypes('CREATE_SIGNED_CERT');

export const GET_LATEST_CERT = ReduxHelpers.generateAsyncActionTypes('GET_LATEST_CERT');

export const GET_CERT_TEXT = ReduxHelpers.generateAsyncActionTypes('GET_CERT_TEXT');

// Action creator
export const loadCertificationText = ReduxHelpers.generateAsyncActionCreator('GET_CERT_TEXT', GetCertificationText);

const createSignedCertification = ReduxHelpers.generateAsyncActionCreator('CREATE_SIGNED_CERT', CreateCertification);

const SIGN_AND_SUBMIT_FOR_APPROVAL = ReduxHelpers.generateAsyncActionTypes(signAndSubmitForApprovalType);

const signAndSubmitForApprovalActions = ReduxHelpers.generateAsyncActions(signAndSubmitForApprovalType);
const signAndSubmitPpmForApprovalActions = ReduxHelpers.generateAsyncActions(signAndSubmitPpmForApprovalType);

export const signAndSubmitForApproval = (moveId, certificationText, signature, dateSigned) => {
  return async function(dispatch, getState) {
    dispatch(signAndSubmitForApprovalActions.start());
    try {
      await dispatch(
        createSignedCertification({
          moveId,
          createSignedCertificationPayload: {
            certification_text: certificationText,
            signature,
            date: dateSigned,
          },
        }),
      );
      const response = await dispatch(SubmitForApproval(moveId));
      const data = normalize(response.payload, move);
      const filtered = pick(data.entities, ['shipments', 'moves']);
      dispatch(addEntities(filtered));
      return dispatch(signAndSubmitForApprovalActions.success());
    } catch (error) {
      return dispatch(signAndSubmitForApprovalActions.error(error));
    }
  };
};

export const signAndSubmitPpm = (moveId, certificationText, signature, dateSigned, ppmId) => {
  return async function(dispatch, getState) {
    dispatch(signAndSubmitPpmForApprovalActions.start());
    try {
      await dispatch(
        createSignedCertification({
          moveId,
          createSignedCertificationPayload: {
            certification_text: certificationText,
            signature,
            date: dateSigned,
          },
        }),
      );
      await dispatch(submitPpm(ppmId));
      return dispatch(signAndSubmitPpmForApprovalActions.success());
    } catch (error) {
      console.log(error);
      return dispatch(signAndSubmitPpmForApprovalActions.error(error));
    }
  };
};

export function submitPpm(personallyProcuredMoveId) {
  return swaggerRequest(
    getClient,
    'ppm.submitPersonallyProcuredMove',
    { personallyProcuredMoveId },
    { label: 'submit_ppm' },
  );
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
