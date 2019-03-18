import React from 'react';
import { shallow, mount } from 'enzyme';

import SignIn from './SignIn';
import Alert from 'shared/Alert';

describe('SignIn tests', () => {
  it('renders without crashing', () => {
    const div = document.createElement('div');
    shallow(<SignIn />, div);
  });
  it('renders errors', () => {
    const context = { siteName: 'TestMove' };
    const location = { search: '?error=SOME_ERROR' };
    const wrapper = mount(<SignIn location={location} />, { context });
    expect(wrapper.find(Alert).text()).toContain('SOME_ERROR');
  });
});
