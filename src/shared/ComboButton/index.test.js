import React from 'react';
import { mount } from 'enzyme';
import ComboButton from './index';

describe('ComboButton tests', () => {
  const renderComboButton = ({ buttonText = '', disabled = false }) =>
    mount(
      <ComboButton buttonText={buttonText} disabled={disabled}>
        <div className="dropdown">dropDownText</div>
      </ComboButton>,
    );

  describe('when the button is disabled', () => {
    const buttonProps = { buttonText: 'buttonText', disabled: true };
    const disabledComboButton = renderComboButton(buttonProps);

    describe('button', () => {
      const button = disabledComboButton.find('button');

      it('renders using value from buttonText', () => {
        expect(button.render().text()).toEqual(buttonProps.buttonText);
      });

      it('renders in a disabled state', () => {
        expect(button.props().disabled).toBe(true);
      });
    });
  });

  describe('when the button is enabled', () => {
    const buttonProps = { disabled: false, buttonText: 'buttonText' };
    const defaultEnabledComboButton = renderComboButton(buttonProps);

    describe('button', () => {
      it('renders in enabled state', () => {
        const button = defaultEnabledComboButton.find('button');

        expect(button.props().disabled).toBe(false);
      });
    });

    describe('dropdown menu', () => {
      it('is not displayed', () => {
        const dropDown = defaultEnabledComboButton.find('.dropdown');

        expect(dropDown.exists()).toBe(false);
      });

      it('is displayed on click', () => {
        const enabledComboButton = renderComboButton(buttonProps);
        enabledComboButton.find('button').simulate('click');
        const dropDown = enabledComboButton.find('.dropdown');

        expect(dropDown.exists()).toBe(true);
      });

      it('disappears on second click', () => {
        const enabledComboButton = renderComboButton(buttonProps);
        enabledComboButton.find('button').simulate('click');
        const dropDownAfterFirstClick = enabledComboButton.find('.dropdown');

        expect(dropDownAfterFirstClick.exists()).toBe(true);
        enabledComboButton.find('button').simulate('click');
        const dropDownAfterSecondClick = enabledComboButton.find('.dropdown');
        expect(dropDownAfterSecondClick.exists()).toBe(false);
      });

      it('state.displayDropDown is false after click outside', () => {
        const newButtonProps = { toolTipText: 'toolTipText', disabled: false, buttonText: 'buttonText' };
        const enabledComboButton = renderComboButton(newButtonProps);
        enabledComboButton.setState({ displayDropDown: true });
        const enabledComboButtonInstance = enabledComboButton.instance();
        enabledComboButtonInstance.handleClickOutside({});

        expect(enabledComboButton.state().displayDropDown).toBe(false);
      });

      it('state.displayDropDown is toggled on click', function() {
        const newButtonProps = { toolTipText: 'toolTipText', disabled: false, buttonText: 'buttonText' };
        const enabledComboButton = renderComboButton(newButtonProps);
        enabledComboButton.setState({ displayDropDown: true });
        enabledComboButton.find('button').simulate('click');

        expect(enabledComboButton.state().displayDropDown).toBe(false);
      });
    });
  });
});
