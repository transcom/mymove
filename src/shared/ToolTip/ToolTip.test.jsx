import React from 'react';
import { mount } from 'enzyme';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';

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

  it('should display tooltip with title when provided', () => {
    const titleText = 'Tooltip Title';
    const bodyText = 'Tooltip Body';

    const component = mount(<ToolTip text={bodyText} title={titleText} icon="circle-question" />);

    component.find('.tooltipContainer').simulate('click');

    const tooltipTitle = component.find('.popoverHeader');
    const tooltipBody = component.find('.popoverBody');

    expect(tooltipTitle.text()).toBe(titleText);
    expect(tooltipBody.text()).toBe(bodyText);
  });

  it('should not display title when it is not provided', () => {
    const bodyText = 'Tooltip Body';

    const component = mount(<ToolTip text={bodyText} icon="circle-question" />);

    component.find('.tooltipContainer').simulate('click');

    // Find the tooltip content and assert only the body is displayed
    const tooltipTitle = component.find('.popoverHeader');
    const tooltipBody = component.find('.popoverBody');

    expect(tooltipTitle.exists()).toBe(false);
    expect(tooltipBody.text()).toBe(bodyText);
  });

  it('verify data-testid is present', () => {
    const text = 'Test Text';
    render(<ToolTip text={text} icon="circle-question" position="left" />);

    // Verify data-testid is present
    const tooltipIcon = screen.getByTestId('tooltip-container');
    expect(tooltipIcon).toBeInTheDocument();
  });

  it('should display a large tooltip', () => {
    const text = 'Test Text';
    const component = mount(<ToolTip text={text} icon="circle-question" position="top" textAreaSize="large" />);

    // Simulate a click on the tooltip container
    component.find('.tooltipContainer').simulate('click');

    // Find the tooltip content after the click with left class
    const tooltipContent = component.find('.tooltipTextTop.toolTipTextAreaLarge');

    // Assert that the tooltip content is displayed
    expect(tooltipContent.text()).toBe(text);
  });
  it('should close tooltip on mouseleave if closeOnLeave prop is passed in', async () => {
    const text = 'Test Text';
    const wrapper = render(<ToolTip text={text} icon="circle-question" position="left" closeOnLeave={true} />);

    // Verify data-testid is present
    const tooltipIcon = screen.getByTestId('tooltip-container');
    expect(tooltipIcon).toBeInTheDocument();

    // Mouseover and view tooltip
    fireEvent.mouseEnter(tooltipIcon);
    const tooltipText = await waitFor(() => wrapper.getByTestId('tooltipText'));
    expect(tooltipText).toBeVisible();

    // Assert that the tooltip content is displayed
    expect(tooltipText.textContent).toBe(text);

    fireEvent.mouseLeave(tooltipIcon);
    expect(tooltipText).not.toBeVisible();
  });
});
