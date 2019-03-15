import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import './dropdown.css';

export const DropDown = props => {
  return <div className="dropdown">{props.children}</div>;
};

export const DropDownItem = props => {
  const liClasses = props => classNames({ disabled: props.disabled });
  return <p className={liClasses(props)}>{props.value}</p>;
};

DropDown.propTypes = {
  children: PropTypes.arrayOf(DropDownItem),
};

DropDownItem.propTypes = {
  value: PropTypes.string,
  disabled: PropTypes.bool,
};
