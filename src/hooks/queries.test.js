import { renderHook } from '@testing-library/react-hooks';

import { SHIPMENT_OPTIONS } from '../shared/constants';

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

import { serviceItemCodes } from 'content/serviceItems';

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
  getMTOShipments: (key, id, normalize) => {
    if (normalize) {
      return Promise.resolve({
        mtoShipments: {
          a1: {
            shipmentType: 'HHG_LONGHAUL_DOMESTIC',
            mtoAgents: [
              {
                agentType: 'RELEASING_AGENT',
                mtoShipmentID: 'a1',
              },
              {
                agentType: 'RECEIVING_AGENT',
                mtoShipmentID: 'a1',
              },
            ],
            mtoServiceItems: [
              {
                reServiceName: 'Domestic linehaul',
              },
              {
                reServiceName: 'Fuel surcharge',
              },
            ],
          },
          b2: {
            shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
            mtoAgents: [
              {
                agentType: 'RELEASING_AGENT',
                mtoShipmentID: 'b2',
              },
              {
                agentType: 'RECEIVING_AGENT',
                mtoShipmentID: 'b2',
              },
            ],
            mtoServiceItems: [
              {
                reServiceName: 'Domestic origin price',
              },
              {
                reServiceName: 'Domestic unpacking',
              },
            ],
          },
        },
      });
    }
    return Promise.resolve([
      {
        shipmentType: 'HHG_LONGHAUL_DOMESTIC',
        mtoAgents: [
          {
            agentType: 'RELEASING_AGENT',
            mtoShipmentID: 'a1',
          },
          {
            agentType: 'RECEIVING_AGENT',
            mtoShipmentID: 'a1',
          },
        ],
        mtoServiceItems: [
          {
            reServiceName: 'Domestic linehaul',
          },
          {
            reServiceName: 'Fuel surcharge',
          },
        ],
      },
      {
        shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
        mtoAgents: [
          {
            agentType: 'RELEASING_AGENT',
            mtoShipmentID: 'b2',
          },
          {
            agentType: 'RECEIVING_AGENT',
            mtoShipmentID: 'b2',
          },
        ],
        mtoServiceItems: [
          {
            reServiceName: 'Domestic origin price',
          },
          {
            reServiceName: 'Domestic unpacking',
          },
        ],
      },
    ]);
  },
  getMTOServiceItems: (key, id, normalize) => {
    if (normalize) {
      return Promise.resolve({
        mtoServiceItems: {
          a: {
            reServiceName: 'Counseling',
          },
          b: {
            reServiceName: 'Move management',
          },
        },
      });
    }
    return Promise.resolve([
      {
        id: 'a',
        reServiceName: 'Counseling',
      },
      {
        id: 'b',
        reServiceName: 'Move management',
      },
    ]);
  },
  getMove: () =>
    Promise.resolve({
      id: '1234',
      ordersId: '4321',
      moveCode: 'ABCDEF',
    }),
  getOrder: (key, id) =>
    Promise.resolve({
      orders: {
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
      order: undefined,
      customerData: undefined,
      isLoading: true,
      isError: false,
      isSuccess: false,
    });

    await waitForNextUpdate();

    expect(result.current).toEqual({
      customerData: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
      order: {
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
    const testOrderId = 'a1b2';
    const { result, waitForNextUpdate } = renderHook(() => useMoveTaskOrderQueries(testOrderId));

    expect(result.current).toEqual({
      orders: undefined,
      moveTaskOrders: undefined,
      mtoShipments: undefined,
      mtoServiceItems: undefined,
      isLoading: true,
      isError: false,
      isSuccess: false,
    });

    await waitForNextUpdate();

    expect(result.current).toEqual({
      orders: {
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
      mtoShipments: [
        {
          shipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoAgents: [
            {
              agentType: 'RELEASING_AGENT',
              mtoShipmentID: 'a1',
            },
            {
              agentType: 'RECEIVING_AGENT',
              mtoShipmentID: 'a1',
            },
          ],
          mtoServiceItems: [
            {
              reServiceName: serviceItemCodes.DLH,
            },
            {
              reServiceName: serviceItemCodes.FSC,
            },
          ],
        },
        {
          shipmentType: SHIPMENT_OPTIONS.NTSR,
          mtoAgents: [
            {
              agentType: 'RELEASING_AGENT',
              mtoShipmentID: 'b2',
            },
            {
              agentType: 'RECEIVING_AGENT',
              mtoShipmentID: 'b2',
            },
          ],
          mtoServiceItems: [
            {
              reServiceName: serviceItemCodes.DOP,
            },
            {
              reServiceName: serviceItemCodes.DUPK,
            },
          ],
        },
      ],
      mtoServiceItems: [
        {
          id: 'a',
          reServiceName: serviceItemCodes.CS,
        },
        {
          id: 'b',
          reServiceName: serviceItemCodes.MS,
        },
      ],
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
      order: {
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
      mtoShipments: undefined,
      mtoServiceItems: undefined,
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
      order: {
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
      mtoShipments: [
        {
          shipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
          mtoAgents: [
            {
              agentType: 'RELEASING_AGENT',
              mtoShipmentID: 'a1',
            },
            {
              agentType: 'RECEIVING_AGENT',
              mtoShipmentID: 'a1',
            },
          ],
          mtoServiceItems: [
            {
              reServiceName: serviceItemCodes.DLH,
            },
            {
              reServiceName: serviceItemCodes.FSC,
            },
          ],
        },
        {
          shipmentType: SHIPMENT_OPTIONS.NTSR,
          mtoAgents: [
            {
              agentType: 'RELEASING_AGENT',
              mtoShipmentID: 'b2',
            },
            {
              agentType: 'RECEIVING_AGENT',
              mtoShipmentID: 'b2',
            },
          ],
          mtoServiceItems: [
            {
              reServiceName: serviceItemCodes.DOP,
            },
            {
              reServiceName: serviceItemCodes.DUPK,
            },
          ],
        },
      ],
      mtoServiceItems: [
        {
          id: 'a',
          reServiceName: serviceItemCodes.CS,
        },
        {
          id: 'b',
          reServiceName: serviceItemCodes.MS,
        },
      ],
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
      orders: {
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
