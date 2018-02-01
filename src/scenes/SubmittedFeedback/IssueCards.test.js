import React from 'react';
import ReactDOM from 'react-dom';
import { shallow } from 'enzyme';
import IssueCards from './IssueCards';

describe('Null Issues on IssueCards', () => {
  const issues = null;

  it('renders without crashing', () => {
    const div = document.createElement('div');
    ReactDOM.render(<IssueCards issues={issues} />, div);
  });
});

describe('Empty Issues on IssueCards', () => {
  const issues = [];

  it('renders without crashing', () => {
    const div = document.createElement('div');
    ReactDOM.render(<IssueCards issues={issues} />, div);
  });
});

describe('Issues on IssueCards', () => {
  const issues = ['Too few dogs.'];

  it('renders without crashing', () => {
    const div = document.createElement('div');
    ReactDOM.render(<IssueCards issues={issues} />, div);
  });
});
