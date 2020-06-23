import React from 'react';
import { Provider } from 'react-redux';
import { MemoryRouter } from 'react-router-dom';
import { mount } from 'enzyme';

import { store } from '../shared/store';

import TXOMoveInfo from './TXOMoveInfo';

describe('TXO Move Info Container', () => {
  it('should render the move tab container', () => {
    const wrapper = mount(
      <Provider store={store}>
        <MemoryRouter initialEntries={['/moves/10000/details']}>
          <TXOMoveInfo too tag="main" />
        </MemoryRouter>
      </Provider>,
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
