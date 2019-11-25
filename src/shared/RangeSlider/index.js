import React from 'react';
import PropTypes from 'prop-types';

let stateToChange;

const RangeSlider = ({ id, min, max, step, defaultValue, onChange, stateChangeObject }) => {
  stateToChange = stateChangeObject;
  return (
    <>
      <div className="rangeslider__container">
        <output htmlFor={id}> </output>
        <input
          id={id}
          className="usa-range"
          type="range"
          min={min}
          max={max}
          step={step}
          defaultValue={defaultValue}
          onInput={onInput}
          onChange={onChange}
        />
      </div>
    </>
  );
};

let onInput = event => {
  let output = document.getElementById('slider__output');
  let slider = document.getElementById('progear__weight__selector');
  let ticks = event.target.valueAsNumber / event.target.step;
  let possibleTicks = event.target.max / event.target.step - 1;
  let pxPerTick = slider.offsetWidth / possibleTicks;
  if (
    pxPerTick * ticks + output.offsetWidth < slider.offsetWidth + output.offsetWidth / 2 &&
    pxPerTick * ticks > output.offsetWidth / 2
  ) {
    output.style.marginLeft = pxPerTick * ticks - output.offsetWidth / 2 + 'px';
  }
  if (stateToChange !== null) {
    this.setState({
      stateChangeObject: event.target.value,
    });
  }
  output.value = event.target.value;
};

RangeSlider.propTypes = {
  id: PropTypes.string.isRequired,
  min: PropTypes.number.isRequired,
  max: PropTypes.number.isRequired,
  step: PropTypes.number.isRequired,
  defaultValue: PropTypes.number.isRequired,
  alwaysShowTooltip: PropTypes.bool,
};

export default RangeSlider;
