import { CreateDocument } from './api.js';

// Types
export const CREATE_DOCUMENT = 'CREATE_DOCUMENT';
export const CREATE_DOCUMENT_SUCCESS = 'CREATE_DOCUMENT_SUCCESS';
export const CREATE_DOCUMENT_FAILURE = 'CREATE_DOCUMENT_FAILURE';

// Actions
// creating document
export const createDocumentRequest = () => ({
  type: CREATE_DOCUMENT,
});

export const createDocumentSuccess = item => ({
  type: CREATE_DOCUMENT_SUCCESS,
  item,
});

export const createDocumentFailure = error => ({
  type: CREATE_DOCUMENT_FAILURE,
  error,
});

// Action creator
export function createDocument(fileUpload, moveID) {
  return function(dispatch, getState) {
    dispatch(createDocumentRequest());
    CreateDocument(fileUpload, moveID)
      .then(item => dispatch(createDocumentSuccess(item)))
      .catch(error => dispatch(createDocumentFailure(error)));
  };
}

// Reducer
const initialState = {
  hasErrored: false,
  hasSucceeded: false,
  confirmationText: '',
  upload: null,
};

export function documentReducer(state = initialState, action) {
  switch (action.type) {
    case CREATE_DOCUMENT_SUCCESS:
      return Object.assign({}, state, {
        hasSucceeded: true,
        hasErrored: false,
        upload: action.item,
      });
    case CREATE_DOCUMENT_FAILURE:
      return Object.assign({}, state, {
        hasSucceeded: false,
        hasErrored: true,
        upload: {},
      });
    default:
      return state;
  }
}

export default documentReducer;
