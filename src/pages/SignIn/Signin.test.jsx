import React from 'react';
import { shallow, mount } from 'enzyme';

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
});
