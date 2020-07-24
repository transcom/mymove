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

    it('renders the ReviewServiceItems sidebar', () => {
      expect(wrapper.find('ReviewServiceItems').exists()).toBe(true);
    });
  });

  describe('with data loaded', () => {
    const dataProps = {
      paymentRequest: {
        id: testPaymentRequestId,
        moveTaskOrderID: testMTOID,
        serviceItems: [
          {
            id: '1',
            mtoServiceItemID: 'a',
            priceCents: 12399,
          },
          {
            id: '2',
            mtoServiceItemID: 'b',
            priceCents: 45600,
          },
          {
            id: '3',
            mtoServiceItemID: 'c',
            priceCents: 12312,
          },
        ],
      },
      mtoServiceItems: [
        {
          id: 'a',
          mtoShipmentID: 'a1',
          reServiceName: 'Test Service Item',
          status: 'SUBMITTED',
          createdAt: '2020-01-01T00:09:00.999Z',
        },
        {
          id: 'b',
          mtoShipmentID: 'b2',
          reServiceName: 'Test Service Item 2',
          status: 'APPROVED',
          createdAt: '2020-01-01T00:09:00.999Z',
        },
        {
          id: 'c',
          mtoShipmentID: 'a1',
          reServiceName: 'Test Service Item 3',
          status: 'REJECTED',
          createdAt: '2020-01-01T00:09:00.999Z',
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

    it('maps the service item card data into the expected format and passes it into the ReviewServiceItems component', () => {
      const reviewServiceItems = wrapper.find('ReviewServiceItems');
      const expectedServiceItemCards = [
        {
          id: 'a',
          shipmentId: 'a1',
          shipmentType: 'HHG',
          serviceItemName: 'Test Service Item',
          amount: 123.99,
          createdAt: '2020-01-01T00:09:00.999Z',
          status: 'SUBMITTED',
        },
        {
          id: 'b',
          shipmentId: 'b2',
          shipmentType: 'NTS',
          serviceItemName: 'Test Service Item 2',
          amount: 456.0,
          createdAt: '2020-01-01T00:09:00.999Z',
          status: 'APPROVED',
        },
        {
          id: 'c',
          shipmentId: 'a1',
          shipmentType: 'HHG',
          serviceItemName: 'Test Service Item 3',
          amount: 123.12,
          createdAt: '2020-01-01T00:09:00.999Z',
          status: 'REJECTED',
        },
      ];

      expect(reviewServiceItems.prop('serviceItemCards')).toEqual(expectedServiceItemCards);
    });
  });
});
