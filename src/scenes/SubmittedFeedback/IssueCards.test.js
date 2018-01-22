import React from 'react';
import ReactDOM from 'react-dom';
import { shallow } from 'enzyme';
import IssueCards from './IssueCards';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<IssueCards issues={null} />, div);
});

//todo: test that variations of issues (null,empty,data) render as expected
