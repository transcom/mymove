import React from 'react';
import PropTypes from 'prop-types';

import SingleDatePicker from '../../../shared/JsonSchemaForm/SingleDatePicker';

import { formatDateForSwagger } from 'shared/dates';

const DateSelectFilter = ({ column: { filterValue, setFilter } }) => {
  // eslint-disable-next-line react/prop-types

  return (
    <SingleDatePicker
      value={filterValue || ''}
      format="DD MMM YYYY"
      data-testid="DateSelectFilter"
      placeholder=""
      onChange={(e) => {
        setFilter(formatDateForSwagger(e) || undefined); // Set undefined to remove the filter entirely
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
};

export default DateSelectFilter;
