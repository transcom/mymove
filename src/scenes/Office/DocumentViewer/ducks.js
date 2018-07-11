import { IndexMoveDocuments } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
export const indexMoveDocumentsType = 'INDEX_MOVE_DOCUMENTS';

// Actions
export const INDEX_MOVE_DOCUMENTS = ReduxHelpers.generateAsyncActionTypes(
  indexMoveDocumentsType,
);

export const indexMoveDocuments = ReduxHelpers.generateAsyncActionCreator(
  indexMoveDocumentsType,
  IndexMoveDocuments,
);

// Reducer
const initialState = {
  moveDocuments: null,
  indexMoveDocumentsSuccess: false,
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
    default:
      return state;
  }
}
