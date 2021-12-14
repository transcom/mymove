import React from 'react';
import { Dropdown } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './ButtonDropdown.module.scss';

const ButtonDropdown = ({ children, onChange, ariaLabel }) => (
  <div className={styles.ButtonDropdown}>
    <Dropdown aria-label={ariaLabel} onChange={onChange} className={classnames(styles.ButtonDropdown, 'usa-button')}>
      {children}
    </Dropdown>
  </div>
);

ButtonDropdown.defaultProps = {
  ariaLabel: '',
};

ButtonDropdown.propTypes = {
  children: PropTypes.node.isRequired,
  onChange: PropTypes.func.isRequired,
  ariaLabel: PropTypes.string,
};

export default ButtonDropdown;
