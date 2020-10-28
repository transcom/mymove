import React from 'react';
import PropTypes from 'prop-types';

const TextBoxFilter = ({ column: { filterValue, preFilteredRows, setFilter } }) => {
  // eslint-disable-next-line react/prop-types
  const count = preFilteredRows.length;

  return (
    <input
      data-testid="TextBoxFilter"
      value={filterValue || ''}
      onChange={(e) => {
        setFilter(e.target.value || undefined); // Set undefined to remove the filter entirely
      }}
      placeholder={`Search ${count} records...`}
    />
  );
};

// Values come from react-table
TextBoxFilter.propTypes = {
  column: PropTypes.shape({
    filterValue: PropTypes.node,
    preFilteredRows: PropTypes.node,
    setFilter: PropTypes.func,
  }).isRequired,
};

export default TextBoxFilter;
