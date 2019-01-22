import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import './index.css';
const BigButton = ({ selected, children, onClick, className, isDisabled }) => (
  <button className={classnames('big-button', className, { selected })} onClick={onClick} disabled={isDisabled}>
    {children}
  </button>
);

BigButton.propTypes = {
  className: PropTypes.string,
  children: PropTypes.node,
  selected: PropTypes.bool,
  onClick: PropTypes.func,
  isDisabled: PropTypes.bool,
};

export default BigButton;
