import 'raf/polyfill';
import React from 'react';

import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { shallow } from 'enzyme';

import { PpmLanding } from '.';
import { PpmSummary } from './PpmSummary';

describe('PPM landing page tests', () => {
  let wrapper;
  const mockStore = configureStore();
  let store;

  const minProps = {
    showLoggedInUser: () => {},
    context: {
      flags: {
        hhgFlow: false,
      },
    },
  };

  describe('when not loggedIn', () => {
    it('renders without crashing', () => {
      const div = document.createElement('div');
      wrapper = shallow(<PpmLanding isLoggedIn={false} {...minProps} />, div);
      expect(wrapper.find('.grid-container').length).toEqual(1);
    });
  });

  describe('When loggedIn', () => {
    let service_member = { id: 'foo' };
    it('renders without crashing', () => {
      const div = document.createElement('div');
      wrapper = shallow(<PpmLanding isLoggedIn={true} {...minProps} />, div);
      expect(wrapper.find('.grid-container').length).toEqual(1);
    });

    describe('When the user has never logged in before', () => {
      it('redirects to enter profile page', () => {
        const mockPush = jest.fn();
        store = mockStore({});
        const wrapper = shallow(
          <Provider store={store}>
            <PpmLanding
              isLoggedIn={true}
              serviceMember={undefined}
              createdServiceMemberIsLoading={false}
              loggedInUserSuccess={true}
              isProfileComplete={false}
              push={mockPush}
              reduxState={{}}
              {...minProps}
            />
          </Provider>,
        );
        const landing = wrapper.find(PpmLanding).dive();
        const resumeMoveFn = jest.spyOn(landing.instance(), 'resumeMove');
        landing.setProps({ serviceMember: service_member });
        expect(resumeMoveFn).toHaveBeenCalledTimes(1);
      });
    });

    describe('When the user profile has started but is not complete', () => {
      it('PpmSummary does not render', () => {
        const div = document.createElement('div');
        wrapper = shallow(
          <PpmLanding
            serviceMember={service_member}
            isLoggedIn={true}
            loggedInUserSuccess={true}
            isProfileComplete={false}
            push={jest.fn()}
            {...minProps}
          />,
          div,
        );
        expect(wrapper.find('.grid-container').length).toEqual(1);
        expect(wrapper.find(PpmSummary).length).toEqual(0);
      });
    });
  });
});
