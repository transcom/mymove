import React from 'react';
import ReactDOM from 'react-dom';
import { shallow } from 'enzyme';
import FeedbackForm from './FeedbackForm';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<FeedbackForm />, div);
});
