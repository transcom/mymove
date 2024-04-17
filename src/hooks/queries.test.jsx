import { renderHook } from '@testing-library/react-hooks';
import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

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
  useEditShipmentQueries,
  useEvaluationReportQueries,
  useServicesCounselingQueuePPMQueries,
  useReviewShipmentWeightsQuery,
} from './queries';

import { serviceItemCodes } from 'content/serviceItems';

const queryClient = new QueryClient();
const wrapper = ({ children }) => {
  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
};

jest.mock('services/ghcApi', () => ({
  getCustomer: (key, id) =>
    Promise.resolve({
      customer: { [id]: { id, last_name: 'Kerry', first_name: 'Smith', dodID: '999999999', agency: 'NAVY' } },
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
            shipmentType: 'HHG',
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
        id: 'a1',
        shipmentType: 'HHG',
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
        id: 'b2',
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
      {
        id: 'c3',
        shipmentType: 'PPM',
        ppmShipment: {
          id: 'p1',
          shipmentId: 'c3',
          estimatedWeight: 100,
        },
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
          uploadedAmendedOrderID: '3',
          departmentIndicator: 'Navy',
          grade: 'E-6',
          originDutyLocation: {
            name: 'JBSA Lackland',
          },
          destinationDutyLocation: {
            name: 'JB Lewis-McChord',
          },
          report_by_date: '2018-08-01',
        },
      },
    }),
  getDocument: (key, id) =>
    Promise.resolve({
      documents: {
        [id]: {
          id,
          uploads: [`${id}0`],
        },
      },
      upload: {
        id: `${id}0`,
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
  getEvaluationReportByID: () =>
    Promise.resolve({
      id: '1234',
      type: 'SHIPMENT',
      moveReferenceID: '4321',
      shipmentID: '123',
    }),
  getMTOShipmentByID: () =>
    Promise.resolve({
      id: '12345',
      moveTaskOrderId: '67890',
      customerRemarks: 'mock remarks',
      requestedPickupDate: '2020-03-01',
      requestedDeliveryDate: '2020-03-30',
    }),
  getReportViolationsByReportID: () =>
    Promise.resolve([
      {
        id: '123',
        reportID: '456',
        violationID: '789',
      },
    ]),
  getShipmentsPaymentSITBalance: () =>
    Promise.resolve({
      shipmentsPaymentSITBalance: {
        a1: {
          pendingBilledEndDate: '2021-08-29',
          pendingSITDaysInvoiced: 30,
          previouslyBilledDays: 0,
          shipmentID: 'a1',
          totalSITDaysAuthorized: 90,
          totalSITDaysRemaining: 60,
          totalSITEndDate: '2021-10-29',
        },
      },
    }),
  getServicesCounselingPPMQueue: () =>
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
  getPPMDocuments: (key, shipmentID) => {
    if (shipmentID !== 'c3') {
      return { MovingExpenses: [], ProGearWeightTickets: [], WeightTickets: [] };
    }
    return Promise.resolve({
      MovingExpenses: [],
      ProGearWeightTickets: [],
      WeightTickets: [
        {
          emptyWeight: 14500,
          fullWeight: 18500,
          id: 'ppmDoc1',
          missingEmptyWeightTicket: false,
          missingFullWeightTicket: false,
          ownsTrailer: false,
          vehicleDescription: '2022 Honda CR-V Hybrid',
        },
      ],
    });
  },
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
    const { result, waitFor } = renderHook(() => useTXOMoveInfoQueries(testMoveCode), { wrapper });

    expect(result.current).toEqual({
      order: undefined,
      customerData: undefined,
      isLoading: true,
      isError: false,
      isSuccess: false,
    });

    await waitFor(() => result.current.isSuccess);

    expect(result.current).toEqual({
      customerData: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999', agency: 'NAVY' },
      order: {
        id: '4321',
        customerID: '2468',
        customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
        uploaded_order_id: '2',
        uploadedAmendedOrderID: '3',
        departmentIndicator: 'Navy',
        grade: 'E-6',
        originDutyLocation: {
          name: 'JBSA Lackland',
        },
        destinationDutyLocation: {
          name: 'JB Lewis-McChord',
        },
        report_by_date: '2018-08-01',
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
      move: {
        id: '1234',
        ordersId: '4321',
        moveCode: 'ABCDEF',
      },
    });
  });
});

describe('usePaymentRequestQueries', () => {
  it('loads data', async () => {
    const testId = 'a1b2';
    const { result, waitFor } = renderHook(() => usePaymentRequestQueries(testId), { wrapper });

    expect(result.current).toEqual({
      paymentRequest: undefined,
      paymentRequests: undefined,
      paymentServiceItems: undefined,
      mtoShipments: undefined,
      shipmentsPaymentSITBalance: undefined,
      isLoading: true,
      isError: false,
      isSuccess: false,
    });

    await waitFor(() => result.current.isSuccess);

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
      mtoShipments: [
        {
          id: 'a1',
          shipmentType: 'HHG',
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
          id: 'b2',
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
        {
          id: 'c3',
          shipmentType: 'PPM',
          ppmShipment: {
            id: 'p1',
            shipmentId: 'c3',
            estimatedWeight: 100,
          },
        },
      ],
      shipmentsPaymentSITBalance: {
        a1: {
          pendingBilledEndDate: '2021-08-29',
          pendingSITDaysInvoiced: 30,
          previouslyBilledDays: 0,
          shipmentID: 'a1',
          totalSITDaysAuthorized: 90,
          totalSITDaysRemaining: 60,
          totalSITEndDate: '2021-10-29',
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
    const { result, waitFor } = renderHook(() => useMoveDetailsQueries(moveCode), { wrapper });

    expect(result.current).toEqual({
      move: {
        id: '1234',
        ordersId: '4321',
        moveCode: 'ABCDEF',
      },
      closeoutOffice: undefined,
      customerData: {
        id: '2468',
        last_name: 'Kerry',
        first_name: 'Smith',
        dodID: '999999999',
        agency: 'NAVY',
      },
      order: {
        id: '4321',
        customerID: '2468',
        customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
        uploaded_order_id: '2',
        uploadedAmendedOrderID: '3',
        departmentIndicator: 'Navy',
        grade: 'E-6',
        originDutyLocation: {
          name: 'JBSA Lackland',
        },
        destinationDutyLocation: {
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

    await waitFor(() => result.current.isSuccess);

    expect(result.current).toEqual({
      move: {
        id: '1234',
        ordersId: '4321',
        moveCode: 'ABCDEF',
      },
      closeoutOffice: undefined,
      customerData: {
        id: '2468',
        last_name: 'Kerry',
        first_name: 'Smith',
        dodID: '999999999',
        agency: 'NAVY',
      },
      order: {
        id: '4321',
        customerID: '2468',
        customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
        uploaded_order_id: '2',
        uploadedAmendedOrderID: '3',
        departmentIndicator: 'Navy',
        grade: 'E-6',
        originDutyLocation: {
          name: 'JBSA Lackland',
        },
        destinationDutyLocation: {
          name: 'JB Lewis-McChord',
        },
        report_by_date: '2018-08-01',
      },
      mtoShipments: [
        {
          id: 'a1',
          shipmentType: SHIPMENT_OPTIONS.HHG,
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
          id: 'b2',
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
        {
          id: 'c3',
          shipmentType: 'PPM',
          ppmShipment: {
            id: 'p1',
            shipmentId: 'c3',
            estimatedWeight: 100,
            weightTickets: [
              {
                emptyWeight: 14500,
                fullWeight: 18500,
                id: 'ppmDoc1',
                missingEmptyWeightTicket: false,
                missingFullWeightTicket: false,
                ownsTrailer: false,
                vehicleDescription: '2022 Honda CR-V Hybrid',
              },
            ],
            reviewShipmentWeightsURL: '/counseling/moves/ABCDEF/shipments/c3/document-review',
          },
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

describe('useMoveTaskOrderQueries', () => {
  it('loads data', async () => {
    const moveId = 'ABCDEF';
    const { result, waitForNextUpdate } = renderHook(() => useMoveTaskOrderQueries(moveId), { wrapper });

    await waitForNextUpdate();

    expect(result.current).toEqual({
      orders: {
        4321: {
          id: '4321',
          customerID: '2468',
          customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
          uploaded_order_id: '2',
          uploadedAmendedOrderID: '3',
          departmentIndicator: 'Navy',
          grade: 'E-6',
          originDutyLocation: {
            name: 'JBSA Lackland',
          },
          destinationDutyLocation: {
            name: 'JB Lewis-McChord',
          },
          report_by_date: '2018-08-01',
        },
      },
      move: {
        id: '1234',
        moveCode: 'ABCDEF',
        ordersId: '4321',
      },
      mtoShipments: [
        {
          id: 'a1',
          shipmentType: SHIPMENT_OPTIONS.HHG,
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
          id: 'b2',
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
        {
          id: 'c3',
          shipmentType: 'PPM',
          ppmShipment: {
            id: 'p1',
            shipmentId: 'c3',
            estimatedWeight: 100,
            weightTickets: [
              {
                emptyWeight: 14500,
                fullWeight: 18500,
                id: 'ppmDoc1',
                missingEmptyWeightTicket: false,
                missingFullWeightTicket: false,
                ownsTrailer: false,
                vehicleDescription: '2022 Honda CR-V Hybrid',
              },
            ],
            reviewShipmentWeightsURL: '/counseling/moves/ABCDEF/shipments/c3/document-review',
          },
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

describe('useEditShipmentQueries', () => {
  it('loads data', async () => {
    const moveCode = 'ABCDEF';
    const { result, waitForNextUpdate } = renderHook(() => useEditShipmentQueries(moveCode), { wrapper });

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
        uploadedAmendedOrderID: '3',
        departmentIndicator: 'Navy',
        grade: 'E-6',
        originDutyLocation: {
          name: 'JBSA Lackland',
        },
        destinationDutyLocation: {
          name: 'JB Lewis-McChord',
        },
        report_by_date: '2018-08-01',
      },
      mtoShipments: [
        {
          id: 'a1',
          shipmentType: SHIPMENT_OPTIONS.HHG,
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
          id: 'b2',
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
        {
          id: 'c3',
          shipmentType: 'PPM',
          ppmShipment: {
            id: 'p1',
            shipmentId: 'c3',
            estimatedWeight: 100,
          },
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

    const { result, waitForNextUpdate } = renderHook(() => useOrdersDocumentQueries(testLocatorId), {
      wrapper,
    });

    await waitForNextUpdate();

    expect(result.current).toEqual({
      move: { id: '1234', ordersId: '4321', moveCode: testLocatorId },
      orders: {
        4321: {
          id: '4321',
          customerID: '2468',
          customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
          uploaded_order_id: '2',
          uploadedAmendedOrderID: '3',
          departmentIndicator: 'Navy',
          grade: 'E-6',
          originDutyLocation: {
            name: 'JBSA Lackland',
          },
          destinationDutyLocation: {
            name: 'JB Lewis-McChord',
          },
          report_by_date: '2018-08-01',
        },
      },
      documents: {
        2: {
          id: '2',
          uploads: ['20'],
        },
      },
      upload: {
        id: '20',
      },
      amendedDocuments: {
        3: {
          id: '3',
          uploads: ['30'],
        },
      },
      amendedUpload: {
        id: '30',
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    });
  });
});

describe('useMovesQueueQueries', () => {
  it('loads data', async () => {
    const { result, waitForNextUpdate } = renderHook(
      () => useMovesQueueQueries({ filters: [], currentPage: 1, currentPageSize: 100 }),
      { wrapper },
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
    const { result, waitForNextUpdate } = renderHook(
      () => usePaymentRequestQueueQueries({ filters: [], currentPage: 1, currentPageSize: 100 }),
      { wrapper },
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
    const { result, waitForNextUpdate } = renderHook(() => useUserQueries(), { wrapper });

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

describe('useEvaluationReportQueries', () => {
  it('loads data', async () => {
    const { result, waitFor } = renderHook(() => useEvaluationReportQueries('1234'), { wrapper });

    expect(result.current).toEqual({
      evaluationReport: {},
      mtoShipment: {},
      reportViolations: [],
      isLoading: true,
      isError: false,
      isSuccess: false,
    });

    await waitFor(() => result.current.isSuccess);

    expect(result.current).toEqual({
      evaluationReport: { id: '1234', moveReferenceID: '4321', type: 'SHIPMENT', shipmentID: '123' },
      mtoShipment: {
        id: '12345',
        moveTaskOrderId: '67890',
        customerRemarks: 'mock remarks',
        requestedPickupDate: '2020-03-01',
        requestedDeliveryDate: '2020-03-30',
      },
      reportViolations: [
        {
          id: '123',
          reportID: '456',
          violationID: '789',
        },
      ],
      isLoading: false,
      isError: false,
      isSuccess: true,
    });
  });
});

describe('useServicesCounselingQueuePPMQueries', () => {
  it('loads data', async () => {
    const { result, waitFor } = renderHook(() => useServicesCounselingQueuePPMQueries('1234'), { wrapper });

    await waitFor(() => result.current.isSuccess);

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

describe('useReviewShipmentWeightsQuery', () => {
  it('loads data', async () => {
    const { result, waitFor } = renderHook(() => useReviewShipmentWeightsQuery('ABCDEF'), { wrapper });

    await waitFor(() => result.current.isSuccess);

    expect(result.current).toEqual({
      move: { id: '1234', ordersId: '4321', moveCode: 'ABCDEF' },
      orders: {
        4321: {
          id: '4321',
          customerID: '2468',
          customer: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
          uploaded_order_id: '2',
          uploadedAmendedOrderID: '3',
          departmentIndicator: 'Navy',
          grade: 'E-6',
          originDutyLocation: {
            name: 'JBSA Lackland',
          },
          destinationDutyLocation: {
            name: 'JB Lewis-McChord',
          },
          report_by_date: '2018-08-01',
        },
      },
      mtoShipments: [
        {
          id: 'a1',
          shipmentType: 'HHG',
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
          id: 'b2',
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
        {
          id: 'c3',
          shipmentType: 'PPM',
          ppmShipment: {
            id: 'p1',
            shipmentId: 'c3',
            estimatedWeight: 100,
            weightTickets: [
              {
                emptyWeight: 14500,
                fullWeight: 18500,
                id: 'ppmDoc1',
                missingEmptyWeightTicket: false,
                missingFullWeightTicket: false,
                ownsTrailer: false,
                vehicleDescription: '2022 Honda CR-V Hybrid',
              },
            ],
            reviewShipmentWeightsURL: '/counseling/moves/ABCDEF/shipments/c3/document-review',
          },
        },
      ],
      isLoading: false,
      isError: false,
      isSuccess: true,
    });
  });
});
