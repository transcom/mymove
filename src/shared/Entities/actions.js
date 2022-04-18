export const ADD_ENTITIES = 'ADD_ENTITIES';
export const addEntities = (entities) => ({
  type: ADD_ENTITIES,
  entities,
});

export const updateMTOShipmentsEntity = (entities) => ({
  type: 'UPDATE_MTO_SHIPMENTS_ENTITIY',
  entities,
});
