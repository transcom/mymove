import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import './dropdown.css';

export const DropDown = props => {
  return <div className="dropdown">{props.children}</div>;
};

export class DropDownItem extends Component {
  handleClick = () => {
    if (!this.props.disabled) {
      this.props.onClick();
    }
  };

  render() {
    const liClasses = props => classNames({ disabled: props.disabled });
    return (
      <p className={liClasses(this.props)} onClick={this.handleClick}>
        {this.props.value}
      </p>
    );
  }
}

DropDown.propTypes = {
  children: PropTypes.oneOfType([PropTypes.func, PropTypes.node]),
};

DropDownItem.propTypes = {
  onClick: PropTypes.func,
  value: PropTypes.string,
  disabled: PropTypes.bool,
};
