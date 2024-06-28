import React from 'react';
import PropTypes from 'prop-types';
import { TextInput } from '@trussworks/react-uswds';
import { v4 as uuidv4 } from 'uuid';

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
      // Generate unique key. This is used to ensure filter column headers are cleared properly
      // for the filtering toggling via pill widget.
      key={uuidv4()}
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
