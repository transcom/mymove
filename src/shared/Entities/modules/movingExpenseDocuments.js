import { includes, get, filter, map } from 'lodash';
import { moveDocuments } from '../schema';
import { addEntities } from '../actions';
import { denormalize, normalize } from 'normalizr';

import { getClient, checkResponse } from 'shared/Swagger/api';

const expenseTypes = ['EXPENSE', 'STORAGE_EXPENSE'];

export function isMovingExpenseDocument(moveDocument) {
  const type = get(moveDocument, 'move_document_type', '');
  return includes(expenseTypes, type);
}

export function createMovingExpenseDocument(
  moveId,
  personallyProcuredMoveId,
  uploadIds,
  title,
  movingExpenseType,
  moveDocumentType,
  requestedAmountCents,
  paymentMethod,
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
        requested_amount_cents: requestedAmountCents,
        payment_method: paymentMethod,
        notes: notes,
      },
    });
    checkResponse(
      response,
      'failed to create moving expense document due to server error',
    );
    const data = normalize(response.body, schema.moveDocument);
    dispatch(addEntities(data.entities));
    return response;
  };
}

export const selectAllMovingExpenseDocumentsForMove = (state, id) => {
  const movingExpenseDocs = filter(state.entities.moveDocuments, doc => {
    return doc.move_id === id && doc.move_document_type === 'EXPENSE';
  });
  return denormalize(
    map(movingExpenseDocs, 'id'),
    moveDocuments,
    state.entities,
  );
};
