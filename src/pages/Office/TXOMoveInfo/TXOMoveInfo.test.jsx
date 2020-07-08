import React from 'react';
import { mount } from 'enzyme';

import TXOMoveInfo from './TXOMoveInfo';

import { MockProviders } from 'testUtils';

const testMoveId = '10000';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveOrderId: '10000' }),
}));

describe('TXO Move Info Container', () => {
  it('should render the move tab container', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveId}/details`]}>
        <TXOMoveInfo />
      </MockProviders>,
    );

    expect(wrapper.find('header.nav-header').exists()).toBe(true);
    expect(wrapper.find('nav.tabNav').exists()).toBe(true);
    expect(wrapper.find('li.tabItem').length).toEqual(4);

    expect(wrapper.find('span.tab-title').at(0).text()).toContain('Move details');
    expect(wrapper.find('span.tab-title').at(1).text()).toContain('Move task order');
    expect(wrapper.find('span.tab-title').at(2).text()).toContain('Payment requests');
    expect(wrapper.find('span.tab-title').at(3).text()).toContain('History');

    expect(wrapper.find('li.tabItem a').at(0).prop('href')).toEqual(`/moves/${testMoveId}/details`);
    expect(wrapper.find('li.tabItem a').at(1).prop('href')).toEqual(`/moves/${testMoveId}/mto`);
    expect(wrapper.find('li.tabItem a').at(2).prop('href')).toEqual(`/moves/${testMoveId}/payment-requests`);
    expect(wrapper.find('li.tabItem a').at(3).prop('href')).toEqual(`/moves/${testMoveId}/history`);
  });
});
