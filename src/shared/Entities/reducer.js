import { last, startsWith } from 'lodash';

import { ADD_ENTITIES } from 'shared/Entities/actions';

// merge new entities into existing entities
function mergeEntities(entities, newEntities) {
  Object.keys(newEntities).forEach(function(key) {
    entities[key] = {
      ...entities[key],
      ...newEntities[key],
    };
  });
  return entities;
}

const initialState = {
  shipments: {},
};

// Actions of either of these types will be merged into the store:
//   @@swagger/tag.operationId/SUCCESS
//   ADD_ENTITIES
export function entitiesReducer(state = initialState, action) {
  if (action.type === ADD_ENTITIES) {
    return mergeEntities(state, action.entities);
  }
  if (startsWith(action.type, '@@swagger')) {
    const parts = action.type.split('/');
    if (last(parts) === 'SUCCESS') {
      return mergeEntities(state, action.entities);
    }
  }
  return state;
}
