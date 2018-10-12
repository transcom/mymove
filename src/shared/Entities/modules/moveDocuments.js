import { filter, map } from 'lodash';
import { moveDocuments } from '../schema';
import { ADD_ENTITIES, addEntities } from '../actions';
import { denormalize, normalize } from 'normalizr';

import { getClient, getPublicClient, checkResponse } from 'shared/Swagger/api';

export const STATE_KEY = 'moveDocuments';

// Reducer
export default function reducer(state = {}, action) {
  switch (action.type) {
    case ADD_ENTITIES:
      return {
        ...state,
        ...action.payload.moveDocuments,
      };

    default:
      return state;
  }
}

// Actions
export const getMoveDocumentsForMove = moveId => {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.move_docs.indexMoveDocuments({
      moveId,
    });
    checkResponse(response, 'failed to get move documents due to server error');

    const data = normalize(response.body, schema.moveDocuments);
    dispatch(addEntities(data.entities));
    return response;
  };
};

export function createMoveDocument(
  moveId,
  personallyProcuredMoveId,
  uploadIds,
  title,
  moveDocumentType,
  notes,
) {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.move_docs.createGenericMoveDocument({
      moveId,
      createGenericMoveDocumentPayload: {
        personally_procured_move_id: personallyProcuredMoveId,
        upload_ids: uploadIds,
        title: title,
        move_document_type: moveDocumentType,
        notes: notes,
      },
    });
    checkResponse(
      response,
      'failed to create move document due to server error',
    );
    const data = normalize(response.body, schema.moveDocument);
    dispatch(addEntities(data.entities));
    return response;
  };
}

export function createShipmentDocument(
  shipmentId,
  createGenericMoveDocumentPayload,
) {
  return async function(dispatch, getState, { schema }) {
    const client = await getPublicClient();
    const response = await client.apis.move_docs.createGenericMoveDocument({
      shipmentId,
      createGenericMoveDocumentPayload,
    });
    checkResponse(
      response,
      'failed to create move document due to server error',
    );
    const data = normalize(response.body, schema.moveDocument);
    dispatch(addEntities(data.entities));
    return response;
  };
}

export const updateMoveDocument = (moveId, moveDocumentId, payload) => {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.move_docs.updateMoveDocument({
      moveId,
      moveDocumentId,
      updateMoveDocument: payload,
    });
    checkResponse(
      response,
      'failed to update move document due to server error',
    );
    const data = normalize(response.body, schema.moveDocument);
    dispatch(addEntities(data.entities));
    return response;
  };
};

// Selectors
export const selectMoveDocument = (state, id) => {
  return denormalize([id], moveDocuments, state.entities)[0];
};

export const selectAllDocumentsForMove = (state, id) => {
  const moveDocs = filter(state.entities.moveDocuments, doc => {
    return doc.move_id === id;
  });
  return denormalize(map(moveDocs, 'id'), moveDocuments, state.entities);
};
