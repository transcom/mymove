import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './index.css';

class ToolTip extends Component {
  render() {
    let { children, toolTipText, disabled, textStyle } = this.props;
    return (
      <span className="tooltip">
        {!disabled && toolTipText && <span className={`tooltiptext ${textStyle}`}>{toolTipText}</span>}
        {children}
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
