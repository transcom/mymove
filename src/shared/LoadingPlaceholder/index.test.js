import React from 'react';
import ReactDOM from 'react-dom';
import { shallow } from 'enzyme';
import LoadingPlaceholder from '.';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<LoadingPlaceholder />, div);
});
