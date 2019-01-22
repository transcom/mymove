import React from 'react';
import { mount } from 'enzyme';
import { SubmittedFeedback } from '.';

const loadIssues = () => {};

describe('No Issues and Errors', () => {
  let wrapper;

  beforeEach(() => {
    const issues = null;
    const hasError = true;
    wrapper = mount(<SubmittedFeedback hasError={hasError} issues={issues} loadIssues={loadIssues} />);
  });

  it('renders an alert', () => {
    expect(wrapper.find('Alert').length).toBe(1);
  });

  it('does not render issue cards', () => {
    expect(wrapper.find('IssueCards').length).toBe(0);
  });
});

describe('Has issues', () => {
  let wrapper;

  beforeEach(() => {
    const issues = [{ id: '10', description: 'Too few dogs.' }, { id: '20', description: 'Too little barking.' }];
    const hasError = false;
    wrapper = mount(<SubmittedFeedback hasError={hasError} issues={issues} loadIssues={loadIssues} />);
  });

  it('renders without an alert', () => {
    expect(wrapper.find('Alert').length).toBe(0);
  });

  it('renders issue cards without crashing', () => {
    expect(wrapper.find('IssueCards').length).toBe(1);
  });
});
