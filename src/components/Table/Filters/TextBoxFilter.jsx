import React from 'react';
import PropTypes from 'prop-types';

const TextBoxFilter = ({ column: { filterValue, setFilter } }) => {
  // eslint-disable-next-line react/prop-types

  return (
    <input
      data-testid="TextBoxFilter"
      defaultValue={filterValue || ''}
      onBlur={(e) => {
        setFilter(e.target.value || undefined); // Set undefined to remove the filter entirely
      }}
    />
  );
};

// Values come from react-table
TextBoxFilter.propTypes = {
  column: PropTypes.shape({
    filterValue: PropTypes.node,
    setFilter: PropTypes.func,
  }).isRequired,
};

export default TextBoxFilter;
