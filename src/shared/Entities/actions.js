export { getMoveDocumentsForMove } from './modules/moveDocuments';

export const ADD_ENTITIES = 'ADD_ENTITIES';
export const addEntities = entities => ({
  type: ADD_ENTITIES,
  payload: entities,
});
