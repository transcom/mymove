import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PaymentRequestCard from './PaymentRequestCard';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';
import { tioRoutes } from 'constants/routes';

jest.mock('hooks/queries', () => ({
  useMovePaymentRequestsQueries: () => {
    const order = {
      sac: '1234456',
      tac: '1213',
    };

    const contractor = {
      contractNumber: 'HTC-123-3456',
    };

    const move = {
      contractor,
      orders: order,
    };

    return {
      paymentRequests: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512e5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
          paymentRequestNumber: '1843-9061-1',
          status: 'REVIEWED',
          reviewedAt: '2020-12-01T00:00:00.000Z',
          moveTaskOrder: move,
          serviceItems: [
            {
              id: '09474c6a-69b6-4501-8e08-670a12512a5f',
              createdAt: '2020-12-01T00:00:00.000Z',
              mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
              priceCents: 2000001,
              status: 'APPROVED',
            },
            {
              id: '39474c6a-69b6-4501-8e08-670a12512a5f',
              createdAt: '2020-12-01T00:00:00.000Z',
              mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
              priceCents: 4000001,
              status: 'DENIED',
              rejectionReason: 'Requested amount exceeds guideline',
            },
          ],
        },
        {
          id: '29474c6a-69b6-4501-8e08-670a12512e5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
          paymentRequestNumber: '1843-9061-2',
          status: 'PENDING',
          moveTaskOrder: move,
          serviceItems: [
            {
              id: '09474c6a-69b6-4501-8e08-670a12512a5f',
              createdAt: '2020-12-01T00:00:00.000Z',
              mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
              mtoShipmentID: 'd81175b7-e26d-4e1e-b1d1-47b17bf4b7f3',
              priceCents: 2000001,
              status: 'REQUESTED',
            },
            {
              id: '39474c6a-69b6-4501-8e08-670a12512a5f',
              createdAt: '2020-12-01T00:00:00.000Z',
              mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
              mtoShipmentID: '9e8222e4-9cdb-4994-8294-6d918a4c684d',
              priceCents: 4000001,
              status: 'REQUESTED',
            },
          ],
        },
      ],
    };
  },
}));

const shipmentInfo = [
  {
    mtoShipmentID: 'd81175b7-e26d-4e1e-b1d1-47b17bf4b7f3',
    shipmentAddress: 'Columbia, SC 29212 to Princeton, NJ 08540',
    departureDate: '2020-12-03T00:00:00.000Z',
  },
  {
    mtoShipmentID: '9e8222e4-9cdb-4994-8294-6d918a4c684d',
    shipmentAddress: 'TBD to Fairfield, CA 94535',
    departureDate: '2020-12-02T00:00:00.000Z',
  },
];
const moveCode = 'AF7K1P';
const dateRegex = /\d{2} [A-Za-z]{3} \d{4}/; // Regex match for DD MMM YYYY

describe('PaymentRequestCard', () => {
  const order = {
    sac: '1234456',
    tac: '1213',
  };

  const contractor = {
    contractNumber: 'HTC-123-3456',
  };

  const move = {
    contractor,
    orders: order,
  };
  const pendingPaymentRequest = {
    id: '29474c6a-69b6-4501-8e08-670a12512e5f',
    createdAt: '2020-12-01T00:00:00.000Z',
    moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
    paymentRequestNumber: '1843-9061-2',
    moveTaskOrder: move,
    status: 'PENDING',
    serviceItems: [
      {
        id: '09474c6a-69b6-4501-8e08-670a12512a5f',
        createdAt: '2020-12-01T00:00:00.000Z',
        mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
        priceCents: 2000001,
        status: 'REQUESTED',
      },
      {
        id: '39474c6a-69b6-4501-8e08-670a12512a5f',
        createdAt: '2020-12-01T00:00:00.000Z',
        mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
        priceCents: 4000001,
        status: 'REQUESTED',
      },
    ],
  };
  const ediErrorPaymentRequest = {
    id: '29474c6a-69b6-4501-8e08-670a12512e5f',
    createdAt: '2020-12-01T00:00:00.000Z',
    moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
    paymentRequestNumber: '1843-9061-2',
    moveTaskOrder: move,
    status: 'PENDING',
    serviceItems: [
      {
        id: '09474c6a-69b6-4501-8e08-670a12512a5f',
        createdAt: '2020-12-01T00:00:00.000Z',
        mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
        priceCents: 2000001,
        status: 'REQUESTED',
      },
      {
        id: '39474c6a-69b6-4501-8e08-670a12512a5f',
        createdAt: '2020-12-01T00:00:00.000Z',
        mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
        priceCents: 4000001,
        status: 'REQUESTED',
      },
    ],
    ediErrorType: '858',
    ediErrorCode: '1A',
    ediErrorDescription: 'Test description',
  };
  const nonWeightReliantPaymentRequest = {
    id: '29474c6a-69b6-4501-8e08-670a12512e5f',
    createdAt: '2020-12-01T00:00:00.000Z',
    moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
    paymentRequestNumber: '1843-9061-2',
    moveTaskOrder: move,
    status: 'PENDING',
    serviceItems: [
      {
        id: '09474c6a-69b6-4501-8e08-670a12512a5f',
        createdAt: '2020-12-01T00:00:00.000Z',
        mtoServiceItemCode: 'MS',
        mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
        priceCents: 2000001,
        status: 'REQUESTED',
      },
    ],
  };
  describe('pending payment request', () => {
    const wrapper = mount(
      <MockProviders
        path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH}
        params={{ moveCode }}
        permissions={[permissionTypes.readPaymentServiceItemStatus]}
      >
        <PaymentRequestCard
          hasBillableWeightIssues={false}
          paymentRequest={pendingPaymentRequest}
          shipmentInfo={shipmentInfo}
        />
      </MockProviders>,
    );

    it('renders the needs review status tag', () => {
      expect(wrapper.find({ 'data-testid': 'tag' }).contains('Needs review')).toBe(true);
    });

    it('sums the service items total', () => {
      expect(wrapper.find('.amountRequested').contains('$60,000.02')).toBe(true);
    });

    it('displays the payment request details ', () => {
      const prDetails = wrapper.find('.footer dd');
      expect(prDetails.contains(contractor.contractNumber)).toBe(true);
    });

    it('renders the view orders link', () => {
      const viewLink = wrapper.find('.footer a');

      expect(viewLink.contains('View orders')).toBe(true);
      expect(viewLink.prop('href')).toBe('/orders');
    });

    it('renders request details toggle drawer by default and hides button', () => {
      const showRequestDetailsButton = wrapper.find('button[data-testid="showRequestDetailsButton"]');

      expect(showRequestDetailsButton.length).toBe(0);
      expect(wrapper.find('[data-testid="toggleDrawer"]').length).toBe(1);
    });

    it('does not render error details toggle drawer by default and hides button', () => {
      render(
        <MockProviders
          path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH}
          params={{ moveCode }}
          permissions={[permissionTypes.updatePaymentServiceItemStatus]}
        >
          <PaymentRequestCard
            paymentRequest={pendingPaymentRequest}
            shipmentInfo={shipmentInfo}
            hasBillableWeightIssues
          />
        </MockProviders>,
      );
      const showErrorDetailsButton = wrapper.find('button[data-testid="showErrorDetailsButton"]');

      expect(showErrorDetailsButton.length).toBe(0);
      expect(wrapper.find('[data-testid="toggleDrawer"]').length).toBe(1);
    });

    it('renders review payment request button disabled when shipment and/or move has billable weight issues', () => {
      render(
        <MockProviders
          path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH}
          params={{ moveCode }}
          permissions={[permissionTypes.updatePaymentServiceItemStatus]}
        >
          <PaymentRequestCard
            paymentRequest={pendingPaymentRequest}
            shipmentInfo={shipmentInfo}
            hasBillableWeightIssues
          />
        </MockProviders>,
      );
      const reviewButton = screen.getByRole('button', { name: 'Review service items' });
      expect(reviewButton).toHaveAttribute('disabled', '');
    });

    it('does not render the review payment request button disabled when shipment and/or move has no billable weight issues', () => {
      render(
        <MockProviders
          path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH}
          params={{ moveCode }}
          permissions={[permissionTypes.updatePaymentServiceItemStatus]}
        >
          <PaymentRequestCard
            paymentRequest={pendingPaymentRequest}
            shipmentInfo={shipmentInfo}
            hasBillableWeightIssues={false}
          />
        </MockProviders>,
      );
      const reviewButton = screen.getByRole('button', { name: 'Review service items' });
      expect(reviewButton).not.toHaveAttribute('disabled', '');
    });
  });

  describe('deprecated payment requests', () => {
    const deprecatedPaymentRequest = {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-2',
      moveTaskOrder: move,
      status: 'DEPRECATED',
      serviceItems: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 2000001,
          status: 'REQUESTED',
        },
      ],
    };

    it('does not have a view documents link', () => {
      render(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
          <PaymentRequestCard
            hasBillableWeightIssues={false}
            paymentRequest={deprecatedPaymentRequest}
            shipmentInfo={shipmentInfo}
          />
        </MockProviders>,
      );
      expect(screen.queryByText('View documents')).not.toBeInTheDocument();
    });

    it('does not have service item status', async () => {
      render(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
          <PaymentRequestCard
            hasBillableWeightIssues={false}
            paymentRequest={deprecatedPaymentRequest}
            shipmentInfo={shipmentInfo}
          />
        </MockProviders>,
      );
      await userEvent.click(screen.getByTestId('showRequestDetailsButton'));
      expect(screen.queryByText('Needs review')).not.toBeInTheDocument();
      expect(screen.getByTestId('deprecated-marker')).toBeInTheDocument();
    });
  });

  describe('reviewed payment request', () => {
    const reviewedPaymentRequest = {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-2',
      status: 'REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED',
      moveTaskOrder: move,
      reviewedAt: '2020-12-01T00:00:00.000Z',
      serviceItems: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 2000001,
          status: 'APPROVED',
        },
        {
          id: '39474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 4000001,
          status: 'DENIED',
          rejectionReason: 'duplicate charge',
        },
      ],
    };

    const mockPDFUpload = {
      contentType: 'application/pdf',
      createdAt: '2020-09-17T16:00:48.099137Z',
      filename: 'test.pdf',
      id: '10',
      status: 'PROCESSING',
      updatedAt: '2020-09-17T16:00:48.099142Z',
      url: '/storage/prime/99/uploads/10?contentType=application%2Fpdf',
    };

    const reviewedPaymentRequestWithDocuments = {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-2',
      status: 'REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED',
      moveTaskOrder: move,
      reviewedAt: '2020-12-01T00:00:00.000Z',
      serviceItems: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 2000001,
          status: 'APPROVED',
        },
        {
          id: '39474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 4000001,
          status: 'DENIED',
          rejectionReason: 'duplicate charge',
        },
      ],
      proofOfServiceDocs: [
        {
          uploads: [mockPDFUpload],
        },
      ],
    };

    const rejectedPaymentRequest = {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-2',
      status: 'REVIEWED',
      reviewedAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrder: move,
      serviceItems: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 2000001,
          status: 'DENIED',
        },
        {
          id: '39474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 4000001,
          status: 'DENIED',
          rejectionReason: 'duplicate charge',
        },
      ],
    };

    const wrapper = mount(
      <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
        <PaymentRequestCard
          hasBillableWeightIssues={false}
          paymentRequest={reviewedPaymentRequest}
          shipmentInfo={shipmentInfo}
        />
      </MockProviders>,
    );

    it('renders the rejected status tag', () => {
      expect(wrapper.find({ 'data-testid': 'tag' }).contains('Rejected')).toBe(true);
    });

    it('sums the approved service items total', () => {
      expect(wrapper.find('.amountAccepted h2').contains('$20,000.01')).toBe(true);
    });

    it('displays the reviewed at date', () => {
      const reviewedAtDate = wrapper.find('.amountAccepted span').at(1).text();
      const reviewedAtDateResult = dateRegex.test(reviewedAtDate);
      expect(reviewedAtDateResult).toBe(true);
    });

    it('sums the rejected service items total', () => {
      expect(wrapper.find('.amountRejected h2').contains('$40,000.01')).toBe(true);
    });

    it('displays the reviewed at date', () => {
      const reviewedAtDate = wrapper.find('.amountRejected span').at(1).text();
      const reviewedAtDateResult = dateRegex.test(reviewedAtDate);
      expect(reviewedAtDateResult).toBe(true);
    });

    it('displays the payment request details ', () => {
      const prDetails = wrapper.find('.footer dd');
      expect(prDetails.contains(contractor.contractNumber)).toBe(true);
    });

    it('renders the view documents link', () => {
      const requestWithDocuments = mount(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
          <PaymentRequestCard
            hasBillableWeightIssues={false}
            paymentRequest={reviewedPaymentRequestWithDocuments}
            shipmentInfo={shipmentInfo}
          />
        </MockProviders>,
      );

      const viewLink = requestWithDocuments.find('.footer a');

      expect(viewLink.text()).toEqual('View documents');
      expect(viewLink.prop('href')).toBe(`payment-requests/${reviewedPaymentRequest.id}`);
    });

    it('renders the no documents text', () => {
      const viewLink = wrapper.find('.footer span');

      expect(viewLink.text()).toEqual('No documents provided');
    });

    it('shows only rejected if no service items are approved', () => {
      const rejected = mount(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
          <PaymentRequestCard
            hasBillableWeightIssues={false}
            paymentRequest={rejectedPaymentRequest}
            shipmentInfo={shipmentInfo}
          />
        </MockProviders>,
      );

      expect(rejected.find('.amountRejected h2').contains('$60,000.02')).toBe(true);
      expect(rejected.find('.amountAccepted').exists()).toBe(false);
    });

    it('renders request details toggle drawer after click', () => {
      const showRequestDetailsButton = wrapper.find('button[data-testid="showRequestDetailsButton"]');
      showRequestDetailsButton.simulate('click');

      expect(wrapper.find('[data-testid="toggleDrawer"]').length).toBe(1);
    });

    it('renders EDI error details toggle drawer after click', () => {
      const ediErrors = mount(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
          <PaymentRequestCard
            hasBillableWeightIssues={false}
            paymentRequest={ediErrorPaymentRequest}
            shipmentInfo={shipmentInfo}
          />
        </MockProviders>,
      );

      const showErrorDetailsButton = ediErrors.find('button[data-testid="showErrorDetailsButton"]');
      showErrorDetailsButton.simulate('click');

      expect(ediErrors.find('[data-testid="toggleDrawer"]').length).toBe(2);
      expect(ediErrors.find('[data-testid="paymentRequestEDIErrorType"]').length).toBe(1);
      expect(ediErrors.find('[data-testid="paymentRequestEDIErrorTypeText"]').length).toBe(1);
      expect(ediErrors.find('[data-testid="paymentRequestEDIErrorCode"]').length).toBe(1);
      expect(ediErrors.find('[data-testid="paymentRequestEDIErrorCodeText"]').length).toBe(1);
      expect(ediErrors.find('[data-testid="paymentRequestEDIErrorDescription"]').length).toBe(1);
      expect(ediErrors.find('[data-testid="paymentRequestEDIErrorDescriptionText"]').length).toBe(1);
    });

    it('renders - for the date it was reviewed at if reviewedAt is null', () => {
      reviewedPaymentRequest.reviewedAt = '';
      const wrapperNoReviewedAtDate = mount(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
          <PaymentRequestCard
            hasBillableWeightIssues={false}
            paymentRequest={reviewedPaymentRequest}
            shipmentInfo={shipmentInfo}
          />
        </MockProviders>,
      );
      const reviewedAtDate = wrapperNoReviewedAtDate.find('.amountRejected span').at(1).text();
      expect(reviewedAtDate).toBe(' on ');
    });
  });

  describe('payment request gex statuses', () => {
    it('renders the EDI Error status tag for edi_error', () => {
      const sentToGexPaymentRequest = {
        id: '29474c6a-69b6-4501-8e08-670a12512e5f',
        createdAt: '2020-12-01T00:00:00.000Z',
        moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
        paymentRequestNumber: '1843-9061-2',
        status: 'EDI_ERROR',
        moveTaskOrder: move,
        serviceItems: [
          {
            id: '09474c6a-69b6-4501-8e08-670a12512a5f',
            createdAt: '2020-12-01T00:00:00.000Z',
            mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
            priceCents: 2000001,
            status: 'DENIED',
          },
          {
            id: '39474c6a-69b6-4501-8e08-670a12512a5f',
            createdAt: '2020-12-01T00:00:00.000Z',
            mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
            priceCents: 4000001,
            status: 'DENIED',
            rejectionReason: 'duplicate charge',
          },
        ],
      };
      const sentToGex = mount(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
          <PaymentRequestCard
            hasBillableWeightIssues={false}
            paymentRequest={sentToGexPaymentRequest}
            shipmentInfo={shipmentInfo}
          />
        </MockProviders>,
      );
      expect(sentToGex.find({ 'data-testid': 'tag' }).contains('EDI Error')).toBe(true);
    });

    const sentToGexPaymentRequest = {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
      paymentRequestNumber: '1843-9061-2',
      status: 'SENT_TO_GEX',
      moveTaskOrder: move,
      serviceItems: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 2000001,
          status: 'APPROVED',
        },
        {
          id: '39474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 4000001,
          status: 'DENIED',
          rejectionReason: 'duplicate charge',
        },
      ],
      sentToGexAt: '2020-12-13T00:00:00.000Z',
    };
    it('renders the Sent to GEX status tag and the date it was sent to gex for sent_to_gex', () => {
      const sentToGex = mount(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
          <PaymentRequestCard
            hasBillableWeightIssues={false}
            paymentRequest={sentToGexPaymentRequest}
            shipmentInfo={shipmentInfo}
          />
        </MockProviders>,
      );
      expect(sentToGex.find({ 'data-testid': 'tag' }).contains('Sent to GEX')).toBe(true);
      expect(sentToGex.find({ 'data-testid': 'sentToGexDetails' }).exists()).toBe(true);
      // displays the sent to gex sum, milmove accepted amount, and milmove rejected amount
      expect(sentToGex.find({ 'data-testid': 'sentToGexDetailsDollarAmountTotal' }).contains('$20,000.01')).toBe(true);
      expect(sentToGex.find({ 'data-testid': 'milMoveAcceptedDetailsDollarAmountTotal' }).contains('$20,000.01')).toBe(
        true,
      );
      expect(sentToGex.find({ 'data-testid': 'milMoveRejectedDetailsDollarAmountTotal' }).contains('$40,000.01')).toBe(
        true,
      );
    });

    it('renders - for the date it was sent to gex if sentToGexAt is null', () => {
      sentToGexPaymentRequest.sentToGexAt = '';
      sentToGexPaymentRequest.reviewedAt = '';

      const sentToGex = mount(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
          <PaymentRequestCard
            hasBillableWeightIssues={false}
            paymentRequest={sentToGexPaymentRequest}
            shipmentInfo={shipmentInfo}
          />
        </MockProviders>,
      );
      expect(sentToGex.find({ 'data-testid': 'tag' }).contains('Sent to GEX')).toBe(true);
      expect(sentToGex.find({ 'data-testid': 'sentToGexDetails' }).exists()).toBe(true);
    });

    it('renders the Tpps Received Status status tag for TPPS_RECEIVED', () => {
      const receivedByGexPaymentRequest = {
        id: '29474c6a-69b6-4501-8e08-670a12512e5f',
        createdAt: '2020-12-01T00:00:00.000Z',
        moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
        paymentRequestNumber: '1843-9061-2',
        status: 'TPPS_RECEIVED',
        moveTaskOrder: move,
        receivedByGexAt: '2020-12-01T00:00:00.000Z',
        serviceItems: [
          {
            id: '09474c6a-69b6-4501-8e08-670a12512a5f',
            createdAt: '2020-12-01T00:00:00.000Z',
            mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
            priceCents: 2000001,
            status: 'APPROVED',
          },
          {
            id: '39474c6a-69b6-4501-8e08-670a12512a5f',
            createdAt: '2020-12-01T00:00:00.000Z',
            mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
            priceCents: 4000001,
            status: 'DENIED',
            rejectionReason: 'duplicate charge',
          },
        ],
      };
      const receivedByGex = mount(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
          <PaymentRequestCard
            hasBillableWeightIssues={false}
            paymentRequest={receivedByGexPaymentRequest}
            shipmentInfo={shipmentInfo}
          />
        </MockProviders>,
      );
      expect(receivedByGex.find({ 'data-testid': 'tag' }).contains('TPPS Received')).toBe(true);
      // displays the tpps received sum, milmove accepted amount, and milmove rejected amount
      expect(receivedByGex.find({ 'data-testid': 'tppsReceivedDetailsDollarAmountTotal' }).contains('$20,000.01')).toBe(
        true,
      );
      expect(
        receivedByGex.find({ 'data-testid': 'milMoveAcceptedDetailsDollarAmountTotal' }).contains('$20,000.01'),
      ).toBe(true);
      expect(
        receivedByGex.find({ 'data-testid': 'milMoveRejectedDetailsDollarAmountTotal' }).contains('$40,000.01'),
      ).toBe(true);
    });

    it('renders the paid status tag for paid request', () => {
      const paidPaymentRequest = {
        id: '29474c6a-69b6-4501-8e08-670a12512e5f',
        createdAt: '2020-12-01T00:00:00.000Z',
        moveTaskOrderID: 'f8c2f97f-99e7-4fb1-9cc4-473debd04dbc',
        paymentRequestNumber: '1843-9061-2',
        status: 'PAID',
        moveTaskOrder: move,
        serviceItems: [
          {
            id: '09474c6a-69b6-4501-8e08-670a12512a5f',
            createdAt: '2020-12-01T00:00:00.000Z',
            mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
            priceCents: 2000001,
            status: 'APPROVED',
          },
          {
            id: '39474c6a-69b6-4501-8e08-670a12512a5f',
            createdAt: '2020-12-01T00:00:00.000Z',
            mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
            priceCents: 4000001,
            status: 'DENIED',
            rejectionReason: 'duplicate charge',
          },
        ],
        tppsInvoiceAmountPaidTotalMillicents: 115155000,
        tppsInvoiceSellerPaidDate: '2024-07-30T00:00:00.000Z',
      };
      const paid = mount(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }}>
          <PaymentRequestCard
            hasBillableWeightIssues={false}
            paymentRequest={paidPaymentRequest}
            shipmentInfo={shipmentInfo}
          />
        </MockProviders>,
      );
      expect(paid.find({ 'data-testid': 'tag' }).contains('TPPS Paid')).toBe(true);
      expect(paid.find({ 'data-testid': 'tppsPaidDetails' }).exists()).toBe(true);
      expect(paid.find({ 'data-testid': 'tppsPaidDetailsDollarAmountTotal' }).exists()).toBe(true);
      // displays the tpps paid sum, milmove accepted amount, and milmove rejected amount
      expect(paid.find({ 'data-testid': 'tppsPaidDetailsDollarAmountTotal' }).contains('$1,151.55')).toBe(true);
      expect(paid.find({ 'data-testid': 'milMoveAcceptedDetailsDollarAmountTotal' }).contains('$20,000.01')).toBe(true);
      expect(paid.find({ 'data-testid': 'milMoveRejectedDetailsDollarAmountTotal' }).contains('$40,000.01')).toBe(true);
    });
  });

  describe('permission dependent rendering', () => {
    it('renders a review service items button when user has TIO permission', () => {
      render(
        <MockProviders
          path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH}
          params={{ moveCode }}
          permissions={[permissionTypes.updatePaymentServiceItemStatus]}
        >
          <PaymentRequestCard paymentRequest={pendingPaymentRequest} shipmentInfo={shipmentInfo} />
        </MockProviders>,
      );
      expect(screen.getByRole('button', { name: 'Review service items' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Review service items' })).not.toHaveAttribute('disabled');
    });

    it('renders the disabled review service items button when user has TIO permission and billable weight issues', () => {
      render(
        <MockProviders
          path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH}
          params={{ moveCode }}
          permissions={[permissionTypes.updatePaymentServiceItemStatus]}
        >
          <PaymentRequestCard
            paymentRequest={pendingPaymentRequest}
            shipmentInfo={shipmentInfo}
            hasBillableWeightIssues
          />
        </MockProviders>,
      );

      expect(screen.getByRole('button', { name: 'Review service items' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Review service items' })).toHaveAttribute('disabled');
    });

    it('renders the disabled review service items button when user has TOO permission', () => {
      render(
        <MockProviders
          path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH}
          params={{ moveCode }}
          permissions={[permissionTypes.readPaymentServiceItemStatus]}
        >
          <PaymentRequestCard paymentRequest={pendingPaymentRequest} shipmentInfo={shipmentInfo} />
        </MockProviders>,
      );

      expect(screen.getByRole('button', { name: 'Review service items' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Review service items' })).toHaveAttribute('disabled');
    });

    it('renders the disabled review service items button when user has TOO permission and billable weight issues', () => {
      render(
        <MockProviders
          path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH}
          params={{ moveCode }}
          permissions={[permissionTypes.readPaymentServiceItemStatus]}
        >
          <PaymentRequestCard
            paymentRequest={pendingPaymentRequest}
            shipmentInfo={shipmentInfo}
            hasBillableWeightIssues
          />
        </MockProviders>,
      );

      expect(screen.getByRole('button', { name: 'Review service items' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Review service items' })).toHaveAttribute('disabled');
    });

    it('does not render the review service items button when user does not have permission', () => {
      render(
        <MockProviders path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH} params={{ moveCode }} permissions={[]}>
          <PaymentRequestCard
            paymentRequest={pendingPaymentRequest}
            shipmentInfo={shipmentInfo}
            hasBillableWeightIssues
          />
        </MockProviders>,
      );

      expect(screen.queryByRole('button', { name: 'Review service items' })).not.toBeInTheDocument();
    });

    it('does not render buttons and disables buttons when the move is locked', () => {
      const isMoveLocked = true;
      render(
        <MockProviders
          path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH}
          params={{ moveCode }}
          permissions={[permissionTypes.updatePaymentServiceItemStatus]}
        >
          <PaymentRequestCard
            paymentRequest={pendingPaymentRequest}
            shipmentInfo={shipmentInfo}
            isMoveLocked={isMoveLocked}
          />
        </MockProviders>,
      );

      expect(screen.queryByRole('button', { name: 'View orders' })).not.toBeInTheDocument();
      expect(screen.queryByRole('button', { name: 'View documents' })).not.toBeInTheDocument();
      expect(screen.queryByRole('button', { name: 'Review service items' })).toBeDisabled();
    });
  });

  it('Review service items is enabled when payment request only contains non weight reliant service items', () => {
    render(
      <MockProviders
        path={tioRoutes.BASE_PAYMENT_REQUESTS_PATH}
        params={{ moveCode }}
        permissions={[permissionTypes.updatePaymentServiceItemStatus]}
      >
        <PaymentRequestCard paymentRequest={nonWeightReliantPaymentRequest} shipmentInfo={shipmentInfo} />
      </MockProviders>,
    );

    expect(screen.queryByRole('button', { name: 'Review service items' })).toBeEnabled();
  });
});
