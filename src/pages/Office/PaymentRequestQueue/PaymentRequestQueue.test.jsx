import React from 'react';
import { mount } from 'enzyme';

import PaymentRequestQueue from './PaymentRequestQueue';

import { MockProviders } from 'testUtils';

describe('PaymentRequestQueue', () => {
  const wrapper = mount(
    <MockProviders initialEntries={['invoicing/queue']}>
      <PaymentRequestQueue />
    </MockProviders>,
  );

  it('should render the h1', () => {
    expect(wrapper.find('h1').text()).toBe('Payment requests (0)');
  });

  it('should render the table', () => {
    expect(wrapper.find('Table').exists()).toBe(true);
  });
});
