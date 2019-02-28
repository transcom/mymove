import { pick, without, cloneDeep, get, every } from 'lodash';
import { NULL_UUID } from 'shared/constants';

import {
  GetServiceMember,
  UpdateServiceMember,
  CreateServiceMember,
  IndexBackupContactsAPI,
  CreateBackupContactAPI,
  UpdateBackupContactAPI,
} from './api.js';
import { GET_LOGGED_IN_USER } from 'shared/Data/users';
import * as ReduxHelpers from 'shared/ReduxHelpers';
import { upsert } from 'shared/utils';
// Types
export const GET_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes('GET_SERVICE_MEMBER');
export const UPDATE_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes('UPDATE_SERVICE_MEMBER');

const createServiceMemberType = 'CREATE_SERVICE_MEMBER';
export const CREATE_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes(createServiceMemberType);

export const createServiceMember = ReduxHelpers.generateAsyncActionCreator(
  createServiceMemberType,
  CreateServiceMember,
);

const createBackupContactType = 'CREATE_BACKUP_CONTACT';
const indexBackupContactsType = 'INDEX_BACKUP_CONTACTS';
const updateBackupContactType = 'UPDATE_BACKUP_CONTACT';

export const CREATE_BACKUP_CONTACT = ReduxHelpers.generateAsyncActionTypes(createBackupContactType);

export const INDEX_BACKUP_CONTACTS = ReduxHelpers.generateAsyncActionTypes(indexBackupContactsType);

export const UPDATE_BACKUP_CONTACT = ReduxHelpers.generateAsyncActionTypes(updateBackupContactType);

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
    const { currentServiceMember } = state.serviceMember;
    if (currentServiceMember) {
      return UpdateServiceMember(currentServiceMember.id, serviceMember)
        .then(item => dispatch(action.success(Object.assign({}, currentServiceMember, item))))
        .catch(error => dispatch(action.error(error)));
    } else {
      return Promise.reject();
    }
  };
}

export function loadServiceMember(serviceMemberId) {
  const action = ReduxHelpers.generateAsyncActions('GET_SERVICE_MEMBER');
  return function(dispatch, getState) {
    dispatch(action.start);
    const state = getState();
    const { currentServiceMember } = state.serviceMember;
    if (!currentServiceMember) {
      return GetServiceMember(serviceMemberId)
        .then(item => dispatch(action.success(item)))
        .catch(error => dispatch(action.error(error)));
    } else {
      return Promise.resolve();
    }
  };
}

//this is similar to go service_member.IsProfileComplete and we should figure out how to use just one if possible
export const isProfileComplete = state => {
  const sm = get(state, 'serviceMember.currentServiceMember') || {};
  return every([
    sm.rank,
    sm.edipi,
    sm.affiliation,
    sm.first_name,
    sm.last_name,
    sm.telephone,
    sm.personal_email,
    get(sm, 'current_station.id', NULL_UUID) !== NULL_UUID,
    get(sm, 'residential_address.postal_code'),
    get(sm, 'backup_mailing_address.postal_code'),
    get(state, 'serviceMember.currentBackupContacts', []).length > 0,
  ]);
};
// Reducer
const initialState = {
  currentServiceMember: null,
  currentBackupContacts: [],
  hasSubmitError: false,
  hasSubmitSuccess: false,
  createBackupContactSuccess: false,
  updateBackupContactSuccess: false,
};
const reshape = sm => {
  if (!sm) return null;
  return pick(sm, without(Object.keys(sm || {}), 'orders', 'backup_contacts'));
};
const upsertBackUpContact = (contact, state) => {
  const newState = cloneDeep(state);
  upsert(newState.currentBackupContacts, contact);
  return newState;
};
export function serviceMemberReducer(state = initialState, action) {
  switch (action.type) {
    case GET_LOGGED_IN_USER.success:
      return Object.assign({}, state, {
        currentServiceMember: reshape(action.payload.service_member),
        currentBackupContacts: get(action, 'payload.service_member.backup_contacts', []),
        hasLoadError: false,
        hasLoadSuccess: true,
      });
    case CREATE_SERVICE_MEMBER.start:
      return Object.assign({}, state, {
        isLoading: true,
        hasSubmitSuccess: false,
      });
    case CREATE_SERVICE_MEMBER.success:
      return Object.assign({}, state, {
        currentServiceMember: reshape(action.payload),
        isLoading: false,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case CREATE_SERVICE_MEMBER.failure:
      return Object.assign({}, state, {
        isLoading: false,
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
        currentServiceMember: reshape(action.payload),
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
        currentServiceMember: reshape(action.payload),
        currentBackupContacts: action.payload.backup_contacts,
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
      return {
        ...upsertBackUpContact(action.payload, state),
        createBackupContactSuccess: true,
        createBackupContactError: false,
      };
    case CREATE_BACKUP_CONTACT.failure:
      return Object.assign({}, state, {
        createBackupContactSuccess: false,
        createBackupContactError: true,
        error: action.error,
      });
    case UPDATE_BACKUP_CONTACT.start:
      return Object.assign({}, state, {
        updateBackupContactSuccess: false,
      });
    case UPDATE_BACKUP_CONTACT.success:
      return {
        ...upsertBackUpContact(action.payload, state),
        updateBackupContactSuccess: true,
        updateBackupContactError: false,
      };
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
        indexBackupContactsSuccess: false,
        indexBackupContactsError: true,
        error: action.error,
      });
    default:
      return state;
  }
}
