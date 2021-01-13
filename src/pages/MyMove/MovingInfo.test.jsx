import React from 'react';
import { mount } from 'enzyme';

import { MovingInfo } from './MovingInfo';

const testProps = {
  entitlementWeight: 7000,
  fetchLatestOrders: jest.fn(),
  serviceMemberId: '1234567890',
};

// eslint-disable-next-line react/jsx-props-no-spreading
const wrapper = mount(<MovingInfo {...testProps} />);

describe('MovingInfo component', () => {
  it('renders', () => {
    expect(wrapper.find('MovingInfo').length).toBe(1);
    expect(wrapper.text()).toContain('Tips for planning your shipments');
    expect(wrapper.text()).toContain('7,000 lbs');
    expect(wrapper.find('[data-testid="shipmentsHeader"]').length).toBe(1);
    expect(wrapper.find('[data-testid="shipmentsSubHeader"]').length).toBe(4);
  });
});
