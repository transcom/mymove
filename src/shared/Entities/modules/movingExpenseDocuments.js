import { filter, map } from 'lodash';
import { movingExpenseDocuments } from '../schema';
import { ADD_ENTITIES, addEntities } from '../actions';
import { denormalize, normalize } from 'normalizr';

import { getClient, checkResponse } from 'shared/api';

export const STATE_KEY = 'movingExpenseDocuments';

// Reducer
export default function reducer(state = {}, action) {
  switch (action.type) {
    case ADD_ENTITIES:
      return {
        ...state,
        ...action.payload.movingExpenseDocuments,
      };

    default:
      return state;
  }
}

// Actions
export const getMovingExpenseDocumentsForMove = moveId => {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.move_docs.indexMovingExpenseDocuments({
      moveId,
    });
    checkResponse(
      response,
      'failed to get moving expense documents due to server error',
    );

    const data = normalize(response.body, schema.movingExpenseDocuments);
    dispatch(addEntities(data.entities));
    return response;
  };
};

export function createMovingExpenseDocument(
  moveId,
  personallyProcuredMoveId,
  uploadIds,
  title,
  movingExpenseType,
  moveDocumentType,
  reimbursement,
  notes,
) {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.move_docs.createMovingExpenseDocument({
      moveId,
      createMovingExpenseDocumentPayload: {
        personally_procured_move_id: personallyProcuredMoveId,
        upload_ids: uploadIds,
        title: title,
        moving_expense_type: movingExpenseType,
        move_document_type: moveDocumentType,
        reimbursement: reimbursement,
        notes: notes,
      },
    });
    checkResponse(
      response,
      'failed to create moving expense document due to server error',
    );
    const data = normalize(response.body, schema.movingExpenseDocuments);
    dispatch(addEntities(data.entities));
    return response;
  };
}

export const updateMovingExpenseDocument = (
  moveId,
  movingExpenseDocumentId,
  payload,
) => {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.move_docs.updateMovingExpenseDocument({
      moveId,
      movingExpenseDocumentId,
      updateMovingExpenseDocument: payload,
    });
    checkResponse(
      response,
      'failed to update movinge expense document due to server error',
    );
    const data = normalize(response.body, schema.movingExpenseDocuments);
    dispatch(addEntities(data.entities));
    return response;
  };
};

// Selectors
export const selectMovingExpenseDocument = (state, id) => {
  return denormalize([id], movingExpenseDocuments, state.entities)[0];
};

export const selectAllMovingExpenseDocumentsForMove = (state, id) => {
  const movingExpenseDocs = filter(
    state.entities.movingExpenseDocuments,
    doc => {
      return doc.move_id === id;
    },
  );
  return denormalize(
    map(movingExpenseDocs, 'id'),
    movingExpenseDocuments,
    state.entities,
  );
};
