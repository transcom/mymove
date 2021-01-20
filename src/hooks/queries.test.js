import { renderHook } from '@testing-library/react-hooks';

import {
  usePaymentRequestQueries,
  useMoveTaskOrderQueries,
  useOrdersDocumentQueries,
  useMovesQueueQueries,
  usePaymentRequestQueueQueries,
  useUserQueries,
  useTXOMoveInfoQueries,
  useMoveDetailsQueries,
} from './queries';

jest.mock('services/ghcApi', () => ({
  getCustomer: (key, id) =>
    Promise.resolve({
      customer: { [id]: { id, last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' } },
    }),
  getPaymentRequest: (key, id) =>
    Promise.resolve({
      paymentRequests: {
        [id]: {
          moveTaskOrderID: '123',
        },
      },
      paymentServiceItems: {},
    }),
  getMTOShipments: () =>
    Promise.resolve({
      mtoShipments: {
        a1: {
          shipmentType: 'HHG',
          // mtoAgents: [],
          // mtoServiceItems: [],
        },
        b2: {
          shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
          // mtoAgents: [],
          // mtoServiceItems: [],
        },
      },
    }),
  getMTOServiceItems: () =>
    Promise.resolve({
      mtoServiceItems: {
        a: {
          reServiceName: 'Test Service Item',
        },
        b: {
          reServiceName: 'Test Service Item 2',
        },
      },
    }),
  getMove: () =>
    Promise.resolve({
      id: '1234',
      ordersId: '4321',
      moveCode: 'ABCDEF',
    }),
  getMoveOrder: (key, id) =>
    Promise.resolve({
      moveOrders: {
        [id]: {
          id,
          customerID: '2468',
          customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
          uploaded_order_id: '2',
          departmentIndicator: 'Navy',
          grade: 'E-6',
          originDutyStation: {
            name: 'JBSA Lackland',
          },
          destinationDutyStation: {
            name: 'JB Lewis-McChord',
          },
          report_by_date: '2018-08-01',
        },
      },
    }),
  getMoveTaskOrderList: () =>
    Promise.resolve({
      moveTaskOrders: {
        1: {
          id: '1',
        },
      },
    }),
  getDocument: () =>
    Promise.resolve({
      documents: {
        2: {
          id: '2',
          uploads: ['z'],
        },
      },
      upload: {
        id: 'z',
      },
    }),
  getMovesQueue: () =>
    Promise.resolve({
      page: 1,
      perPage: 100,
      totalCount: 2,
      data: [
        {
          id: 'move1',
        },
        {
          id: 'move2',
        },
      ],
    }),
  getPaymentRequestsQueue: () =>
    Promise.resolve({
      page: 0,
      perPage: 100,
      totalCount: 2,
      data: [
        {
          id: 'payment1',
        },
        {
          id: 'payment2',
        },
      ],
    }),
  getLoggedInUserQueries: () =>
    Promise.resolve({
      data: {},
    }),
}));

jest.mock('services/internalApi', () => ({
  getLoggedInUserQueries: () =>
    Promise.resolve({
      office_user: { transportation_office: { gbloc: 'LMKG' } },
    }),
}));

describe('useTXOMoveInfoQueries', () => {
  it('loads data', async () => {
    const testMoveCode = 'ABCDEF';
    const { result, waitForNextUpdate } = renderHook(() => useTXOMoveInfoQueries(testMoveCode));

    expect(result.current).toEqual({
      moveOrder: undefined,
      customerData: undefined,
      isLoading: true,
      isError: false,
      isSuccess: false,
    });

    await waitForNextUpdate();

    expect(result.current).toEqual({
      customerData: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
      moveOrder: {
        id: '4321',
        customerID: '2468',
        customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
        uploaded_order_id: '2',
        departmentIndicator: 'Navy',
        grade: 'E-6',
        originDutyStation: {
          name: 'JBSA Lackland',
        },
        destinationDutyStation: {
          name: 'JB Lewis-McChord',
        },
        report_by_date: '2018-08-01',
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    });
  });
});

describe('usePaymentRequestQueries', () => {
  it('loads data', async () => {
    const testId = 'a1b2';
    const { result, waitForNextUpdate } = renderHook(() => usePaymentRequestQueries(testId));

    expect(result.current).toEqual({
      paymentRequest: undefined,
      paymentRequests: undefined,
      paymentServiceItems: undefined,
      isLoading: true,
      isError: false,
      isSuccess: false,
    });

    await waitForNextUpdate();

    expect(result.current).toEqual({
      paymentRequest: {
        moveTaskOrderID: '123',
      },
      paymentRequests: {
        a1b2: {
          moveTaskOrderID: '123',
        },
      },
      paymentServiceItems: {},
      isLoading: false,
      isError: false,
      isSuccess: true,
    });
  });
});

describe('useMoveTaskOrderQueries', () => {
  it('loads data', async () => {
    const testMoveOrderId = 'a1b2';
    const { result, waitForNextUpdate } = renderHook(() => useMoveTaskOrderQueries(testMoveOrderId));

    expect(result.current).toEqual({
      moveOrders: undefined,
      moveTaskOrders: undefined,
      mtoShipments: undefined,
      mtoServiceItems: undefined,
      isLoading: true,
      isError: false,
      isSuccess: false,
    });

    await waitForNextUpdate();

    expect(result.current).toEqual({
      moveOrders: {
        4321: {
          id: '4321',
          customerID: '2468',
          customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
          uploaded_order_id: '2',
          departmentIndicator: 'Navy',
          grade: 'E-6',
          originDutyStation: {
            name: 'JBSA Lackland',
          },
          destinationDutyStation: {
            name: 'JB Lewis-McChord',
          },
          report_by_date: '2018-08-01',
        },
      },
      moveTaskOrders: {
        1: {
          id: '1',
        },
      },
      mtoShipments: {
        a1: {
          shipmentType: 'HHG',
        },
        b2: {
          shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
        },
      },
      mtoServiceItems: {
        a: {
          reServiceName: 'Test Service Item',
        },
        b: {
          reServiceName: 'Test Service Item 2',
        },
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    });
  });
});

describe('useMoveDetailsQueries', () => {
  it('loads data', async () => {
    const moveCode = 'ABCDEF';
    const { result, waitForNextUpdate } = renderHook(() => useMoveDetailsQueries(moveCode));

    expect(result.current).toEqual({
      move: {
        id: '1234',
        ordersId: '4321',
        moveCode: 'ABCDEF',
      },
      moveOrders: {
        4321: {
          id: '4321',
          customerID: '2468',
          customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
          uploaded_order_id: '2',
          departmentIndicator: 'Navy',
          grade: 'E-6',
          originDutyStation: {
            name: 'JBSA Lackland',
          },
          destinationDutyStation: {
            name: 'JB Lewis-McChord',
          },
          report_by_date: '2018-08-01',
        },
      },
      mtoShipments: undefined,
      isLoading: true,
      isError: false,
      isSuccess: false,
    });

    await waitForNextUpdate();

    expect(result.current).toEqual({
      move: {
        id: '1234',
        ordersId: '4321',
        moveCode: 'ABCDEF',
      },
      moveOrders: {
        4321: {
          id: '4321',
          customerID: '2468',
          customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
          uploaded_order_id: '2',
          departmentIndicator: 'Navy',
          grade: 'E-6',
          originDutyStation: {
            name: 'JBSA Lackland',
          },
          destinationDutyStation: {
            name: 'JB Lewis-McChord',
          },
          report_by_date: '2018-08-01',
        },
      },
      mtoShipments: {
        a1: {
          shipmentType: 'HHG',
        },
        b2: {
          shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
        },
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    });
  });
});

describe('useOrdersDocumentQueries', () => {
  it('loads data', async () => {
    const testLocatorId = 'ABCDEF';
    const { result, waitForNextUpdate } = renderHook(() => useOrdersDocumentQueries(testLocatorId));

    await waitForNextUpdate();

    expect(result.current).toEqual({
      move: { id: '1234', ordersId: '4321', moveCode: testLocatorId },
      moveOrders: {
        4321: {
          id: '4321',
          customerID: '2468',
          customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
          uploaded_order_id: '2',
          departmentIndicator: 'Navy',
          grade: 'E-6',
          originDutyStation: {
            name: 'JBSA Lackland',
          },
          destinationDutyStation: {
            name: 'JB Lewis-McChord',
          },
          report_by_date: '2018-08-01',
        },
      },
      documents: {
        2: {
          id: '2',
          uploads: ['z'],
        },
      },
      upload: {
        id: 'z',
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    });
  });
});

describe('useMovesQueueQueries', () => {
  it('loads data', async () => {
    const { result, waitForNextUpdate } = renderHook(() =>
      useMovesQueueQueries({ filters: [], currentPage: 1, currentPageSize: 100 }),
    );

    await waitForNextUpdate();

    expect(result.current).toEqual({
      queueResult: {
        page: 1,
        perPage: 100,
        totalCount: 2,
        data: [
          {
            id: 'move1',
          },
          {
            id: 'move2',
          },
        ],
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    });
  });
});

describe('usePaymentRequestsQueueQueries', () => {
  it('loads data', async () => {
    const { result, waitForNextUpdate } = renderHook(() =>
      usePaymentRequestQueueQueries({ filters: [], currentPage: 1, currentPageSize: 100 }),
    );

    await waitForNextUpdate();

    expect(result.current).toEqual({
      queueResult: {
        page: 0,
        perPage: 100,
        totalCount: 2,
        data: [
          {
            id: 'payment1',
          },
          {
            id: 'payment2',
          },
        ],
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    });
  });
});

describe('useUserQueries', () => {
  it('loads data', async () => {
    const { result, waitForNextUpdate } = renderHook(() => useUserQueries());

    await waitForNextUpdate();

    expect(result.current).toEqual({
      data: {
        office_user: { transportation_office: { gbloc: 'LMKG' } },
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    });
  });
});
