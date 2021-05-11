import React from 'react';
import { mount, shallow } from 'enzyme';

import SignIn from './SignIn';

import Alert from 'shared/Alert';

describe('SignIn tests', () => {
  it('renders without crashing', () => {
    const div = document.createElement('div');
    shallow(<SignIn />, div);
  });
  it('does not render content of error parameter', () => {
    const context = { siteName: 'TestMove' };
    const location = { search: '?error=SOME_ERROR' };
    const wrapper = mount(<SignIn location={location} />, { context });
    expect(wrapper.find(Alert).text()).not.toContain('SOME_ERROR');
  });

  it('shows the EULA when the signin button is clicked and hides the EULA when cancel is clicked', () => {
    const context = { siteName: 'TestMove' };
    const location = { search: '' };
    const wrapper = mount(<SignIn location={location} />, { context });
    expect(wrapper.find('[data-testid="modal"]').length).toEqual(0);
    wrapper.find('button[data-testid="signin"]').simulate('click');
    expect(wrapper.find('[data-testid="modal"]').length).toEqual(1);
    const CancelButton = wrapper.find('button[aria-label="Cancel"]');
    CancelButton.simulate('click');
    expect(wrapper.find('[data-testid="modal"]').length).toEqual(0);
  });

  it('show logout message when hasLoggedOut state is true', () => {
    const context = { siteName: 'TestMove' };
    const location = { state: { hasLoggedOut: true } };
    const wrapper = mount(<SignIn location={location} />, { context });
    expect(wrapper.find(Alert).text()).toContain('You have signed out of MilMove');
  });

  it('does not show logout message when hasLoggedOut state is false', () => {
    const context = { siteName: 'TestMove' };
    const location = { state: { hasLoggedOut: false } };
    const wrapper = mount(<SignIn location={location} />, { context });
    expect(wrapper.find(Alert).length).toEqual(0);
  });

  it('show logout message when timedout state is true', () => {
    const context = { siteName: 'TestMove' };
    const location = { state: { timedout: true } };
    const wrapper = mount(<SignIn location={location} />, { context });
    expect(wrapper.find(Alert).text()).toContain('You have been logged out due to inactivity');
  });

  it('does not show logout message when timedout state is false', () => {
    const context = { siteName: 'TestMove' };
    const location = { state: { timedout: false } };
    const wrapper = mount(<SignIn location={location} />, { context });
    expect(wrapper.find(Alert).length).toEqual(0);
  });
});
