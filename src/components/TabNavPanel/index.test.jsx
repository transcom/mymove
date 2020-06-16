import React from 'react';
import { shallow } from 'enzyme';
import { TabPanel } from 'react-tabs';

import TabNavPanel from '.';

describe('TabNavPanel', () => {
  it('should render the tab navigation', () => {
    const wrapper = shallow(<TabNavPanel>Body Of Tab 1</TabNavPanel>);
    expect(wrapper.find(TabPanel).first().props().children).toBe('Body Of Tab 1');
  });
});
