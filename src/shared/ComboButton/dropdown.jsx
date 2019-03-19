import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import './dropdown.css';

export const DropDown = props => {
  return <div className="dropdown">{props.children}</div>;
};

export function DropDownItem(props) {
  const dropdownItemClasses = classNames({ disabled: props.disabled });
  return (
    <div className={dropdownItemClasses} onClick={props.disabled ? null : props.onClick}>
      {props.value}
    </div>
  );
}

DropDown.propTypes = {
  children: PropTypes.oneOfType([PropTypes.func, PropTypes.node]),
};

DropDownItem.propTypes = {
  onClick: PropTypes.func,
  value: PropTypes.string,
  disabled: PropTypes.bool,
};
