import React from 'react';
import { mount } from 'enzyme';

import { PaymentRequestReview } from './PaymentRequestReview';

const testPaymentRequestId = 'test-payment-id-123';
const testMTOID = 'test-mto-id-456';

const mockGetPaymentRequest = jest.fn(() =>
  Promise.resolve({
    entities: {
      paymentRequests: {
        [testPaymentRequestId]: {
          moveTaskOrderID: testMTOID,
        },
      },
    },
  }),
);

describe('PaymentRequestReview', () => {
  const requiredProps = {
    match: { params: { paymentRequestId: testPaymentRequestId } },
    getPaymentRequest: mockGetPaymentRequest,
    getMTOServiceItems: jest.fn(() => Promise.resolve()),
    getMTOShipments: jest.fn(() => Promise.resolve()),
    patchPaymentServiceItemStatus: jest.fn(),
    history: { push: jest.fn() },
    updatePaymentRequest: jest.fn(),
  };

  describe('with or without data loaded', () => {
    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = mount(<PaymentRequestReview {...requiredProps} />);

    it('renders without errors', () => {
      expect(wrapper.find('[data-testid="PaymentRequestReview"]').exists()).toBe(true);
    });

    it('loads the required API data on mount', () => {
      expect(requiredProps.getPaymentRequest).toHaveBeenCalledWith(testPaymentRequestId);
      expect(requiredProps.getMTOServiceItems).toHaveBeenCalledWith(testMTOID);
      expect(requiredProps.getMTOShipments).toHaveBeenCalledWith(testMTOID);
    });

    it('renders the document viewer', () => {
      expect(wrapper.find('DocumentViewer').exists()).toBe(true);
    });
  });

  describe('with data loaded', () => {
    const dataProps = {
      paymentServiceItems: [
        {
          id: '1',
          mtoServiceItemID: 'a',
          priceCents: 12399,
          createdAt: '2020-01-01T00:09:00.999Z',
          status: 'APPROVED',
        },
        {
          id: '2',
          mtoServiceItemID: 'b',
          priceCents: 45600,
          createdAt: '2020-01-01T00:09:00.999Z',
        },
        {
          id: '3',
          mtoServiceItemID: 'c',
          priceCents: 12312,
          createdAt: '2020-01-01T00:09:00.999Z',
          status: 'DENIED',
        },
        {
          id: '4',
          mtoServiceItemID: 'd',
          priceCents: 99999,
          createdAt: '2020-01-01T00:09:00.999Z',
        },
      ],
      mtoServiceItems: [
        {
          id: 'a',
          mtoShipmentID: 'a1',
          reServiceName: 'Test Service Item',
        },
        {
          id: 'b',
          mtoShipmentID: 'b2',
          reServiceName: 'Test Service Item 2',
        },
        {
          id: 'c',
          mtoShipmentID: 'a1',
          reServiceName: 'Test Service Item 3',
        },
        {
          id: 'd',
          reServiceName: 'Test Service Item 4',
        },
      ],
      mtoShipments: [
        {
          id: 'a1',
          shipmentType: 'HHG',
        },
        {
          id: 'b2',
          shipmentType: 'NTS',
        },
      ],
    };

    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = mount(<PaymentRequestReview {...requiredProps} {...dataProps} />);

    it('renders the ReviewServiceItems sidebar', () => {
      expect(wrapper.find('ReviewServiceItems').exists()).toBe(true);
    });

    it('maps the service item card data into the expected format and passes it into the ReviewServiceItems component', () => {
      const reviewServiceItems = wrapper.find('ReviewServiceItems');
      const expectedServiceItemCards = [
        {
          id: '1',
          shipmentId: 'a1',
          shipmentType: 'HHG',
          serviceItemName: 'Test Service Item',
          amount: 123.99,
          createdAt: '2020-01-01T00:09:00.999Z',
          status: 'APPROVED',
        },
        {
          id: '2',
          shipmentId: 'b2',
          shipmentType: 'NTS',
          serviceItemName: 'Test Service Item 2',
          amount: 456.0,
          createdAt: '2020-01-01T00:09:00.999Z',
        },
        {
          id: '3',
          shipmentId: 'a1',
          shipmentType: 'HHG',
          serviceItemName: 'Test Service Item 3',
          amount: 123.12,
          createdAt: '2020-01-01T00:09:00.999Z',
          status: 'DENIED',
        },
        {
          id: '4',
          serviceItemName: 'Test Service Item 4',
          amount: 999.99,
          createdAt: '2020-01-01T00:09:00.999Z',
        },
      ];

      expect(reviewServiceItems.prop('serviceItemCards')).toEqual(expectedServiceItemCards);
    });
  });
});
