import React from 'react';
import { shallow } from 'enzyme';
import { UserGreeting } from './UserGreeting';

describe('UserGreeting tests', () => {
  describe('User is not authenticated', () => {
    it('should not render if the user is not logged in', () => {
      const wrapper = shallow(<UserGreeting isLoggedIn={false} email="" />);
      expect(wrapper.children().length).toBe(0);
    });
  });
  describe('User is authenticated', () => {
    it('should render the Welcome message if the user provides a firstName', () => {
      const wrapper = shallow(<UserGreeting isLoggedIn={true} firstName="Kevin" email="kevin@test.com" />);
      expect(wrapper.find('strong').text()).toBe('Welcome, Kevin');
    });
    it('should render the email if the user does not have a first name yet', () => {
      const wrapper = shallow(<UserGreeting isLoggedIn={true} email="kevin@test.com" />);
      expect(wrapper.find('strong').text()).toBe('kevin@test.com');
    });
  });
});
