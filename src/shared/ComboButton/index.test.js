import React from 'react';
import { shallow } from 'enzyme';
import ComboButton from './index';

describe('ComboButton tests', () => {
  const renderComboButton = ({ buttonText = '', isDisabled = false, toolTipText = undefined }) =>
    shallow(<ComboButton buttonText={buttonText} isDisabled={isDisabled} toolTipText={toolTipText} />);

  describe('button text', () => {
    it('renders button using buttonText', () => {
      const comboButton = renderComboButton({ buttonText: 'Text' });
      const button = comboButton.find('button');
      expect(button.render().text()).toEqual('Text');
    });
  });

  describe('disabled state', () => {
    it('renders a disabled button when isDisabled is true', () => {
      const comboButton = renderComboButton({ isDisabled: true });
      const button = comboButton.find('button');
      expect(button.props().disabled).toBe(true);
    });

    it('renders an enabled button when isDisabled is false', () => {
      const comboButton = renderComboButton({});
      const button = comboButton.find('button');
      expect(button.props().disabled).toBe(false);
    });
  });

  describe('tool tip', () => {
    it('renders the tool tip when passed tool tip text', () => {
      const text = 'tooltipText';
      const comboButton = renderComboButton({ toolTipText: text });
      const tooltipText = comboButton.find('.tooltiptext');

      expect(tooltipText.text()).toBe(text);
    });

    it('does not render the tool tip when no text is passed', () => {
      const comboButton = renderComboButton({ toolTipText: null });
      const tooltipText = comboButton.find('.tooltiptext');

      expect(tooltipText.exists()).toBe(false);
    });
  });
});
