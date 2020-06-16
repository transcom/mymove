import { documents } from '../schema';
import { ADD_ENTITIES } from '../actions';
import { denormalize } from 'normalizr';
import { getClient } from 'shared/Swagger/api';
import { swaggerRequest } from 'shared/Swagger/request';

export const STATE_KEY = 'documents';
export const createUploadLabel = 'documents.createUpload';

export default function reducer(state = {}, action) {
  switch (action.type) {
    case ADD_ENTITIES:
      return {
        ...state,
        ...action.payload.documents,
      };

    default:
      return state;
  }
}

// Actions
export function createUpload(fileUpload, documentId, isPublic = false, label = createUploadLabel) {
  const swaggerTag = 'uploads.createUpload';
  return swaggerRequest(
    getClient,
    swaggerTag,
    {
      documentId,
      file: fileUpload,
    },
    { label },
  );
}

// Selectors
export const selectDocument = (state, id) => {
  if (!id) {
    return {};
  }
  return denormalize([id], documents, state.entities)[0];
};
