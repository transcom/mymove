import { cloneDeep } from 'lodash';
import * as ReduxHelpers from 'shared/ReduxHelpers';

import { IndexMoveDocuments, CreateMoveDocument } from './api.js';
import { upsert } from 'shared/utils';

// Types
export const indexMoveDocumentsType = 'INDEX_MOVE_DOCUMENTS';
const createMoveDocumentType = 'CREATE_MOVE_DOCUMENT';

// Action types
export const INDEX_MOVE_DOCUMENTS = ReduxHelpers.generateAsyncActionTypes(
  indexMoveDocumentsType,
);
const CREATE_MOVE_DOCUMENT = ReduxHelpers.generateAsyncActionTypes(
  createMoveDocumentType,
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
        createdMoveDocument: action.payload,
      };
    case CREATE_MOVE_DOCUMENT.failure:
      return Object.assign({}, state, {
        moveDocumentCreateError: true,
        error: action.error.message,
      });
    default:
      return state;
  }
}
