import React from 'react';
import { mount } from 'enzyme';
import { ReactQueryCacheProvider, makeQueryCache } from 'react-query';

import PaymentRequestIndex from './PaymentRequestIndex';

import { MockProviders } from 'testUtils';
import { PAYMENT_REQUESTS } from 'constants/queryKeys';

const mockGetPaymentRequestListSuccess = jest.fn(() =>
  Promise.resolve({
    paymentRequests: {
      a1: {
        id: 'a1',
        moveTaskOrderID: '999',
        isFinal: false,
        rejectionReason: '',
        serviceItemIDs: [],
        status: 'PENDING',
      },
    },
  }),
);

describe('PaymentRequestIndex', () => {
  describe('loading state', () => {
    it('shows the loader', () => {
      const wrapper = mount(
        <MockProviders initialEntries={['/']}>
          <PaymentRequestIndex />
        </MockProviders>,
      );

      expect(wrapper.find('LoadingPlaceholder').exists()).toBe(true);
    });
  });

  describe('with data loaded', () => {
    const cache = makeQueryCache();

    it('renders without errors', async () => {
      await cache.prefetchQuery(PAYMENT_REQUESTS, mockGetPaymentRequestListSuccess);

      const wrapper = mount(
        <ReactQueryCacheProvider queryCache={cache}>
          <MockProviders initialEntries={['/']}>
            <PaymentRequestIndex />
          </MockProviders>
        </ReactQueryCacheProvider>,
      );

      expect(wrapper.find('[data-testid="PaymentRequestIndex"]').exists()).toBe(true);
    });
  });
});
