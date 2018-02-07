import React from 'react';
import { shallow } from 'enzyme';
import IssueCards from './IssueCards';

describe('Null Issues on IssueCards', () => {
  let wrapper;
  const issues = null;

  beforeEach(() => {
    wrapper = shallow(<IssueCards issues={issues} />);
  });

  it('renders without crashing', () => {
    expect(wrapper.find('IssueCards').toExist);
  });
});

describe('Empty Issues on IssueCards', () => {
  let wrapper;
  const issues = [];

  beforeEach(() => {
    wrapper = shallow(<IssueCards issues={issues} />);
  });

  it('renders without crashing', () => {
    expect(wrapper.find('IssueCards').toExist);
  });
});

describe('Issues on IssueCards', () => {
  let wrapper;
  const issues = [{ id: '13', description: 'Too few dogs.' }];

  beforeEach(() => {
    wrapper = shallow(<IssueCards issues={issues} />);
  });

  it('renders without crashing', () => {
    expect(wrapper.find('IssueCards').toExist);
  });
});
