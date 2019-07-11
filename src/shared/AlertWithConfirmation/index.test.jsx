import React from 'react';
import { shallow } from 'enzyme';
import Alert from '.';

describe('basic alert component', () => {
  const text = 'some text';
  const heading = 'a heading';
  const wrapper = shallow(<Alert heading={heading}>{text}</Alert>);
  it('should render children and heading', () => {
    expect(wrapper.find('.usa-alert-heading').text()).toBe(heading);
    expect(wrapper.find('.usa-alert-text').text()).toBe(text);
  });
  it('should not have a close button', () => {
    expect(wrapper.find('.icon.remove-icon')).toHaveLength(0);
  });
  it('should not display a spinner', () => {
    expect(wrapper.find('.heading--icon')).toHaveLength(0);
  });
  describe('loading alert', () => {
    const wrapper = shallow(
      <Alert heading={heading} type="loading">
        {text}
      </Alert>,
    );
    it('should display a spinner', () => {
      expect(wrapper.find('.heading--icon')).toHaveLength(1);
    });
  });
  describe('close button', () => {
    const mockOnRemove = jest.fn();
    const wrapper = shallow(
      <Alert heading={heading} onRemove={mockOnRemove}>
        {text}
      </Alert>,
    );
    it('should render a close button', () => {
      expect(wrapper.find('.icon.remove-icon')).toHaveLength(1);
      wrapper.find('.icon.remove-icon').simulate('click');
      expect(mockOnRemove).toHaveBeenCalled();
    });
  });
});
