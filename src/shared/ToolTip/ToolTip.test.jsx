import React from 'react';
import { render, screen } from '@testing-library/react';
import { mount } from 'enzyme';

import ToolTip from './ToolTip';

describe('ToolTip', () => {
  it('should display tooltip on top', () => {
    const text = 'Test Text';
    const component = mount(<ToolTip text={text} icon="circle-question" />);

    // Simulate a click on the tooltip container
    component.find('.tooltipContainer').simulate('click');

    // Find the tooltip content after the click with top class
    const tooltipContent = component.find('.tooltipTextTop');

    // Assert that the tooltip content is displayed
    expect(tooltipContent.text()).toBe(text);
  });
  it('should display tooltip on bottom', () => {
    const text = 'Test Text';
    const component = mount(<ToolTip text={text} icon="circle-question" position="bottom" />);

    // Simulate a click on the tooltip container
    component.find('.tooltipContainer').simulate('click');

    // Find the tooltip content after the click with bottom class
    const tooltipContent = component.find('.tooltipTextBottom');

    // Assert that the tooltip content is displayed
    expect(tooltipContent.text()).toBe(text);
  });
  it('should display tooltip on right', () => {
    const text = 'Test Text';
    const component = mount(<ToolTip text={text} icon="circle-question" position="right" />);

    // Simulate a click on the tooltip container
    component.find('.tooltipContainer').simulate('click');

    // Find the tooltip content after the click with right class
    const tooltipContent = component.find('.tooltipTextRight');

    // Assert that the tooltip content is displayed
    expect(tooltipContent.text()).toBe(text);
  });
  it('should display tooltip on left', () => {
    const text = 'Test Text';
    const component = mount(<ToolTip text={text} icon="circle-question" position="left" />);

    // Simulate a click on the tooltip container
    component.find('.tooltipContainer').simulate('click');

    // Find the tooltip content after the click with left class
    const tooltipContent = component.find('.tooltipTextLeft');

    // Assert that the tooltip content is displayed
    expect(tooltipContent.text()).toBe(text);
  });
  it('verify data-testid is present', () => {
    const text = 'Test Text';
    render(<ToolTip text={text} icon="circle-question" position="left" />);

    // Verify data-testid is present
    const tooltipIcon = screen.getByTestId('tooltip-container');
    expect(tooltipIcon).toBeInTheDocument();
  });
});
