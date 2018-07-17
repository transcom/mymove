import { filter, map } from 'lodash';
import { moveDocuments } from '../schema';
import { ADD_ENTITIES, addEntities } from '../actions';
import { denormalize, normalize } from 'normalizr';

import { getClient, checkResponse } from 'shared/api';

export const STATE_KEY = 'moveDocuments';

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

export const getMoveDocumentsForMove = moveId => {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const response = await client.apis.moves.indexMoveDocuments({
      moveId,
    });
    checkResponse(response, 'failed to get move documents due to server error');

    const data = normalize(response.body, schema.moveDocuments);
    dispatch(addEntities(data.entities));
    return response;
  };
};

export const selectMoveDocument = (state, id) => {
  return denormalize([id], moveDocuments, state.entities)[0];
};

export const selectAllDocumentsForMove = (state, id) => {
  debugger;
  const moveDocs = filter(state.entities.moveDocuments, doc => {
    return doc.move_id === id;
  });
  return denormalize(map(moveDocs, 'id'), moveDocuments, state.entities);
};
