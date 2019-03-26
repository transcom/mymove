import React from 'react';
import { mount } from 'enzyme';
import ToolTip from './index';

describe('ToolTip tests', () => {
  const renderToolTip = ({ disabled = false, toolTipText = 'toolTipText', textStyle = 'tooltip-style' }) =>
    mount(<ToolTip disabled={disabled} toolTipText={toolTipText} textStyle={textStyle} />);

  describe('when ToolTip is disabled', () => {
    const toolTipProps = { toolTipText: 'toolTipText', buttonText: 'buttonText', disabled: true };
    const disabledToolTip = renderToolTip(toolTipProps);

    it('does not render toolTipText', () => {
      const tooltipText = disabledToolTip.find('.tooltiptext');

      expect(tooltipText.exists()).toBe(false);
    });
  });

  describe('when ToolTip is enabled', () => {
    it('renders toolTipText', () => {
      const toolTipProps = { toolTipText: 'toolTipText', buttonText: 'buttonText', disabled: false };
      const enabledToolTip = renderToolTip(toolTipProps);
      const tooltipText = enabledToolTip.find('.tooltiptext');

      expect(tooltipText.exists()).toBe(true);
      expect(tooltipText.text()).toBe(toolTipProps.toolTipText);
    });

    it('does not render when toolTipText is null', () => {
      const toolTipProps = { toolTipText: null, buttonText: 'buttonText', disabled: false };
      const enabledToolTip = renderToolTip(toolTipProps);
      const tooltipText = enabledToolTip.find('.tooltiptext');

      expect(tooltipText.exists()).toBe(false);
    });
  });
});
