import React from 'react';
import { shallow } from 'enzyme';
import SignIn from './SignIn';

describe('SignIn tests', () => {
  it('renders without crashing', () => {
    const div = document.createElement('div');
    shallow(<SignIn />, div);
  });
});
