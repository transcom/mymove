/* eslint-disable react/no-array-index-key */
import React from 'react';
import PropTypes from 'prop-types';
import { Dropdown } from '@trussworks/react-uswds';

const SelectFilter = ({ options, column: { filterValue, preFilteredRows, setFilter } }) => {
  return (
    <Dropdown
      data-testid="SelectFilter"
      defaultValue={filterValue}
      onChange={(e) => {
        setFilter(e.target.value);
      }}
      style={{ width: 'auto' }}
    >
      <option>Search {preFilteredRows.length} records...</option>
      {options.map(({ value, label }) => (
        <option value={value} key={`filterOption_${value}`}>
          {label}
        </option>
      ))}
    </Dropdown>
  );
};

// Values come from react-table
SelectFilter.propTypes = {
  options: PropTypes.arrayOf(PropTypes.object).isRequired,
  column: PropTypes.shape({
    filterValue: PropTypes.node,
    preFilteredRows: PropTypes.array,
    setFilter: PropTypes.func,
  }).isRequired,
};

export default SelectFilter;
