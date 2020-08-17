import { renderHook } from '@testing-library/react-hooks';

import { usePaymentRequestQueries, useMoveTaskOrderQueries } from './queries';

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
  getMoveOrder: () =>
    Promise.resolve({
      moveOrders: {
        1: {
          id: '1',
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
    const testId = 'a1b2';
    const { result, waitForNextUpdate } = renderHook(() => useMoveTaskOrderQueries(testId));

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
        1: {
          id: '1',
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
