import React from 'react';
import ReactDOM from 'react-dom';
import { shallow } from 'enzyme';
import SignIn from './SignIn';

describe('SignIn tests', () => {
  let wrapper;
  it('renders without crashing', () => {
    const div = document.createElement('div');
    wrapper = shallow(<SignIn />, div);
    expect(wrapper.find('.usa-grid').length).toEqual(1);
  });
});
