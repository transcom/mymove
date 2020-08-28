import React from 'react';
import { Provider } from 'react-redux';
import store from 'shared/store';
import { mount } from 'enzyme';
import AlertWithDeleteConfirmation from '.';

describe('basic alert with delete confirmation component', () => {
  const text = 'some text';
  const heading = 'a heading';
  const wrapper = mount(
    <Provider store={store}>
      <AlertWithDeleteConfirmation heading={heading} message={text} />
    </Provider>,
  );
  it('should render children and heading', () => {
    expect(wrapper.find('.delete-alert-heading').text()).toBe(heading);
    expect(wrapper.find('.usa-alert__text').text()).toBe(text);
  });
  it('should not have a close button', () => {
    expect(wrapper.find('.icon.remove-icon')).toHaveLength(0);
  });
  it('should not display a spinner', () => {
    expect(wrapper.find('.heading--icon')).toHaveLength(0);
  });
  describe('cancel confirmation', () => {
    const mockCancelActionHandler = jest.fn();
    const wrapper = mount(
      <Provider store={store}>
        <AlertWithDeleteConfirmation heading={heading} message={text} cancelActionHandler={mockCancelActionHandler} />
      </Provider>,
    );
    it('should render cancel button', () => {
      expect(wrapper.find('.usa-button--secondary')).toHaveLength(1);
      wrapper.find('.usa-button--secondary').simulate('click');
      expect(mockCancelActionHandler).toHaveBeenCalled();
    });
  });
  describe('delete confirmation', () => {
    const mockDeleteActionHandler = jest.fn();
    const wrapper = mount(
      <Provider store={store}>
        <AlertWithDeleteConfirmation heading={heading} message={text} deleteActionHandler={mockDeleteActionHandler} />
      </Provider>,
    );
    it('should render delete and ok buttons', () => {
      expect(wrapper.find('.usa-button')).toHaveLength(2);
      wrapper.find('.usa-button').at(0).simulate('click');
      expect(mockDeleteActionHandler).toHaveBeenCalled();
    });
  });
});
