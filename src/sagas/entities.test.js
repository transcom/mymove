import { all, takeLatest, put, call } from 'redux-saga/effects';

import { watchUpdateEntities, updateServiceMember, updateBackupContact } from './entities';

import { UPDATE_SERVICE_MEMBER, UPDATE_BACKUP_CONTACT } from 'store/entities/actions';
import { normalizeResponse } from 'services/swaggerRequest';
import { addEntities } from 'shared/Entities/actions';

describe('watchUpdateEntities', () => {
  const generator = watchUpdateEntities();

  it('takes the latest update actions and calls update sagas', () => {
    expect(generator.next().value).toEqual(
      all([
        takeLatest(UPDATE_SERVICE_MEMBER, updateServiceMember),
        takeLatest(UPDATE_BACKUP_CONTACT, updateBackupContact),
      ]),
    );
  });

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});

describe('updateServiceMember', () => {
  const testAction = {
    payload: {
      service_member: {
        id: 'testServiceMemberId',
        orders: [{ id: 'testorder1' }, { id: 'testorder2' }],
      },
    },
  };

  const normalizedServiceMember = normalizeResponse(testAction.payload, 'serviceMember');

  const generator = updateServiceMember(testAction);

  it('normalizes the payload', () => {
    expect(generator.next().value).toEqual(call(normalizeResponse, testAction.payload, 'serviceMember'));
  });

  it('stores the normalized data in entities', () => {
    expect(generator.next(normalizedServiceMember).value).toEqual(put(addEntities(normalizedServiceMember)));
  });

  it('calls the legacy UPDATE_SERVICE_MEMBER_SUCCESS action with the raw payload', () => {
    expect(generator.next().value).toEqual(
      put({
        type: 'UPDATE_SERVICE_MEMBER_SUCCESS',
        payload: testAction.payload,
      }),
    );
  });

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});

describe('updateBackupContact', () => {
  const testAction = {
    payload: {
      created_at: '2020-11-17T17:55:43.745Z',
      email: 'newron@example.com',
      id: '3bbf4b50-975f-459e-922c-47e1d5afb538',
      name: 'John Lee New',
      permission: 'NONE',
      service_member_id: '39c1dcea-6e3b-4d80-9b4d-5cc66d215f61',
      telephone: '999-999-9999',
      updated_at: '2020-11-17T17:56:33.081Z',
    },
  };

  const normalizedBackupContact = normalizeResponse(testAction.payload, 'backupContact');

  const generator = updateBackupContact(testAction);

  it('normalizes the payload', () => {
    expect(generator.next().value).toEqual(call(normalizeResponse, testAction.payload, 'backupContact'));
  });

  it('stores the normalized data in entities', () => {
    expect(generator.next(normalizedBackupContact).value).toEqual(put(addEntities(normalizedBackupContact)));
  });

  it('calls the legacy UPDATE_BACKUP_CONTACT_SUCCESS action with the raw payload', () => {
    expect(generator.next().value).toEqual(
      put({
        type: 'UPDATE_BACKUP_CONTACT_SUCCESS',
        payload: testAction.payload,
      }),
    );
  });

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});
