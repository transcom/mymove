import React from 'react';
import { mount } from 'enzyme';
import { ReactQueryCacheProvider, makeQueryCache } from 'react-query';

import PaymentRequestIndex from './paymentRequestIndex';

import { MockProviders } from 'testUtils';

const mockGetPaymentRequestListError = jest.fn(() => Promise.reject('API error'));

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

  // eslint-disable-next-line no-only-tests/no-only-tests
  describe.skip('error state', () => {
    const cache = makeQueryCache();
    it('shows an error', async () => {
      await cache.prefetchQuery('paymentRequestList', mockGetPaymentRequestListError, { retry: false });

      const wrapper = mount(
        <ReactQueryCacheProvider queryCache={cache}>
          <MockProviders initialEntries={['/']}>
            <PaymentRequestIndex />
          </MockProviders>
        </ReactQueryCacheProvider>,
      );

      expect(wrapper.find('SomethingWentWrong').exists()).toBe(true);
    });
  });

  describe('with data loaded', () => {
    const cache = makeQueryCache();

    it('renders without errors', async () => {
      await cache.prefetchQuery('paymentRequestList', mockGetPaymentRequestListSuccess);

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
