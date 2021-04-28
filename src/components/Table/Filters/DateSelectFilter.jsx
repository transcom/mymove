import React from 'react';
import PropTypes from 'prop-types';

import SingleDatePicker from '../../../shared/JsonSchemaForm/SingleDatePicker';

import { formatDateForSwagger, formatDateTime } from 'shared/dates';

// Return function with proper type
const DateSelectFilter = ({ type, column: { filterValue, setFilter } }) => {
  return (
    <SingleDatePicker
      value={filterValue || ''}
      format="DD MMM YYYY"
      data-testid="DateSelectFilter"
      placeholder=""
      onChange={(e) => {
        setFilter(type === 'DateTime' ? formatDateTime(e) || undefined : formatDateForSwagger(e) || undefined); // Set undefined to remove the filter entirely
      }}
    />
  );
};

// Values come from react-table
DateSelectFilter.propTypes = {
  column: PropTypes.shape({
    filterValue: PropTypes.string,
    setFilter: PropTypes.func,
  }).isRequired,
  type: PropTypes.string.isRequired,
};

export default DateSelectFilter;
