import React from 'react';
import PropTypes from 'prop-types';

import SingleDatePicker from '../../../shared/JsonSchemaForm/SingleDatePicker';

import { formatDateForSwagger, formatDateTime } from 'shared/dates';

// Return function with proper type
const DateSelectFilter = ({ dateTime, column: { filterValue, setFilter } }) => {
  return (
    <SingleDatePicker
      value={filterValue || ''}
      format="DD MMM YYYY"
      data-testid="DateSelectFilter"
      placeholder=""
      onChange={(e) => {
        // note that `e` here is a Date, not a string. Fortunately
        // formatDateTime accidentally handles that just fine
        setFilter(dateTime ? formatDateTime(e) || undefined : formatDateForSwagger(e) || undefined); // Set undefined to remove the filter entirely
      }}
    />
  );
};

//
DateSelectFilter.defaultProps = {
  dateTime: false,
};

// Values come from react-table
DateSelectFilter.propTypes = {
  column: PropTypes.shape({
    filterValue: PropTypes.string,
    setFilter: PropTypes.func,
  }).isRequired,
  dateTime: PropTypes.bool,
};

export default DateSelectFilter;
