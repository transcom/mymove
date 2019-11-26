import React, { Component } from 'react';
import PropTypes from 'prop-types';
import styles from './RangeSlider.module.scss';

class RangeSlider extends Component {
  onInput = event => {
    let output = document.getElementById('output__' + this.props.id);
    let slider = document.getElementById(this.props.id);
    let ticks = event.target.valueAsNumber / event.target.step;
    let possibleTicks = event.target.max / event.target.step - 1;
    let pxPerTick = slider.offsetWidth / possibleTicks;
    if (
      pxPerTick * ticks + output.offsetWidth < slider.offsetWidth + output.offsetWidth / 2 - slider.offsetWidth / 25 &&
      pxPerTick * ticks > output.offsetWidth / 2
    ) {
      output.style.marginLeft = pxPerTick * ticks - output.offsetWidth / 2 + 'px';
    }
    output.value =
      (this.props.prependTooltipText ? this.props.prependTooltipText + ' ' : '') +
      event.target.valueAsNumber +
      (this.props.appendToolTipText ? ' ' + this.props.appendToolTipText : '');

    if (this.props.stateChangeFunc) {
      this.props.stateChangeFunc(event.target.valueAsNumber);
    }
  };

  onChange = value => {
    this.props.onChange(value);
  };

  render() {
    return (
      <>
        <div className="rangeslider__container">
          <output
            className={`${styles['rangeslider-output']} border-base border-1px radius-lg padding-left-1 padding-right-1`}
            id={'output__' + this.props.id}
            htmlFor={this.props.id}
          >
            {(this.props.prependTooltipText ? this.props.prependTooltipText + ' ' : '') +
              this.props.defaultValue +
              (this.props.appendToolTipText ? ' ' + this.props.appendToolTipText : '')}
          </output>
          <input
            id={this.props.id}
            className="usa-range"
            type="range"
            min={this.props.min}
            max={this.props.max}
            step={this.props.step}
            defaultValue={this.props.defaultValue}
            onInput={this.onInput}
            onChange={this.onChange}
          />
        </div>
      </>
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
  appendToolTipText: PropTypes.string,
  stateChangeFunc: PropTypes.func,
  onChange: PropTypes.func.isRequired,
};

export default RangeSlider;
