import { each, clone, omit, mapValues, last, startsWith } from 'lodash';

import { ADD_ENTITIES } from 'shared/Entities/actions';

// merge new entities into existing entities
function mergeEntities(entities, newEntities) {
  // shallow clone to mutate
  let result = clone(entities);
  each(newEntities, function (_value, key) {
    /* eslint-disable security/detect-object-injection */
    result[key] = {
      ...result[key],
      ...newEntities[key],
    };
    /* eslint-enable security/detect-object-injection */
  });

  return result;
}

// deletes all items from entities with matching key, id in deleteEntities
function deleteEntities(entities, deleteEntities) {
  return mapValues(entities, function (value, key) {
    //RA Summary: eslint - security/detect-object-injection
    //RA: Using square bracket notation with user input can lead to exploitation
    //RA: Uses object square bracket notation
    //RA: Valuable for state management cleanup
    //RA: The threat actor (web application user) already controls the execution environment (web browser)
    //RA Developer Status: Mitigated
    //RA Validator Status: Mitigated
    //RA Modified Severity: N/A
    // eslint-disable-next-line security/detect-object-injection
    const idsToDelete = Object.keys(deleteEntities[key] || {});
    return omit(value, idsToDelete);
  });
}

const initialState = {
  backupContacts: {},
  documents: {},
  invoices: {},
  moveDocuments: {},
  moves: {},
  orders: {},
  personallyProcuredMoves: {},
  mtoShipments: {},
  reimbursements: {},
  signedCertifications: {},
  tariff400ngItems: {},
  transportationServiceProviders: {},
};

// Actions of either of these types will be merged into the store:
//   @@swagger/tag.operationId/SUCCESS
//   ADD_ENTITIES
export function entitiesReducer(state = initialState, action) {
  if (action.type === 'UPDATE_MTO_SHIPMENTS_ENTITIY') {
    return {
      ...state,
      mtoShipments: action.entities['mtoShipments'],
    };
  }
  if (action.type === ADD_ENTITIES) {
    return mergeEntities(state, action.entities);
  }
  if (startsWith(action.type, '@@swagger')) {
    const parts = action.type.split('/');
    if (last(parts) === 'SUCCESS') {
      if (action.method === 'delete') {
        return deleteEntities(state, action.entities);
      }
      return mergeEntities(state, action.entities);
    }
  }
  return state;
}
