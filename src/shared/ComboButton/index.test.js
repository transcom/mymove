import React from 'react';
import { mount } from 'enzyme';
import ComboButton from './index';

describe('ComboButton tests', () => {
  const renderComboButton = ({ buttonText = '', disabled = false, toolTipText = undefined }) =>
    mount(<ComboButton buttonText={buttonText} disabled={disabled} toolTipText={toolTipText} />);

  describe('when the button is disabled', () => {
    const buttonProps = { toolTipText: 'toolTipText', buttonText: 'buttonText', disabled: true };
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

    describe('tool tip', () => {
      it('renders with toolTipText', () => {
        const tooltipText = disabledComboButton.find('.tooltiptext');

        expect(tooltipText.text()).toBe(buttonProps.toolTipText);
      });

      it('does not render when toolTipText is null', () => {
        const comboButton = renderComboButton({ toolTipText: null });
        const tooltipText = comboButton.find('.tooltiptext');

        expect(tooltipText.exists()).toBe(false);
      });
    });
  });

  describe('when the button is enabled', () => {
    const buttonProps = { toolTipText: 'toolTipText', disabled: false, buttonText: 'buttonText' };
    const defaultEnabledComboButton = renderComboButton(buttonProps);

    describe('button', () => {
      it('renders in enabled state', () => {
        const button = defaultEnabledComboButton.find('button');

        expect(button.props().disabled).toBe(false);
      });
    });

    describe('tooltip', () => {
      it('does not render', () => {
        const tooltipText = defaultEnabledComboButton.find('.tooltiptext');

        expect(tooltipText.exists()).toBe(false);
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

      it('disappears on click outside of button', () => {
        const enabledComboButton = mount(
          <div className={'outside'}>
            <ComboButton {...buttonProps} />
          </div>,
        );

        enabledComboButton.find('button').simulate('click');
        const dropDownAfterFirstClick = enabledComboButton.find('.dropdown');
        expect(dropDownAfterFirstClick.exists()).toBe(true);
        enabledComboButton.find('.outside').simulate('click');
        const dropDownAfterSecondClick = enabledComboButton.find('.dropdown');
        expect(dropDownAfterSecondClick.exists()).toBe(false);
      });
    });
  });
});
