import React from 'react';
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
});
