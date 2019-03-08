import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './index.css';

export class ToolTip extends Component {
  render() {
    let { children, text, disabled, textStyle } = this.props;
    return (
      <span className="tooltip">
        {children}
        {text && !disabled && <span className={`tooltiptext ${textStyle}`}>{text}</span>}
      </span>
    );
  }
}

ToolTip.propTypes = {
  text: PropTypes.string,
  disabled: PropTypes.bool,
  textStyle: PropTypes.string,
};

export default ToolTip;
