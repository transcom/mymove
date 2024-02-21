import { includes, get } from 'lodash';
import { normalize } from 'normalizr';

import { addEntities } from '../actions';

import { getClient, checkResponse } from 'shared/Swagger/api';

const expenseTypes = ['EXPENSE', 'STORAGE_EXPENSE'];

export function isMovingExpenseDocument(moveDocument) {
  const type = get(moveDocument, 'move_document_type', '');
  return includes(expenseTypes, type);
}

export function createMovingExpenseDocument({
  moveId,
  personallyProcuredMoveId,
  uploadIds,
  title,
  movingExpenseType,
  moveDocumentType,
  requestedAmountCents,
  paymentMethod,
  notes,
  missingReceipt,
  storage_start_date,
  storage_end_date,
}) {
  return async function (dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.move_docs.createMovingExpenseDocument({
      moveId,
      createMovingExpenseDocumentPayload: {
        personally_procured_move_id: personallyProcuredMoveId,
        upload_ids: uploadIds,
        title,
        moving_expense_type: movingExpenseType,
        move_document_type: moveDocumentType,
        requested_amount_cents: requestedAmountCents,
        payment_method: paymentMethod,
        notes,
        receipt_missing: missingReceipt,
        storage_start_date,
        storage_end_date,
      },
    });
    checkResponse(response, 'failed to create moving expense document due to server error');
    const data = normalize(response.body, schema.moveDocument);
    dispatch(addEntities(data.entities));
    return response;
  };
}
