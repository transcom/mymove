/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { PaymentRequestReview } from './PaymentRequestReview';

import { patchPaymentServiceItemStatus } from 'services/ghcApi';
import { SHIPMENT_OPTIONS, PAYMENT_REQUEST_STATUS, PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { usePaymentRequestQueries } from 'hooks/queries';
import { ReactQueryWrapper } from 'testUtils';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => jest.fn(),
}));

const mockPDFUpload = {
  contentType: 'application/pdf',
  createdAt: '2020-09-17T16:00:48.099137Z',
  filename: 'test.pdf',
  id: '10',
  status: 'PROCESSING',
  updatedAt: '2020-09-17T16:00:48.099142Z',
  url: '/storage/prime/99/uploads/10?contentType=application%2Fpdf',
};

const mockJPGUpload = {
  contentType: 'image/jpg',
  createdAt: '2020-09-17T16:00:48.099137Z',
  filename: 'test.jpg',
  id: '11',
  status: 'PROCESSING',
  updatedAt: '11',
  url: '/storage/prime/99/uploads/10?contentType=image%2Fjpeg',
};

const mockPNGUpload = {
  contentType: 'image/png',
  createdAt: '2020-09-17T16:00:48.099137Z',
  filename: 'test.png',
  id: '12',
  status: 'PROCESSING',
  updatedAt: '12',
  url: '/storage/prime/99/uploads/10?contentType=image%2Fpng',
};

const mockShipmentOptions = SHIPMENT_OPTIONS;

const testPaymentRequestId = 'test-payment-id-123';

jest.mock('hooks/queries', () => ({
  usePaymentRequestQueries: jest.fn(),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  patchPaymentServiceItemStatus: jest.fn(),
}));

// prevents react-fileviewer from throwing errors without mocking relevant DOM elements
jest.mock('components/DocumentViewer/Content/Content', () => {
  const MockContent = () => <div>Content</div>;
  return MockContent;
});

const usePaymentRequestQueriesReturnValue = {
  paymentRequest: {
    id: testPaymentRequestId,
    moveTaskOrderID: '123',
    proofOfServiceDocs: [
      {
        uploads: [mockPDFUpload],
      },
      {
        uploads: [mockJPGUpload, mockPNGUpload],
      },
    ],
  },
  paymentRequests: {
    [testPaymentRequestId]: {
      id: testPaymentRequestId,
      moveTaskOrderID: '123',
    },
  },
  mtoShipments: [
    {
      actualPickupDate: '2021-05-04',
      destinationAddress: {
        city: 'Fairfield',
        id: '1f0054d2-6cf9-41aa-8b9c-27b9fc6667b5',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      id: 'a1',
      status: 'CANCELED',
      moveTaskOrderID: '123',
      pickupAddress: {
        city: 'Beverly Hills',
        id: '0cf43b1f-04e8-4c56-a6a1-06aec192ca07',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
    },
    {
      actualPickupDate: '2021-05-04',
      destinationAddress: {
        city: 'Fairfield',
        id: '1f0054d2-6cf9-41aa-8b9c-27b9fc6667b5',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      id: 'b2',
      moveTaskOrderID: '123',
      pickupAddress: {
        city: 'Beverly Hills',
        id: '0cf43b1f-04e8-4c56-a6a1-06aec192ca07',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
    },
  ],
  paymentServiceItems: {
    1: {
      id: '1',
      mtoServiceItemID: 'a',
      mtoShipmentID: 'a1',
      mtoShipmentType: mockShipmentOptions.HHG,
      mtoServiceItemName: 'Test Service Item 1',
      priceCents: 12399,
      createdAt: '2020-01-01T00:09:00.999Z',
      status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    },
    3: {
      id: '3',
      mtoServiceItemID: 'b',
      mtoShipmentID: 'b2',
      mtoShipmentType: mockShipmentOptions.NTSR,
      mtoServiceItemName: 'Test Service Item 3',
      priceCents: 45600,
      createdAt: '2020-01-02T00:09:00.999Z',
      status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    },
    2: {
      id: '2',
      mtoServiceItemID: 'c',
      mtoShipmentID: 'a1',
      mtoShipmentType: mockShipmentOptions.HHG,
      mtoServiceItemName: 'Test Service Item 2',
      priceCents: 12312,
      createdAt: '2020-01-03T00:09:00.999Z',
      status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
    },
    4: {
      id: '4',
      mtoServiceItemID: 'd',
      priceCents: 99999,
      mtoServiceItemName: 'Test Service Item 4',
      createdAt: '2020-01-04T00:09:00.999Z',
      status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    },
  },
  shipmentsPaymentSITBalance: undefined,
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const usePaymentRequestQueriesReturnNoDocs = {
  ...usePaymentRequestQueriesReturnValue,
  paymentRequest: {
    ...usePaymentRequestQueriesReturnValue.paymentRequest,
    proofOfServiceDocs: [],
  },
};

const usePaymentRequestQueriesReturnValuePending = {
  ...usePaymentRequestQueriesReturnValue,
  paymentRequest: {
    ...usePaymentRequestQueriesReturnValue.paymentRequest,
    status: PAYMENT_REQUEST_STATUS.PENDING,
  },
};

const reviewedPaymentServiceItems = {
  ...usePaymentRequestQueriesReturnValuePending.paymentServiceItems,
  2: {
    ...usePaymentRequestQueriesReturnValuePending.paymentServiceItems[2],
    rejectionReason: 'duplicate charge',
    thing: 'stuff',
  },
  3: {
    ...usePaymentRequestQueriesReturnValuePending.paymentServiceItems[3],
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
  },
  4: {
    ...usePaymentRequestQueriesReturnValuePending.paymentServiceItems[4],
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
  },
};

const usePaymentRequestQueriesReturnValuePendingFinalReview = {
  ...usePaymentRequestQueriesReturnValuePending,
  paymentServiceItems: reviewedPaymentServiceItems,
};

const usePaymentRequestQueriesReturnValueApproved = {
  ...usePaymentRequestQueriesReturnValue,
  paymentServiceItems: reviewedPaymentServiceItems,
};

const requiredProps = {
  order: { tac: '1234', sac: '5678', ntsTac: 'AB12', ntsSac: 'CD34' },
};

const loadingReturnValue = {
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  isLoading: false,
  isError: true,
  isSuccess: false,
};

describe('PaymentRequestReview', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      usePaymentRequestQueries.mockReturnValue(loadingReturnValue);

      render(
        <ReactQueryWrapper>
          <PaymentRequestReview {...requiredProps} />
        </ReactQueryWrapper>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      usePaymentRequestQueries.mockReturnValue(errorReturnValue);

      render(
        <ReactQueryWrapper>
          <PaymentRequestReview {...requiredProps} />
        </ReactQueryWrapper>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('with data loaded and no docs', () => {
    it('renders text instead of the DocumentViewer', () => {
      usePaymentRequestQueries.mockReturnValue(usePaymentRequestQueriesReturnNoDocs);
      render(
        <ReactQueryWrapper>
          <PaymentRequestReview {...requiredProps} />
        </ReactQueryWrapper>,
      );
      const h2 = screen.getByRole('heading', { name: 'No documents provided' });
      expect(h2).toBeInTheDocument();
    });
  });
  describe('with data loaded', () => {
    usePaymentRequestQueries.mockReturnValue(usePaymentRequestQueriesReturnValue);
    const wrapper = mount(
      <ReactQueryWrapper>
        <PaymentRequestReview {...requiredProps} />
      </ReactQueryWrapper>,
    );
    it('renders without errors', () => {
      expect(wrapper.find('[data-testid="PaymentRequestReview"]').exists()).toBe(true);
    });
    it('renders the DocumentViewer', () => {
      const documentViewer = wrapper.find('DocumentViewer');
      expect(documentViewer.exists()).toBe(true);
      expect(documentViewer.prop('files')).toEqual([mockPDFUpload, mockJPGUpload, mockPNGUpload]);
    });
    it('renders the ReviewServiceItems sidebar', () => {
      expect(wrapper.find('ReviewServiceItems').exists()).toBe(true);
    });
    it('maps the service item card data into the expected format and passes it into the ReviewServiceItems component, ordered by timestamp ascending', () => {
      const reviewServiceItems = wrapper.find('ReviewServiceItems');
      const expectedServiceItemCards = [
        {
          id: '1',
          mtoShipmentID: 'a1',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG,
          mtoShipmentTacType: 'HHG',
          mtoShipmentSacType: 'HHG',
          mtoServiceItemName: 'Test Service Item 1',
          mtoShipmentModificationType: 'CANCELED',
          mtoShipmentDepartureDate: '2021-05-04',
          mtoShipmentDestinationAddress: 'Fairfield, CA 94535',
          mtoShipmentPickupAddress: 'Beverly Hills, CA 90210',
          amount: 123.99,
          createdAt: '2020-01-01T00:09:00.999Z',
          status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
        },
        {
          id: '2',
          mtoShipmentID: 'a1',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG,
          mtoShipmentTacType: 'HHG',
          mtoShipmentSacType: 'HHG',
          mtoServiceItemName: 'Test Service Item 2',
          mtoShipmentModificationType: 'CANCELED',
          mtoShipmentDepartureDate: '2021-05-04',
          mtoShipmentDestinationAddress: 'Fairfield, CA 94535',
          mtoShipmentPickupAddress: 'Beverly Hills, CA 90210',
          amount: 123.12,
          createdAt: '2020-01-03T00:09:00.999Z',
          status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
        },
        {
          id: '3',
          mtoShipmentID: 'b2',
          mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
          mtoServiceItemName: 'Test Service Item 3',
          mtoShipmentDepartureDate: '2021-05-04',
          mtoShipmentDestinationAddress: 'Fairfield, CA 94535',
          mtoShipmentPickupAddress: 'Beverly Hills, CA 90210',
          amount: 456.0,
          createdAt: '2020-01-02T00:09:00.999Z',
          status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
        },
        {
          id: '4',
          mtoServiceItemName: 'Test Service Item 4',
          amount: 999.99,
          createdAt: '2020-01-04T00:09:00.999Z',
          status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
        },
      ];
      expect(reviewServiceItems.prop('serviceItemCards')).toEqual(expectedServiceItemCards);
    });
  });

  describe('clicking the next button', () => {
    describe('with pending requests', () => {
      beforeEach(async () => {
        jest.clearAllMocks();
        usePaymentRequestQueries.mockReturnValue(usePaymentRequestQueriesReturnValuePending);
        render(
          <ReactQueryWrapper>
            <PaymentRequestReview {...requiredProps} />
          </ReactQueryWrapper>,
        );
        expect(await screen.findByText('1 OF 4 ITEMS')).toBeInTheDocument();
        expect(await screen.findByText(/Test Service Item 1/)).toBeInTheDocument();
        const nextButton = screen.getByRole('button', { name: 'Next Service Item' });
        const returnValue = {
          paymentServiceItems: {
            ...usePaymentRequestQueriesReturnValue.paymentServiceItems,
          },
        };
        patchPaymentServiceItemStatus.mockReturnValue(returnValue);
        await userEvent.click(nextButton);
      });
      it('submits the form with expected payload', async () => {
        const payload = {
          ifMatchEtag: undefined,
          moveTaskOrderID: '123',
          paymentServiceItemID: '1',
          rejectionReason: '',
          status: 'APPROVED',
        };
        expect(patchPaymentServiceItemStatus).toHaveBeenCalledWith(payload);
      });
      it('navigates forward', async () => {
        expect(await screen.findByText('2 OF 4 ITEMS')).toBeInTheDocument();
      });
      it('shows an error if Reject is selected and no reason is provided', async () => {
        const rejectButton = screen.getByRole('radio', { name: 'Reject' });
        const nextButton = screen.getByRole('button', { name: 'Next Service Item' });
        await fireEvent.click(rejectButton);
        await userEvent.click(nextButton);
        expect(screen.getByText(/Add a reason why this service item is rejected/)).toBeInTheDocument();
      });
      it('rejects an item when a reason is provided', async () => {
        const rejectButton = screen.getByLabelText('Reject');
        const nextButton = screen.getByRole('button', { name: 'Next Service Item' });
        const reasonInput = screen.getByRole('textbox', { name: 'Reason for rejection' });
        fireEvent.click(rejectButton);
        await userEvent.type(reasonInput, 'duplicate charge');
        await userEvent.click(nextButton);
        await waitFor(() => {
          expect(rejectButton).toBeChecked();
        });
        expect(screen.queryByText(/Add a reason why this service item is rejected/)).not.toBeInTheDocument();
        const payload = {
          ifMatchEtag: undefined,
          moveTaskOrderID: '123',
          paymentServiceItemID: '2',
          rejectionReason: 'duplicate charge',
          status: 'DENIED',
        };
        expect(patchPaymentServiceItemStatus).toBeCalledWith(payload);
      });
      describe('can navigate to the final review', () => {
        it('with an incomplete review and and finish reviewing', async () => {
          // second item is loaded from the previous step
          expect(await screen.findByText('2 OF 4 ITEMS')).toBeInTheDocument();
          expect(await screen.findByText('Test Service Item 2')).toBeInTheDocument();
          expect(screen.getByRole('radio', { name: 'Reject' })).toBeChecked();
          const reasonInput = screen.getByRole('textbox', { name: 'Reason for rejection' });
          await userEvent.type(reasonInput, 'duplicate charge');
          const nextButton = screen.getByRole('button', { name: 'Next Service Item' });
          await userEvent.click(nextButton);
          expect(await screen.findByText('3 OF 4 ITEMS')).toBeInTheDocument();
          expect(await screen.findByText('Test Service Item 3')).toBeInTheDocument();
          expect(screen.getByRole('radio', { name: 'Reject' })).not.toBeChecked();
          expect(screen.getByRole('radio', { name: 'Approve' })).not.toBeChecked();
          await userEvent.click(nextButton);
          expect(await screen.findByText('4 OF 4 ITEMS')).toBeInTheDocument();
          expect(await screen.findByText('Test Service Item 4')).toBeInTheDocument();
          expect(screen.getByRole('radio', { name: 'Reject' })).not.toBeChecked();
          expect(screen.getByRole('radio', { name: 'Approve' })).not.toBeChecked();
          await userEvent.click(nextButton);
          expect(screen.getByRole('heading', { level: 2, text: 'Complete request' })).toBeInTheDocument();
          const finishReviewButton = screen.getByRole('button', { name: 'Finish review' });
          await userEvent.click(finishReviewButton);
          expect(await screen.findByText('3 OF 4 ITEMS')).toBeInTheDocument();
        });
      });
    });
    describe('with a complete review, pending request', () => {
      it('can navigate and authorize payment', async () => {
        usePaymentRequestQueries.mockReturnValue(usePaymentRequestQueriesReturnValuePendingFinalReview);
        render(
          <ReactQueryWrapper>
            <PaymentRequestReview {...requiredProps} />
          </ReactQueryWrapper>,
        );
        expect(screen.getByText('1 OF 4 ITEMS')).toBeInTheDocument();
        expect(screen.getByText(/Test Service Item 1/)).toBeInTheDocument();
        expect(screen.getByRole('radio', { name: 'Approve' })).toBeChecked();

        const nextButton = screen.getByRole('button', { name: 'Next Service Item' });
        await userEvent.click(nextButton);
        expect(screen.getByText('2 OF 4 ITEMS')).toBeInTheDocument();
        expect(screen.getByText(/Test Service Item 2/)).toBeInTheDocument();
        expect(screen.getByRole('radio', { name: 'Reject' })).toBeChecked();
        expect(screen.getByText('duplicate charge')).toBeInTheDocument();
        await userEvent.click(nextButton);
        expect(await screen.getByText('3 OF 4 ITEMS')).toBeInTheDocument();
        await userEvent.click(nextButton);
        expect(await screen.getByText('4 OF 4 ITEMS')).toBeInTheDocument();
        await userEvent.click(nextButton);
        expect(screen.getByRole('button', { name: 'Authorize payment' })).toBeInTheDocument();
      });
    });
    describe('with an approved review', () => {
      beforeEach(() => {
        usePaymentRequestQueries.mockReturnValue(usePaymentRequestQueriesReturnValueApproved);
        render(
          <ReactQueryWrapper>
            <PaymentRequestReview {...requiredProps} />
          </ReactQueryWrapper>,
        );
      });
      it('shows the review details', async () => {
        expect(screen.getByRole('heading', { name: 'Review details', level: 3 })).toBeInTheDocument();
        const terms = screen.getAllByRole('term');
        expect(terms[0]).toHaveTextContent('Requested');
        expect(terms[1]).toHaveTextContent('Accepted');
        expect(terms[2]).toHaveTextContent('Rejected');
        const definitions = screen.getAllByRole('definition');
        expect(definitions[0]).toHaveTextContent('$1,703.10');
        expect(definitions[1]).toHaveTextContent('$1,579.98');
        expect(definitions[2]).toHaveTextContent('$123.12');
      });
      it('navigates back, and shows the correct icons for approved and rejected cards', async () => {
        await userEvent.click(screen.getByRole('button', { name: 'Back' }));
        expect(screen.getByTestId('statusHeading')).toHaveTextContent('Accepted');
        await userEvent.click(screen.getByRole('button', { name: 'Previous Service Item' }));
        expect(screen.getByTestId('statusHeading')).toHaveTextContent('Accepted');
        await userEvent.click(screen.getByRole('button', { name: 'Previous Service Item' }));
        expect(screen.getByTestId('statusHeading')).toHaveTextContent('Rejected');
        await userEvent.click(screen.getByRole('button', { name: 'Previous Service Item' }));
        expect(screen.getByTestId('statusHeading')).toHaveTextContent('Accepted');
      });
    });
  });
});
