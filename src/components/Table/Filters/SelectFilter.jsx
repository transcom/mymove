import React from 'react';
import PropTypes from 'prop-types';
import { Dropdown } from '@trussworks/react-uswds';

const SelectFilter = ({ options, column: { filterValue, setFilter } }) => {
  return (
    <Dropdown
      data-testid="SelectFilter"
      defaultValue={filterValue}
      onChange={(e) => {
        setFilter(e.target.value);
      }}
      style={{ width: 'auto' }}
    >
      {options.map(({ value, label }) => (
        <option value={value} key={`filterOption_${value}`}>
          {label}
        </option>
      ))}
    </Dropdown>
  );
};

const OptionsShape = PropTypes.object;

// Values come from react-table
SelectFilter.propTypes = {
  options: PropTypes.arrayOf(OptionsShape).isRequired,
  column: PropTypes.shape({
    filterValue: PropTypes.node,
    setFilter: PropTypes.func,
  }).isRequired,
};

export default SelectFilter;
