import React from 'react';
import ReactDOM from 'react-dom';
import { shallow } from 'enzyme';
import { expect } from 'chai';
import Feedback from './Feedback';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<Feedback />, div);
});

const wrapper = shallow(<Feedback />);
expect(wrapper.instance().handleChange()).equals('something');
