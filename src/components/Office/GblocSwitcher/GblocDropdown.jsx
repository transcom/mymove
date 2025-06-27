import React from 'react';
import PropTypes from 'prop-types';
import { Dropdown } from '@trussworks/react-uswds';

import styles from './GblocSwitcher.module.scss';

export const gblocDropdownTestId = 'gbloc_switcher';

const GblocDropdown = ({ handleChange, ariaLabel, defaultValue, gblocs }) => {
  const testId = gblocDropdownTestId;

  return (
    <Dropdown
      data-testid={testId}
      className={styles.switchGblocButton}
      onChange={(e) => handleChange(e.target.value)}
      defaultValue={defaultValue}
      aria-label={ariaLabel}
    >
      {gblocs?.map((gbloc) => (
        <option value={gbloc} key={`filterOption_${gbloc}`}>
          {gbloc}
        </option>
      ))}
    </Dropdown>
  );
};

GblocDropdown.defaultProps = {
  handleChange: () => {},
  defaultValue: '',
  ariaLabel: undefined,
  gblocs: [],
};

GblocDropdown.propTypes = {
  handleChange: PropTypes.func,
  defaultValue: PropTypes.string,
  ariaLabel: PropTypes.string,
  gblocs: PropTypes.arrayOf(PropTypes.string),
};

export default GblocDropdown;
