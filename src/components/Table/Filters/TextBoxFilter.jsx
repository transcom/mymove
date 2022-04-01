import React from 'react';
import PropTypes from 'prop-types';
import { TextInput } from '@trussworks/react-uswds';

const TextBoxFilter = ({ column: { filterValue, setFilter, id } }) => {
  return (
    <TextInput
      data-testid="TextBoxFilter"
      id={id}
      name={id}
      defaultValue={filterValue || ''}
      onKeyUp={(e) => {
        if (e.key === 'Enter') {
          setFilter(e.target.value || undefined); // Set undefined to remove the filter entirely
        }
      }}
      onBlur={(e) => {
        setFilter(e.target.value || undefined); // Set undefined to remove the filter entirely
      }}
      type="text"
    />
  );
};

// Values come from react-table
TextBoxFilter.propTypes = {
  column: PropTypes.shape({
    filterValue: PropTypes.node,
    setFilter: PropTypes.func,
    id: PropTypes.string,
  }).isRequired,
};

export default TextBoxFilter;
