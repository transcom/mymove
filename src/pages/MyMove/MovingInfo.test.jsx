import React from 'react';
import { mount } from 'enzyme';

import { MovingInfo } from './MovingInfo';

describe('MovingInfo component', () => {
  const testProps = {
    entitlementWeight: 7000,
    fetchLatestOrders: jest.fn(),
    serviceMemberId: '1234567890',
  };

  // eslint-disable-next-line react/jsx-props-no-spreading
  const wrapper = mount(<MovingInfo {...testProps} />);
  it('renders', () => {
    expect(wrapper.find('MovingInfo').length).toBe(1);
    expect(wrapper.text()).toContain('Tips for planning your shipments');
    expect(wrapper.find('[data-testid="shipmentsHeader"]').length).toBe(1);
    expect(wrapper.text()).toContain('7,000 lbs');
    expect(wrapper.find('[data-testid="shipmentsSubHeader"]').length).toBe(4);
  });
});

describe('MovingInfo when entitlement weight is 0', () => {
  const testProps = {
    entitlementWeight: 0,
    fetchLatestOrders: jest.fn(),
    serviceMemberId: '1234567890',
  };

  // eslint-disable-next-line react/jsx-props-no-spreading
  const wrapper = mount(<MovingInfo {...testProps} />);
  it('renders with no errors when entitlement weight is 0', () => {
    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('MovingInfo').length).toBe(1);
    expect(wrapper.text()).toContain('Tips for planning your shipments');
    expect(wrapper.find('[data-testid="shipmentsHeader"]').length).toBe(1);
    expect(wrapper.find('[data-testid="shipmentsAlert"]').length).toBe(0);
    expect(wrapper.find('[data-testid="shipmentsSubHeader"]').length).toBe(4);
  });
});
