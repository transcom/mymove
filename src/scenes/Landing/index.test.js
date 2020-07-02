import 'raf/polyfill';
import React from 'react';

import { Provider } from 'react-redux';
import moment from 'moment';
import configureStore from 'redux-mock-store';
import { shallow } from 'enzyme';

import { Landing } from '.';
import { MoveSummary } from './MoveSummary';
import PpmAlert from './PpmAlert';
import { MOVE_TYPES } from 'shared/constants';

describe('HomePage tests', () => {
  let wrapper;
  const mockStore = configureStore();
  let store;
  describe('when not loggedIn', () => {
    it('renders without crashing', () => {
      const div = document.createElement('div');
      wrapper = shallow(<Landing isLoggedIn={false} />, div);
      expect(wrapper.find('.grid-container').length).toEqual(1);
    });
  });
  describe('When loggedIn', () => {
    let service_member = { id: 'foo' };
    it('renders without crashing', () => {
      const div = document.createElement('div');
      wrapper = shallow(<Landing isLoggedIn={true} />, div);
      expect(wrapper.find('.grid-container').length).toEqual(1);
    });
    describe('When the user has never logged in before', () => {
      it('redirects to enter profile page', () => {
        const mockPush = jest.fn();
        store = mockStore({});
        const wrapper = shallow(
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
        const landing = wrapper.find(Landing).dive();
        const resumeMoveFn = jest.spyOn(landing.instance(), 'resumeMove');
        landing.setProps({ createdServiceMemberIsLoading: true });
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
        expect(wrapper.find('.grid-container').length).toEqual(1);
        expect(wrapper.find(MoveSummary).length).toEqual(0);
      });
    });

    describe('When orders have been entered but the move is not complete', () => {
      const futureFortNight = moment().add(14, 'day');
      const orders = {
        orders_type: 'foo',
        issue_date: '2019-01-01',
        report_by_date: '2019-02-01',
        new_duty_station: { id: 'something' },
      };
      const ppmObj = {
        original_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };

      describe('When a ppm only move is submitted', () => {
        it('renders the ppm only alert', () => {
          const moveObj = { selected_move_type: MOVE_TYPES.PPM, status: 'SUBMITTED', moveSubmitSuccess: true };
          wrapper = shallow(
            <Landing
              move={moveObj}
              moveSubmitSuccess={moveObj.moveSubmitSuccess}
              serviceMember={service_member}
              ppm={ppmObj}
              orders={orders}
              isLoggedIn={true}
              loggedInUserSuccess={true}
              isProfileComplete={true}
            />,
          );
          const ppmAlert = wrapper.find(PpmAlert).shallow();

          expect(ppmAlert.length).toEqual(1);
          expect(ppmAlert.props().heading).toEqual('Congrats - your move is submitted!');
        });
      });
    });
  });
});
