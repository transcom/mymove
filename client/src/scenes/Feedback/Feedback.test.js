import React from 'react';
import ReactDOM from 'react-dom';
import { shallow } from 'enzyme';
import { expect } from 'chai';
import Feedback from './Feedback';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<Feedback />, div);
});

// https://stackoverflow.com/questions/45236277/how-to-test-react-component-class-methods
test('it should change state when provided a key and value', () => {
  const wrapper = shallow(<Feedback />);
  expect(wrapper.instance().handleChange({ target: 'test text' }));
});
