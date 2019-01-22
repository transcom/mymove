import 'raf/polyfill';
import React from 'react';

import { Provider } from 'react-redux';
import moment from 'moment';
import configureStore from 'redux-mock-store';
import { shallow, mount } from 'enzyme';

import { Landing } from '.';
import { MoveSummary } from './MoveSummary';

describe('HomePage tests', () => {
  let wrapper;
  const mockStore = configureStore();
  let store;
  describe('when not loggedIn', () => {
    it('renders without crashing', () => {
      const div = document.createElement('div');
      wrapper = shallow(<Landing isLoggedIn={false} />, div);
      expect(wrapper.find('.usa-grid').length).toEqual(1);
    });
  });
  describe('When loggedIn', () => {
    let service_member = { id: 'foo' };
    it('renders without crashing', () => {
      const div = document.createElement('div');
      wrapper = shallow(<Landing isLoggedIn={true} />, div);
      expect(wrapper.find('.usa-grid').length).toEqual(1);
    });
    describe('When the user has never logged in before', () => {
      it('redirects to enter profile page', () => {
        const mockPush = jest.fn();
        store = mockStore({});
        const wrapper = mount(
          <Provider store={store}>
            <Landing
              isLoggedIn={true}
              serviceMember={service_member}
              createdServiceMemberIsLoading={false}
              loggedInUserSuccess={true}
              isProfileComplete={false}
              push={mockPush}
              reduxState={{}}
            />
          </Provider>,
        );
        const resumeMoveFn = jest.spyOn(wrapper.children().instance(), 'resumeMove');
        wrapper.setProps({ createdServiceMemberIsLoading: true });
        expect(resumeMoveFn).toHaveBeenCalledTimes(1);
      });
    });
    describe('When the user profile has started but is not complete', () => {
      it('MoveSummary does not render', () => {
        const div = document.createElement('div');
        wrapper = shallow(
          <Landing
            serviceMember={service_member}
            isLoggedIn={true}
            loggedInUserSuccess={true}
            isProfileComplete={false}
          />,
          div,
        );
        expect(wrapper.find('.usa-grid').length).toEqual(1);
        expect(wrapper.find(MoveSummary).length).toEqual(0);
      });
    });
    describe('When the user profile is complete but orders have not been entered', () => {
      const moveObj = { id: 'foo', selected_move_type: 'PPM' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };
      it('Renders resume move text', () => {
        const div = document.createElement('div');
        wrapper = shallow(
          <Landing
            serviceMember={service_member}
            ppm={ppmObj}
            move={moveObj}
            isLoggedIn={true}
            loggedInUserSuccess={true}
            isProfileComplete={true}
          />,
          div,
        );
        const moveSummary = shallow(
          wrapper
            .find(MoveSummary)
            .dive()
            .props()
            .children(),
        );

        expect(
          moveSummary
            .find('.status-component')
            .dive()
            .find('.step')
            .find('div.title')
            .first()
            .html(),
        ).toEqual('<div class="title">Next Step: Finish setting up your move</div>');
      });
    });
    describe('When orders have been entered but the move is not complete', () => {
      const moveObj = { id: 'foo', selected_move_type: 'PPM', status: 'DRAFT' };
      const futureFortNight = moment().add(14, 'day');
      const orders = {
        orders_type: 'foo',
        issue_date: '2019-01-01',
        report_by_date: '2019-02-01',
        new_duty_station: { id: 'something' },
      };
      const ppmObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };
      it('Renders resume move text', () => {
        const div = document.createElement('div');
        wrapper = shallow(
          <Landing
            serviceMember={service_member}
            ppm={ppmObj}
            move={moveObj}
            orders={orders}
            isLoggedIn={true}
            loggedInUserSuccess={true}
            isProfileComplete={true}
          />,
          div,
        );
        const moveSummary = shallow(
          wrapper
            .find(MoveSummary)
            .dive()
            .props()
            .children(),
        );
        expect(
          moveSummary
            .find('.status-component')
            .dive()
            .find('.step')
            .find('div.title')
            .first()
            .html(),
        ).toEqual('<div class="title">Next Step: Finish setting up your move</div>');
      });
    });
  });
});
