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

  describe('routing', () => {
    it('should handle the Move Details route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveId}/details`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      expect(wrapper.find('MoveDetails')).toHaveLength(1);
    });

    it('should redirect from move info root to the Move Details route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveId}`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('MoveDetails');
      expect(renderedRoute).toHaveLength(1);
    });

    it('should handle the Move Orders route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveId}/orders`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual(['/moves/:moveCode/allowances', '/moves/:moveCode/orders']);
    });

    it('should handle the Allowances route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveId}/allowances`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual(['/moves/:moveCode/allowances', '/moves/:moveCode/orders']);
    });

    it('should handle the Move Task Order route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveId}/mto`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/moves/:moveCode/mto');
    });

    it('should handle the Move Payment Requests route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveId}/payment-requests`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/moves/:locator/payment-requests');
    });

    it('should handle the Move History route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/moves/${testMoveId}/history`]}>
          <TXOMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('Route');
      expect(renderedRoute).toHaveLength(1);
      expect(renderedRoute.prop('path')).toEqual('/moves/:moveOrderId/history');
    });
  });
});
