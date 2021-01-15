/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import { SHIPMENT_OPTIONS } from '../../../shared/constants';

import { PaymentRequestReview } from './PaymentRequestReview';

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
  id: '11',
  status: 'PROCESSING',
  updatedAt: '11',
  url: '/storage/prime/99/uploads/10?contentType=image%2Fpng',
};

const mockShipmentOptions = SHIPMENT_OPTIONS;

jest.mock('hooks/queries', () => ({
  usePaymentRequestQueries: () => {
    const testPaymentRequestId = 'test-payment-id-123';
    return {
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
      isLoading: false,
      isError: false,
      isSuccess: true,
    };
  },
}));

describe('PaymentRequestReview', () => {
  const testPaymentRequestId = 'test-payment-id-123';

  const requiredProps = {
    match: { params: { paymentRequestId: testPaymentRequestId } },
    history: { push: jest.fn() },
  };

  describe('with data loaded', () => {
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
          mtoServiceItemName: 'Test Service Item',
          amount: 123.99,
          createdAt: '2020-01-01T00:09:00.999Z',
          status: 'APPROVED',
        },
        {
          id: '2',
          mtoShipmentID: 'b2',
          mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
          mtoServiceItemName: 'Test Service Item 2',
          amount: 456.0,
          createdAt: '2020-01-01T00:09:00.999Z',
        },
        {
          id: '3',
          mtoShipmentID: 'a1',
          mtoShipmentType: SHIPMENT_OPTIONS.HHG,
          mtoServiceItemName: 'Test Service Item 3',
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
