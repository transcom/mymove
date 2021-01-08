import React from 'react';
import { mount } from 'enzyme';

import { MovingInfo } from './MovingInfo';

const wrapper = mount(<MovingInfo />);

describe('MovingInfo component', () => {
  it('renders', () => {
    expect(wrapper.find('MovingInfo').length).toBe(1);
    expect(wrapper.text()).toContain('Tips for planning your shipments');
    expect(wrapper.find('[data-testid="shipmentsHeader"]').length).toBe(1);
    expect(wrapper.find('[data-testid="shipmentsSubHeader"]').length).toBe(4);
  });
});
