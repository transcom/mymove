import React from 'react';
import ReactDOM from 'react-dom';

import RangeSlider, { CalculateSliderPositions } from '.';

const onChangeMock = () => {};

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <RangeSlider defaultValue={0} min={0} max={100} step={50} id="tester" onChange={onChangeMock} />,
    div,
  );
});

it('calculateSliderPositions returns proper value', () => {
  let calculations = CalculateSliderPositions(1007, 20, 500, 0, 10000, 500);
  expect(calculations.isTooltipWithinBufferSpaceBoundary).toBe(true);
  expect(calculations.tooltipLeftMargin).toBe('43px');

  // situation where it should return false for isToolTipWithinBufferSpaceBoundary
  calculations = CalculateSliderPositions(1007, 20, 10000, 0, 10000, 500);
  expect(calculations.isTooltipWithinBufferSpaceBoundary).toBe(false);
});
