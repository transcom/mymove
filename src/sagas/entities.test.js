import { all, takeLatest, put, call } from 'redux-saga/effects';

import {
  watchUpdateEntities,
  updateServiceMember,
  updateBackupContact,
  updateMove,
  updateMTOShipment,
  updateOrders,
  updatePPMs,
  updatePPM,
  updatePPMEstimate,
  updatePPMSitEstimate,
} from './entities';

import {
  UPDATE_SERVICE_MEMBER,
  UPDATE_BACKUP_CONTACT,
  UPDATE_MOVE,
  UPDATE_MTO_SHIPMENT,
  UPDATE_ORDERS,
  UPDATE_PPMS,
  UPDATE_PPM,
  UPDATE_PPM_ESTIMATE,
  UPDATE_PPM_SIT_ESTIMATE,
} from 'store/entities/actions';
import { normalizeResponse } from 'services/swaggerRequest';
import { addEntities } from 'shared/Entities/actions';

describe('watchUpdateEntities', () => {
  const generator = watchUpdateEntities();

  it('takes the latest update actions and calls update sagas', () => {
    expect(generator.next().value).toEqual(
      all([
        takeLatest(UPDATE_SERVICE_MEMBER, updateServiceMember),
        takeLatest(UPDATE_BACKUP_CONTACT, updateBackupContact),
        takeLatest(UPDATE_ORDERS, updateOrders),
        takeLatest(UPDATE_MOVE, updateMove),
        takeLatest(UPDATE_MTO_SHIPMENT, updateMTOShipment),
        takeLatest(UPDATE_PPMS, updatePPMs),
        takeLatest(UPDATE_PPM, updatePPM),
        takeLatest(UPDATE_PPM_ESTIMATE, updatePPMEstimate),
        takeLatest(UPDATE_PPM_SIT_ESTIMATE, updatePPMSitEstimate),
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

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});

describe('updateMove', () => {
  const testAction = {
    payload: {
      created_at: '2020-12-07T17:03:58.767Z',
      id: '3a8c9f4f-7344-4f18-9ab5-0de3ef57b901',
      locator: 'ONEHHG',
      orders_id: 'a413144b-137f-4400-85c2-a99c437ef85e',
      selected_move_type: 'HHG',
      service_member_id: '1d06ab96-cb72-4013-b159-321d6d29c6eb',
      status: 'DRAFT',
      updated_at: '2020-12-07T22:41:08.999Z',
    },
  };

  const normalizedMove = normalizeResponse(testAction.payload, 'move');

  const generator = updateMove(testAction);

  it('normalizes the payload', () => {
    expect(generator.next().value).toEqual(call(normalizeResponse, testAction.payload, 'move'));
  });

  it('stores the normalized data in entities', () => {
    expect(generator.next(normalizedMove).value).toEqual(put(addEntities(normalizedMove)));
  });

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});

describe('updateMTOShipment', () => {
  const testAction = {
    payload: {
      createdAt: '2020-12-08T17:39:05.051Z',
      customerRemarks: '',
      eTag: 'MjAyMC0xMi0wOFQxNzozOTowNS4wNTE0Mzha',
      id: 'dfb950b2-075f-4595-b3cd-a56bb7293ca1',
      moveTaskOrderID: 'c7621194-ed1e-4ead-9033-761715f179f3',
      pickupAddress: {
        city: 'New York',
        id: '3eedb9b9-fbbb-47c0-93c7-50daaae84cb4',
        postalCode: '10002',
        state: 'NY',
        streetAddress1: '415 Grand St, #E306, #E306',
        streetAddress2: '#E306',
      },
      requestedPickupDate: '2020-12-22',
      shipmentType: 'HHG_INTO_NTS_DOMESTIC',
      status: 'SUBMITTED',
      updatedAt: '2020-12-08T17:39:05.051Z',
    },
  };

  const normalizedMove = normalizeResponse(testAction.payload, 'mtoShipment');

  const generator = updateMTOShipment(testAction);

  it('normalizes the payload', () => {
    expect(generator.next().value).toEqual(call(normalizeResponse, testAction.payload, 'mtoShipment'));
  });

  it('stores the normalized data in entities', () => {
    expect(generator.next(normalizedMove).value).toEqual(put(addEntities(normalizedMove)));
  });

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});

describe('updateOrders', () => {
  const testAction = {
    payload: {
      created_at: '2020-12-17T15:54:48.853Z',
      has_dependents: false,
      id: 'ef45eb5a-c1bf-4c60-9c22-990500b6badc',
      issue_date: '2020-12-22',
      moves: [
        {
          created_at: '2020-12-17T15:54:48.873Z',
          id: '0ff5ec27-57be-4760-a87f-42998aa94caf',
          locator: 'C8PFDW',
          orders_id: 'ef45eb5a-c1bf-4c60-9c22-990500b6badc',
          selected_move_type: '',
          service_member_id: '15a17300-e1c6-4b3a-8e5d-9c47782a3961',
          status: 'DRAFT',
          updated_at: '2020-12-17T15:54:48.873Z',
        },
      ],
      new_duty_location: {
        address: {
          city: 'Glendale Luke AFB',
          country: 'United States',
          id: 'ce6ec9a4-1bad-4fb3-8b3c-89ebee54e8cf',
          postalCode: '85309',
          state: 'AZ',
          streetAddress1: 'n/a',
        },
        address_id: 'ce6ec9a4-1bad-4fb3-8b3c-89ebee54e8cf',
        affiliation: 'AIR_FORCE',
        created_at: '2020-12-07T17:02:33.987Z',
        id: '9e1b519e-6daa-4c3f-8cfe-413c582b6366',
        name: 'Luke AFB',
        updated_at: '2020-12-07T17:02:33.987Z',
      },
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      report_by_date: '2020-12-28',
      service_member_id: '15a17300-e1c6-4b3a-8e5d-9c47782a3961',
      spouse_has_pro_gear: false,
      status: 'DRAFT',
      updated_at: '2020-12-17T15:54:48.853Z',
      uploaded_orders: {
        id: '251ea83d-4295-4105-9780-3ae2d6549872',
        service_member_id: '15a17300-e1c6-4b3a-8e5d-9c47782a3961',
        uploads: [],
      },
    },
  };

  const normalizedOrders = normalizeResponse(testAction.payload, 'orders');

  const generator = updateOrders(testAction);

  it('normalizes the payload', () => {
    expect(generator.next().value).toEqual(call(normalizeResponse, testAction.payload, 'orders'));
  });

  it('stores the normalized data in entities', () => {
    expect(generator.next(normalizedOrders).value).toEqual(put(addEntities(normalizedOrders)));
  });

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});

describe('updatePPM', () => {
  const testAction = {
    payload: {
      actual_move_date: '2020-12-18',
      advance_worksheet: {
        id: '00000000-0000-0000-0000-000000000000',
        service_member_id: '00000000-0000-0000-0000-000000000000',
        uploads: [],
      },
      approve_date: '2020-12-21T22:45:52.000Z',
      created_at: '2020-12-21T22:43:48.278Z',
      destination_postal_code: '99619',
      has_additional_postal_code: false,
      has_requested_advance: false,
      has_sit: false,
      id: 'd9488eac-eef8-430e-8c4b-05884c3cc6fa',
      move_id: '2b8198ca-e70a-40b7-822e-be5527bf0606',
      original_move_date: '2020-12-19',
      pickup_postal_code: '10002',
      status: 'PAYMENT_REQUESTED',
      submit_date: '2020-12-21T22:45:12.100Z',
      updated_at: '2020-12-21T22:46:50.805Z',
    },
  };

  const normalizedPPM = normalizeResponse(testAction.payload, 'personallyProcuredMove');

  const generator = updatePPM(testAction);

  it('normalizes the payload', () => {
    expect(generator.next().value).toEqual(call(normalizeResponse, testAction.payload, 'personallyProcuredMove'));
  });

  it('stores the normalized data in entities', () => {
    expect(generator.next(normalizedPPM).value).toEqual(put(addEntities(normalizedPPM)));
  });

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});

describe('updatePPMs', () => {
  const testAction = {
    payload: [
      {
        actual_move_date: '2020-12-18',
        advance_worksheet: {
          id: '00000000-0000-0000-0000-000000000000',
          service_member_id: '00000000-0000-0000-0000-000000000000',
          uploads: [],
        },
        approve_date: '2020-12-21T22:45:52.000Z',
        created_at: '2020-12-21T22:43:48.278Z',
        destination_postal_code: '99619',
        has_additional_postal_code: false,
        has_requested_advance: false,
        has_sit: false,
        id: 'd9488eac-eef8-430e-8c4b-05884c3cc6fa',
        move_id: '2b8198ca-e70a-40b7-822e-be5527bf0606',
        original_move_date: '2020-12-19',
        pickup_postal_code: '10002',
        status: 'PAYMENT_REQUESTED',
        submit_date: '2020-12-21T22:45:12.100Z',
        updated_at: '2020-12-21T22:46:50.805Z',
      },
    ],
  };

  const normalizedPPM = normalizeResponse(testAction.payload, 'personallyProcuredMoves');

  const generator = updatePPMs(testAction);

  it('normalizes the payload', () => {
    expect(generator.next().value).toEqual(call(normalizeResponse, testAction.payload, 'personallyProcuredMoves'));
  });

  it('stores the normalized data in entities', () => {
    expect(generator.next(normalizedPPM).value).toEqual(put(addEntities(normalizedPPM)));
  });

  it('is done', () => {
    expect(generator.next().done).toEqual(true);
  });
});
