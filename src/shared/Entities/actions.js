export const ADD_ENTITIES = 'ADD_ENTITIES';
export const SET_OKTA_USER = 'SET_OKTA_USER';
export const SET_ADMIN_USER = 'SET_ADMIN_USER';

export const addEntities = (entities) => ({
  type: ADD_ENTITIES,
  entities,
});

export const updateMTOShipmentsEntity = (entities) => ({
  type: 'UPDATE_MTO_SHIPMENTS_ENTITIY',
  entities,
});

export const setOktaUser = (oktaUser) => ({
  type: SET_OKTA_USER,
  oktaUser,
});

export const setAdminUser = (adminUser) => ({
  type: SET_ADMIN_USER,
  adminUser,
});
