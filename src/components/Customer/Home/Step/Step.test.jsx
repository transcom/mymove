/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import Step from '.';

const defaultProps = {
  actionBtnDisabled: false,
  actionBtnLabel: '',
  children: null,
  complete: false,
  completedHeaderText: '',
  containerClassName: '',
  editBtnDisabled: false,
  editBtnLabel: '',
  headerText: '',
  onActionBtnClick: () => {},
  onEditBtnClick: () => {},
  secondaryBtn: false,
  step: '',
};

function mountStep(props = defaultProps) {
  return mount(<Step {...props}>{props.children}</Step>);
}
describe('Step component', () => {
  it('renders Step with uncompleted step and header', () => {
    const headerText = 'This is header text';
    const step = '1';
    const props = {
      headerText,
      step,
    };
    const wrapper = mountStep(props);
    expect(wrapper.find('NumberCircle').length).toBe(1);
    expect(wrapper.find('.number-circle').text()).toBe('1');
    expect(wrapper.find('strong').text()).toBe(headerText);
  });

  it('renders a completed step', () => {
    const headerText = 'This is header text';
    const step = '1';
    const completedHeaderText = 'This is completed header text';
    const props = {
      headerText,
      complete: true,
      completedHeaderText,
      step,
    };
    const wrapper = mountStep(props);
    expect(wrapper.find('NumberCircle').length).toBe(0);
    expect(wrapper.find('svg').length).toBe(1);
  });

  it('should render a description', () => {
    const headerText = 'This is header text';
    const step = '1';
    const description = 'this is a description';
    const props = {
      headerText,
      step,
      children: <p>{description}</p>,
    };
    const wrapper = mountStep(props);
    expect(wrapper.find('p').text()).toBe(description);
  });

  it('should render children', () => {
    const headerText = 'This is header text';
    const step = '1';
    const children = <div data-testid="children">Hi I am child</div>;
    const props = {
      headerText,
      step,
      children,
    };
    const wrapper = mountStep(props);
    expect(wrapper.find('[data-testid="children"]').text()).toBe('Hi I am child');
  });

  it('renders Step with a call to action', () => {
    const actionBtnLabel = 'Action btn label';
    const editBtnLabel = 'Edit';
    const headerText = 'This is header text';
    const onActionBtnClick = jest.fn();
    const onEditBtnClick = jest.fn();
    const step = '1';
    const props = {
      actionBtnLabel,
      editBtnLabel,
      headerText,
      onActionBtnClick,
      onEditBtnClick,
      step,
    };
    const wrapper = mountStep(props);
    // const EditButton = wrapper.find('[data-testid="button"]').at(0);
    const ActionButton = wrapper.find('[data-testid="button"]').at(0);
    // expect(EditButton.text()).toBe(editBtnLabel);
    // expect(onEditBtnClick.mock.calls.length).toBe(0);
    // EditButton.simulate('click');
    // expect(onEditBtnClick.mock.calls.length).toBe(1);
    expect(ActionButton.text()).toBe(actionBtnLabel);
    ActionButton.simulate('click');
    expect(onActionBtnClick.mock.calls.length).toBe(1);
  });

  it('should not call handlers if disabled', () => {
    const actionBtnLabel = 'Action btn label';
    const editBtnLabel = 'Edit';
    const headerText = 'This is header text';
    const onActionBtnClick = jest.fn();
    const onEditBtnClick = jest.fn();
    const step = '1';
    const props = {
      actionBtnDisabled: true,
      actionBtnLabel,
      editBtnDisabled: true,
      editBtnLabel,
      headerText,
      onActionBtnClick,
      onEditBtnClick,
      step,
    };
    const wrapper = mountStep(props);
    // const EditButton = wrapper.find('[data-testid="button"]').at(0);
    const ActionButton = wrapper.find('[data-testid="button"]').at(0);
    // expect(onEditBtnClick.mock.calls.length).toBe(0);
    // EditButton.simulate('click');
    expect(onEditBtnClick.mock.calls.length).toBe(0);
    ActionButton.simulate('click');
    expect(onActionBtnClick.mock.calls.length).toBe(0);
  });
});
