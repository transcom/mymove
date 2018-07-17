import { cloneDeep } from 'lodash';
import * as ReduxHelpers from 'shared/ReduxHelpers';

import {
  IndexMoveDocuments,
  CreateMoveDocument,
  UpdateMoveDocument,
} from './api.js';
import { upsert } from 'shared/utils';

// Types
const indexMoveDocumentsType = 'INDEX_MOVE_DOCUMENTS';
const createMoveDocumentType = 'CREATE_MOVE_DOCUMENT';
const updateMoveDocumentType = 'UPDATE_MOVE_DOCUMENT';

// Action types
export const INDEX_MOVE_DOCUMENTS = ReduxHelpers.generateAsyncActionTypes(
  indexMoveDocumentsType,
);
export const CREATE_MOVE_DOCUMENT = ReduxHelpers.generateAsyncActionTypes(
  createMoveDocumentType,
);
export const UPDATE_MOVE_DOCUMENT = ReduxHelpers.generateAsyncActionTypes(
  updateMoveDocumentType,
);

// Action creators
export const indexMoveDocuments = ReduxHelpers.generateAsyncActionCreator(
  indexMoveDocumentsType,
  IndexMoveDocuments,
);

export const createMoveDocument = ReduxHelpers.generateAsyncActionCreator(
  createMoveDocumentType,
  CreateMoveDocument,
);
export const updateMoveDocument = ReduxHelpers.generateAsyncActionCreator(
  updateMoveDocumentType,
  UpdateMoveDocument,
);

export function updateMoveDocumentInfo(moveId, moveDocId, moveDocument) {
  const actions = ReduxHelpers.generateAsyncActions(updateMoveDocumentType);
  return async function(dispatch, getState) {
    dispatch(actions.start());
    try {
      await dispatch(updateMoveDocument(moveId, moveDocId, moveDocument));
    } catch (ex) {
      return dispatch(actions.error(ex));
    }
  };
}

// Reducer
const initialState = {
  moveDocuments: [],
  indexMoveDocumentsSuccess: false,
};
const upsertMoveDocument = (moveDocument, state) => {
  const newState = cloneDeep(state);
  upsert(newState.moveDocuments, moveDocument);
  return newState;
};

export function documentsReducer(state = initialState, action) {
  switch (action.type) {
    case INDEX_MOVE_DOCUMENTS.start:
      return Object.assign({}, state, {
        indexMoveDocumentsSuccess: false,
      });
    case INDEX_MOVE_DOCUMENTS.success:
      return Object.assign({}, state, {
        moveDocuments: action.payload,
        indexMoveDocumentsSuccess: true,
        indexMoveDocumentsError: false,
      });
    case INDEX_MOVE_DOCUMENTS.failure:
      return Object.assign({}, state, {
        indexMoveDocumentsSuccess: false,
        indexMoveDocumentsError: true,
        error: action.error,
      });
    case CREATE_MOVE_DOCUMENT.success:
      return {
        ...upsertMoveDocument(action.payload, state),
        moveDocumentCreateError: false,
        updatedMoveDocument: action.payload,
      };
    case CREATE_MOVE_DOCUMENT.failure:
      return Object.assign({}, state, {
        moveDocumentCreateError: true,
        error: action.error.message,
      });
    case UPDATE_MOVE_DOCUMENT.success:
      return {
        ...upsertMoveDocument(action.payload, state),
        moveDocumentUpdateError: false,
        updatedMoveDocument: action.payload,
      };
    case UPDATE_MOVE_DOCUMENT.failure:
      return Object.assign({}, state, {
        moveDocumentUpdateError: true,
        error: action.error.message,
      });
    default:
      return state;
  }
}
