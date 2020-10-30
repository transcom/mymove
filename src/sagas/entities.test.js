import { all, takeLatest, put, call } from 'redux-saga/effects';

import { watchUpdateEntities, updateServiceMember } from './entities';

import { UPDATE_SERVICE_MEMBER } from 'store/entities/actions';
import { normalizeResponse } from 'services/swaggerRequest';
import { addEntities } from 'shared/Entities/actions';

describe('watchUpdateEntities', () => {
  const generator = watchUpdateEntities();

  it('takes the latest update actions and calls update sagas', () => {
    expect(generator.next().value).toEqual(all([takeLatest(UPDATE_SERVICE_MEMBER, updateServiceMember)]));
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
