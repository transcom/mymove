import React, { Component } from 'react';
import PropTypes from 'prop-types';
import styles from './RangeSlider.module.scss';
import { detectIE11 } from '../utils';

const toolTipEOLBufferPercentageModifier = 0.01;

class RangeSlider extends Component {
  onInput = event => {
    let output = document.getElementById('output-' + this.props.id);
    let slider = document.getElementById(this.props.id);
    let ticks = event.target.valueAsNumber / event.target.step;
    let possibleTicks = event.target.max / event.target.step - 1;
    let pxPerTick = slider.offsetWidth / possibleTicks;
    let leftMarginInPx = pxPerTick * ticks;
    let halfToolTipWidth = output.offsetWidth / 2;

    // This if calculation is in place to ensure that the tooltip as it moves does not go onto a new line in mobile.
    // toolTipEOLBufferPercentageModifier is used to create a percentage chunk of the total width of the slider
    // that is used to shorten the possible 'track' for the tooltip just enough to keep this from happening.
    // The other half of the and ensures that the tooltip does not move until the slider thumb has traveled more than
    // half of the tooltip's own width.
    if (
      leftMarginInPx + output.offsetWidth <
        slider.offsetWidth + halfToolTipWidth - slider.offsetWidth * toolTipEOLBufferPercentageModifier &&
      leftMarginInPx > halfToolTipWidth
    ) {
      output.style.marginLeft = leftMarginInPx - halfToolTipWidth + 'px';
    }
    output.innerText = `${this.props.prependTooltipText} ${event.target.valueAsNumber} ${this.props.appendTooltipText}`.trim();

    // This is here to allow someone to pass in a function from the parent calling this that is
    // managing a state object that relies on the value of this slider. The function being passed in needs
    // to take a numerical value as its arguement. Here's an example from Weight.jsx:
    //         onWeightSelecting = value => {
    //           this.setState({
    //             pendingPpmWeight: value,
    //           });
    //         };
    if (this.props.stateChangeFunc) {
      this.props.stateChangeFunc(event.target.valueAsNumber);
    }
  };

  onChange = value => {
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
          className={`${styles['rangeslider-output']} border-base border-1px radius-lg padding-left-1 padding-right-1`}
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
        <span className={styles['slider-min-label']}>{min}</span>
        <span className={styles['slider-max-label']}>{max}</span>
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
