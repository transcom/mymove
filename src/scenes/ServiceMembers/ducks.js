import {
  GetServiceMember,
  UpdateServiceMember,
  CreateServiceMember,
  IndexBackupContactsAPI,
  CreateBackupContactAPI,
  UpdateBackupContactAPI,
} from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
export const GET_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes(
  'GET_SERVICE_MEMBER',
);
export const UPDATE_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes(
  'UPDATE_SERVICE_MEMBER',
);

const createServiceMemberType = 'CREATE_SERVICE_MEMBER';
const CREATE_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes(
  createServiceMemberType,
);

export const createServiceMember = ReduxHelpers.generateAsyncActionCreator(
  createServiceMemberType,
  CreateServiceMember,
);

const createBackupContactType = 'CREATE_BACKUP_CONTACT';
const indexBackupContactsType = 'INDEX_BACKUP_CONTACTS';
const updateBackupContactType = 'UPDATE_BACKUP_CONTACT';

const CREATE_BACKUP_CONTACT = ReduxHelpers.generateAsyncActionTypes(
  createBackupContactType,
);

const INDEX_BACKUP_CONTACTS = ReduxHelpers.generateAsyncActionTypes(
  indexBackupContactsType,
);

const UPDATE_BACKUP_CONTACT = ReduxHelpers.generateAsyncActionTypes(
  updateBackupContactType,
);

export const createBackupContact = ReduxHelpers.generateAsyncActionCreator(
  createBackupContactType,
  CreateBackupContactAPI,
);

export const indexBackupContacts = ReduxHelpers.generateAsyncActionCreator(
  indexBackupContactsType,
  IndexBackupContactsAPI,
);

export const updateBackupContact = ReduxHelpers.generateAsyncActionCreator(
  updateBackupContactType,
  UpdateBackupContactAPI,
);

// Action creation
export function updateServiceMember(serviceMember) {
  const action = ReduxHelpers.generateAsyncActions('UPDATE_SERVICE_MEMBER');
  return function(dispatch, getState) {
    dispatch(action.start());
    const state = getState();
    const currentServiceMember = state.serviceMember.currentServiceMember;
    if (currentServiceMember) {
      UpdateServiceMember(currentServiceMember.id, serviceMember)
        .then(item =>
          dispatch(
            action.success(Object.assign({}, currentServiceMember, item)),
          ),
        )
        .catch(error => dispatch(action.error(error)));
    }
  };
}

export function loadServiceMember(serviceMemberId) {
  const action = ReduxHelpers.generateAsyncActions('GET_SERVICE_MEMBER');
  return function(dispatch, getState) {
    dispatch(action.start);
    const state = getState();
    const currentServiceMember = state.serviceMember.currentServiceMember;
    if (!currentServiceMember) {
      GetServiceMember(serviceMemberId)
        .then(item => dispatch(action.success(item)))
        .catch(error => dispatch(action.error(error)));
    }
  };
}

// Reducer
const initialState = {
  currentServiceMember: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
};
export function serviceMemberReducer(state = initialState, action) {
  console.log('REDUCIN', action);
  switch (action.type) {
    case CREATE_SERVICE_MEMBER.start:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
      });
    case CREATE_SERVICE_MEMBER.success:
      return Object.assign({}, state, {
        currentServiceMember: action.payload,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case CREATE_SERVICE_MEMBER.failure:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case UPDATE_SERVICE_MEMBER.start:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
      });
    case UPDATE_SERVICE_MEMBER.success:
      return Object.assign({}, state, {
        currentServiceMember: action.payload,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case UPDATE_SERVICE_MEMBER.failure:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case GET_SERVICE_MEMBER.start:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
      });
    case GET_SERVICE_MEMBER.success:
      return Object.assign({}, state, {
        currentServiceMember: action.payload,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case GET_SERVICE_MEMBER.failure:
      return Object.assign({}, state, {
        currentServiceMember: null,
        hasSubmitSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    // Backup Contacts!
    case CREATE_BACKUP_CONTACT.start:
      return Object.assign({}, state, {
        createBackupContactSuccess: false,
      });
    case CREATE_BACKUP_CONTACT.success:
      let newBackupContacts = state.currentBackupContacts || [];
      newBackupContacts.push(action.payload);
      return Object.assign({}, state, {
        currentBackupContacts: newBackupContacts,
        createdBackupContact: action.payload,
        createBackupContactSuccess: true,
        createBackupContactError: false,
      });
    case CREATE_BACKUP_CONTACT.failure:
      return Object.assign({}, state, {
        createBackupContactSuccess: false,
        hasSubmitError: true,
        error: action.error,
      });
    case UPDATE_BACKUP_CONTACT.start:
      return Object.assign({}, state, {
        updateBackupContactSuccess: false,
      });
    case UPDATE_BACKUP_CONTACT.success:
      // replace the updated contact in the list
      newBackupContacts = state.currentBackupContacts;
      const staleIndex = newBackupContacts.findIndex(element => {
        return (element.id = action.payload.id);
      });
      newBackupContacts[staleIndex] = action.payload;
      return Object.assign({}, state, {
        currentServiceMember: action.payload,
        currentBackupContacts: newBackupContacts,
        updateBackupContactSuccess: true,
        updateBackupContactError: false,
      });
    case UPDATE_BACKUP_CONTACT.failure:
      return Object.assign({}, state, {
        updateBackupContactSuccess: false,
        updateBackupContactError: true,
        error: action.error,
      });
    case INDEX_BACKUP_CONTACTS.start:
      return Object.assign({}, state, {
        indexBackupContactsSuccess: false,
      });
    case INDEX_BACKUP_CONTACTS.success:
      return Object.assign({}, state, {
        currentBackupContacts: action.payload,
        indexBackupContactsSuccess: true,
        indexBackupContactsError: false,
      });
    case INDEX_BACKUP_CONTACTS.failure:
      return Object.assign({}, state, {
        currentBackupContacts: null,
        indexBackupContactsSuccess: false,
        indexBackupContactsError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
