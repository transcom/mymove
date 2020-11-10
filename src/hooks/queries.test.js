import { renderHook } from '@testing-library/react-hooks';

import {
  usePaymentRequestQueries,
  useMoveTaskOrderQueries,
  useOrdersDocumentQueries,
  useMovesQueueQueries,
  usePaymentRequestQueueQueries,
} from './queries';

jest.mock('services/ghcApi', () => ({
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
        },
        b2: {
          shipmentType: 'NTS',
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
  getMoveOrder: (key, id) =>
    Promise.resolve({
      moveOrders: {
        [id]: {
          id,
          uploaded_order_id: '2',
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
      page: 0,
      perPage: 100,
      totalCount: 2,
      queueMoves: [
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
      queuePaymentRequests: [
        {
          id: 'payment1',
        },
        {
          id: 'payment2',
        },
      ],
    }),
}));

describe('usePaymentRequestQueries', () => {
  it('loads data', async () => {
    const testId = 'a1b2';
    const { result, waitForNextUpdate } = renderHook(() => usePaymentRequestQueries(testId));

    expect(result.current).toEqual({
      paymentRequest: undefined,
      paymentRequests: undefined,
      paymentServiceItems: undefined,
      mtoShipments: undefined,
      mtoServiceItems: undefined,
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
      mtoShipments: {
        a1: {
          shipmentType: 'HHG',
        },
        b2: {
          shipmentType: 'NTS',
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
        a1b2: {
          id: 'a1b2',
          uploaded_order_id: '2',
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
          shipmentType: 'NTS',
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

describe('useOrdersDocumentQueries', () => {
  it('loads data', async () => {
    const testMoveOrderId = 'a1b2';
    const { result, waitForNextUpdate } = renderHook(() => useOrdersDocumentQueries(testMoveOrderId));

    expect(result.current).toEqual({
      moveOrders: {
        a1b2: {
          id: 'a1b2',
          uploaded_order_id: '2',
        },
      },
      documents: undefined,
      upload: undefined,
      isLoading: true,
      isError: false,
      isSuccess: false,
    });

    await waitForNextUpdate();

    expect(result.current).toEqual({
      moveOrders: {
        a1b2: {
          id: 'a1b2',
          uploaded_order_id: '2',
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
    const { result, waitForNextUpdate } = renderHook(() => useMovesQueueQueries());

    await waitForNextUpdate();

    expect(result.current).toEqual({
      queueMovesResult: {
        page: 0,
        perPage: 100,
        totalCount: 2,
        queueMoves: [
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
    const { result, waitForNextUpdate } = renderHook(() => usePaymentRequestQueueQueries());

    await waitForNextUpdate();

    expect(result.current).toEqual({
      queuePaymentRequestsResult: {
        page: 0,
        perPage: 100,
        totalCount: 2,
        queuePaymentRequests: [
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
