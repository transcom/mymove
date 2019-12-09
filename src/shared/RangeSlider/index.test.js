import React from 'react';
import ReactDOM from 'react-dom';

import RangeSlider, { calculateSliderPositions } from '.';

const onChangeMock = () => {};

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <RangeSlider defaultValue={0} min={0} max={100} step={50} id="tester" onChange={onChangeMock} />,
    div,
  );
});

it('calculateSliderPositions returns proper value', () => {
  const calculations = calculateSliderPositions(1007, 20, 500, 0, 10000, 500);
  expect(calculations.isTooltipWithinBufferSpaceBoundary).toBe(true);
  expect(calculations.tooltipLeftMargin).toBe('43px');
});

it('calculateSliderPositions identifies when the buffer space is outside the proper boundary', () => {
  const calculations = calculateSliderPositions(1007, 20, 10000, 0, 10000, 500);
  expect(calculations.isTooltipWithinBufferSpaceBoundary).toBe(false);
});
