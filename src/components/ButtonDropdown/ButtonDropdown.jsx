import React from 'react';
import { Dropdown } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './ButtonDropdown.module.scss';

const ButtonDropdown = ({ children, onChange, value, ariaLabel, divClassName, testId }) => (
  <div className={classnames(styles.ButtonDropdown, divClassName)} data-testid={testId}>
    <Dropdown aria-label={ariaLabel} onChange={onChange} className={styles.ButtonDropdown} value={value}>
      {children}
    </Dropdown>
    <span className={styles.ButtonDropdownIcon} />
  </div>
);

ButtonDropdown.defaultProps = {
  ariaLabel: '',
  divClassName: '',
};

ButtonDropdown.propTypes = {
  children: PropTypes.node.isRequired,
  onChange: PropTypes.func.isRequired,
  ariaLabel: PropTypes.string,
  divClassName: PropTypes.string,
};

export default ButtonDropdown;
