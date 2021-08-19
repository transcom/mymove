import React, { Component } from 'react';
import PropTypes from 'prop-types';
import styles from './RangeSlider.module.scss';
import { detectIE11 } from '../utils';

const toolTipEOLBufferPercentageModifier = 0.01;

const calculateSliderPositions = (sliderOffset, tooltipWidth, numValue, min, max, step) => {
  let ticks = numValue / step;
  let possibleTicks = max / step - 1;
  let pxPerTick = sliderOffset / possibleTicks;
  let leftMarginInPx = pxPerTick * ticks;
  let halfToolTipWidth = tooltipWidth / 2;
  let tooltipLeftMargin = leftMarginInPx - halfToolTipWidth + 'px';

  // This if calculation is in place to ensure that the tooltip as it moves does not go onto a new line in mobile.
  // toolTipEOLBufferPercentageModifier is used to create a percentage chunk of the total width of the slider
  // that is used to shorten the possible 'track' for the tooltip just enough to keep this from happening.
  // The other half of the and ensures that the tooltip does not move until the slider thumb has traveled more than
  // half of the tooltip's own width.
  let isTooltipWithinBufferSpaceBoundary =
    leftMarginInPx + tooltipWidth <
      sliderOffset + halfToolTipWidth - sliderOffset * toolTipEOLBufferPercentageModifier &&
    leftMarginInPx > halfToolTipWidth;

  return { isTooltipWithinBufferSpaceBoundary, tooltipLeftMargin };
};

class RangeSlider extends Component {
  onInput = (event) => {
    const output = document.getElementById('output-' + this.props.id);
    const slider = document.getElementById(this.props.id);

    const calculations = calculateSliderPositions(
      slider.offsetWidth,
      output.offsetWidth,
      event.target.valueAsNumber,
      event.target.min,
      event.target.max,
      event.target.step,
    );
    output.innerText =
      `${this.props.prependTooltipText} ${event.target.valueAsNumber} ${this.props.appendTooltipText}`.trim();

    if (calculations.isTooltipWithinBufferSpaceBoundary) {
      output.style.marginLeft = calculations.tooltipLeftMargin;
    }

    if (this.props.stateChangeFunc) {
      this.props.stateChangeFunc(event.target.valueAsNumber);
    }
  };

  onChange = (value) => {
    if (detectIE11()) {
      this.onInput(value);
    }
    this.props.onChange(value);
  };

  render() {
    const { id, min, max, step, defaultValue, prependTooltipText, appendTooltipText } = this.props;
    return (
      <div className="rangeslider-container">
        <span
          className={`${styles['rangeslider-output']} border-1px radius-lg padding-left-1 padding-right-1`}
          id={'output-' + id}
          htmlFor={id}
        >
          {`${prependTooltipText} ${defaultValue} ${appendTooltipText}`.trim()}
        </span>
        <input
          id={id}
          className="usa-range"
          type="range"
          min={min}
          max={max}
          step={step}
          defaultValue={defaultValue}
          onInput={this.onInput}
          onChange={this.onChange}
        />
        <div className={styles['range-label-container']}>
          <span>{min}</span>
          <span>{max}</span>
        </div>
      </div>
    );
  }
}

RangeSlider.propTypes = {
  id: PropTypes.string.isRequired,
  min: PropTypes.number.isRequired,
  max: PropTypes.number.isRequired,
  step: PropTypes.number.isRequired,
  defaultValue: PropTypes.number.isRequired,
  prependTooltipText: PropTypes.string,
  appendTooltipText: PropTypes.string,
  stateChangeFunc: PropTypes.func,
  onChange: PropTypes.func.isRequired,
};

export default RangeSlider;
export { calculateSliderPositions };
