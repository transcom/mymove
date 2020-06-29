import React from 'react';
import { mount } from 'enzyme';

import TXOMoveInfo from './TXOMoveInfo';

import { MockProviders } from 'testUtils';

describe('TXO Move Info Container', () => {
  it('should render the move tab container', () => {
    const wrapper = mount(
      <MockProviders initialEntries={['/moves/10000/details']}>
        <TXOMoveInfo too tag="main" />
      </MockProviders>,
    );

    expect(wrapper.find('header.nav-header').exists()).toBe(true);
    expect(wrapper.find('nav.tabNav').exists()).toBe(true);
    expect(wrapper.find('li.tabItem').length).toEqual(4);
    expect(wrapper.find('span.tab-title').at(0).text()).toContain('Move details');
    expect(wrapper.find('span.tab-title').at(1).text()).toContain('Move task order');
    expect(wrapper.find('span.tab-title').at(2).text()).toContain('Payment requests');
    expect(wrapper.find('span.tab-title').at(3).text()).toContain('History');
  });
});
