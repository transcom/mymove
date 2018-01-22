import React from 'react';
import ReactDOM from 'react-dom';
import { shallow } from 'enzyme';
import Feedback from './Feedback';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<Feedback />, div);
});

// instance() is needed to instantiate the shallow render. The handleChange method needs to be passed an object, rather than just the key's value, to work. More here: https://stackoverflow.com/questions/45236277/how-to-test-react-component-class-methods
test('it should change state when provided a key and value', () => {
  const wrapper = shallow(<Feedback />);
  wrapper.instance().handleChange({ target: { value: 'test text' } });
  expect(wrapper.instance().state.value).toEqual('test text');
});
