import React from 'react';
import { shallow } from 'enzyme';

import MaintenancePage from './MaintenancePage';

describe('Maintenance Page', () => {
  it(' renders the maintenance page text', () => {
    const wrapper = shallow(<MaintenancePage />);
    expect(wrapper.find('div')).toBeDefined();
    expect(wrapper.text()).toEqual(
      '<CUIHeader /><MilMoveHeader />System MaintenanceThis system is currently undergoing maintenance. Please check back later.',
    );
  });
});
