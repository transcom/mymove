import { filter, map } from 'lodash';
import { denormalize, normalize } from 'normalizr';

import { moveDocuments } from '../schema';
import { ADD_ENTITIES, addEntities } from '../actions';
import { WEIGHT_TICKET_SET_TYPE, MOVE_DOC_TYPE, MOVE_DOC_STATUS } from '../../constants';

import { getClient, checkResponse } from 'shared/Swagger/api';

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

// MoveDocument filter functions
const onlyPending = ({ status }) => ![MOVE_DOC_STATUS.OK, MOVE_DOC_STATUS.EXCLUDE].includes(status);
const onlyOKed = ({ status }) => status === MOVE_DOC_STATUS.OK;
const onlyWeightTickets = ({ move_document_type }) => move_document_type === MOVE_DOC_TYPE.WEIGHT_TICKET_SET;
const onlyProgear = ({ weight_ticket_set_type }) => weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.PRO_GEAR;
const onlyVehicle = ({ weight_ticket_set_type }) => weight_ticket_set_type !== WEIGHT_TICKET_SET_TYPE.PRO_GEAR;

// Common combinations of MoveDocument filters
export function findOKedVehicleWeightTickets(moveDocs) {
  return moveDocs.filter(onlyWeightTickets).filter(onlyVehicle).filter(onlyOKed);
}
export function findOKedProgearWeightTickets(moveDocs) {
  return moveDocs.filter(onlyWeightTickets).filter(onlyProgear).filter(onlyOKed);
}
export function findPendingWeightTickets(moveDocs) {
  return moveDocs.filter(onlyWeightTickets).filter(onlyPending);
}

// Actions
export const getMoveDocumentsForMove = (moveId) => {
  return async function (dispatch, getState, { schema }) {
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

export function createMoveDocument({ moveId, personallyProcuredMoveId, uploadIds, title, moveDocumentType, notes }) {
  return async function (dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.move_docs.createGenericMoveDocument({
      moveId,
      createGenericMoveDocumentPayload: {
        personally_procured_move_id: personallyProcuredMoveId,
        upload_ids: uploadIds,
        title,
        move_document_type: moveDocumentType,
        notes,
      },
    });
    checkResponse(response, 'failed to create move document due to server error');
    const data = normalize(response.body, schema.moveDocument);
    dispatch(addEntities(data.entities));
    return response;
  };
}

export const updateMoveDocument = (moveId, moveDocumentId, payload) => {
  return async function (dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.move_docs.updateMoveDocument({
      moveId,
      moveDocumentId,
      updateMoveDocument: payload,
    });
    checkResponse(response, 'failed to update move document due to server error');
    const data = normalize(response.body, schema.moveDocument);
    dispatch(addEntities(data.entities));
    return response;
  };
};

// Selectors
export const selectMoveDocument = (state, id) => {
  if (!id) {
    return {};
  }
  return denormalize([id], moveDocuments, state.entities)[0];
};

export const selectAllDocumentsForMove = (state, id) => {
  const moveDocs = filter(state.entities.moveDocuments, (doc) => {
    return doc.move_id === id;
  });
  return denormalize(map(moveDocs, 'id'), moveDocuments, state.entities);
};
