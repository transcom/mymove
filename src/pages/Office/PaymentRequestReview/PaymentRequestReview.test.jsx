/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';

import { SHIPMENT_OPTIONS } from '../../../shared/constants';

import { PaymentRequestReview } from './PaymentRequestReview';

import { usePaymentRequestQueries } from 'hooks/queries';

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
      mtoServiceItemName: 'Test Service Item',
      priceCents: 12399,
      createdAt: '2020-01-01T00:09:00.999Z',
      status: 'APPROVED',
    },
    2: {
      id: '2',
      mtoServiceItemID: 'b',
      mtoShipmentID: 'b2',
      mtoShipmentType: mockShipmentOptions.NTSR,
      mtoServiceItemName: 'Test Service Item 2',
      priceCents: 45600,
      createdAt: '2020-01-01T00:09:00.999Z',
    },
    3: {
      id: '3',
      mtoServiceItemID: 'c',
      mtoShipmentID: 'a1',
      mtoShipmentType: mockShipmentOptions.HHG,
      mtoServiceItemName: 'Test Service Item 3',
      priceCents: 12312,
      createdAt: '2020-01-01T00:09:00.999Z',
      status: 'DENIED',
    },
    4: {
      id: '4',
      mtoServiceItemID: 'd',
      priceCents: 99999,
      mtoServiceItemName: 'Test Service Item 4',
      createdAt: '2020-01-01T00:09:00.999Z',
    },
  },
  shipmentsPaymentSITBalance: undefined,
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const requiredProps = {
  match: { params: { paymentRequestId: testPaymentRequestId } },
  history: { push: jest.fn() },
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

      render(<PaymentRequestReview {...requiredProps} />);

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      usePaymentRequestQueries.mockReturnValue(errorReturnValue);

      render(<PaymentRequestReview {...requiredProps} />);

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('with data loaded', () => {
    usePaymentRequestQueries.mockReturnValue(usePaymentRequestQueriesReturnValue);
    const wrapper = mount(<PaymentRequestReview {...requiredProps} />);

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

    it('maps the service item card data into the expected format and passes it into the ReviewServiceItems component', () => {
      const reviewServiceItems = wrapper.find('ReviewServiceItems');
      const expectedServiceItemCards = [
        {
          id: '1',
          mtoShipmentID: 'a1',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG,
          mtoShipmentTacType: 'HHG',
          mtoShipmentSacType: 'HHG',
          mtoServiceItemName: 'Test Service Item',
          mtoShipmentModificationType: 'CANCELED',
          mtoShipmentDepartureDate: '2021-05-04',
          mtoShipmentDestinationAddress: 'Fairfield, CA 94535',
          mtoShipmentPickupAddress: 'Beverly Hills, CA 90210',
          amount: 123.99,
          createdAt: '2020-01-01T00:09:00.999Z',
          status: 'APPROVED',
        },
        {
          id: '2',
          mtoShipmentID: 'b2',
          mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
          mtoServiceItemName: 'Test Service Item 2',
          mtoShipmentDepartureDate: '2021-05-04',
          mtoShipmentDestinationAddress: 'Fairfield, CA 94535',
          mtoShipmentPickupAddress: 'Beverly Hills, CA 90210',
          amount: 456.0,
          createdAt: '2020-01-01T00:09:00.999Z',
        },
        {
          id: '3',
          mtoShipmentID: 'a1',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG,
          mtoShipmentTacType: 'HHG',
          mtoShipmentSacType: 'HHG',
          mtoServiceItemName: 'Test Service Item 3',
          mtoShipmentModificationType: 'CANCELED',
          mtoShipmentDepartureDate: '2021-05-04',
          mtoShipmentDestinationAddress: 'Fairfield, CA 94535',
          mtoShipmentPickupAddress: 'Beverly Hills, CA 90210',
          amount: 123.12,
          createdAt: '2020-01-01T00:09:00.999Z',
          status: 'DENIED',
        },
        {
          id: '4',
          mtoServiceItemName: 'Test Service Item 4',
          amount: 999.99,
          createdAt: '2020-01-01T00:09:00.999Z',
        },
      ];

      expect(reviewServiceItems.prop('serviceItemCards')).toEqual(expectedServiceItemCards);
    });
  });
});
