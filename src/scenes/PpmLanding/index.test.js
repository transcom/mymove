import React from 'react';
import { shallow } from 'enzyme';

import { PpmSummary } from './PpmSummary';

import { PpmLanding } from './index';

describe('PPM landing page tests', () => {
  let wrapper;

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
    const service_member = { id: 'foo' };
    it('renders without crashing', () => {
      const div = document.createElement('div');
      wrapper = shallow(<PpmLanding isLoggedIn {...minProps} />, div);
      expect(wrapper.find('.grid-container').length).toEqual(1);
    });

    describe('When the user profile has started but is not complete', () => {
      it('PpmSummary does not render', () => {
        const div = document.createElement('div');
        wrapper = shallow(
          <PpmLanding
            serviceMember={service_member}
            isLoggedIn
            loggedInUserSuccess
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
