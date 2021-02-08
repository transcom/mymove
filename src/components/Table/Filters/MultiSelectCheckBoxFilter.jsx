/* eslint-disable react/jsx-props-no-spreading,react/prop-types */
import React, { useMemo } from 'react';
import PropTypes from 'prop-types';
import Select, { components } from 'react-select';
import { Checkbox } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './MultiSelectCheckBoxFilter.module.scss';

const Option = (props) => {
  const {
    isSelected,
    label,
    innerProps: { id },
  } = props;
  return (
    <components.Option {...props}>
      <Checkbox id={id} name={label} label={label} checked={isSelected} onChange={() => {}} />
    </components.Option>
  );
};

const DropdownIndicator = (props) => {
  return (
    <components.DropdownIndicator {...props}>
      <FontAwesomeIcon className="fas fa-sort" icon="sort" />
    </components.DropdownIndicator>
  );
};

const ValueContainer = ({ children, ...props }) => {
  return (
    <components.ValueContainer {...props}>
      <div>{children}</div>
    </components.ValueContainer>
  );
};

const MultiValueContainer = ({ data: { label } }) => {
  return <span data-testid="multi-value-container">{label}</span>;
};

const MultiSelectCheckBoxFilter = ({ options, column: { filterValue, setFilter } }) => {
  const selectFilterValue = useMemo(() => {
    return filterValue
      ? filterValue.split(',').map((val) => ({
          label: options.find((option) => option.value === val).label,
          value: val,
        }))
      : [];
  }, [filterValue, options]);

  return (
    <div data-testid="MultiSelectCheckBoxFilter">
      <Select
        classNamePrefix="MultiSelectCheckBoxFilter"
        className={styles.MultiSelectCheckBoxFilterWrapper}
        options={options}
        defaultValue={selectFilterValue || undefined}
        onChange={(value) => {
          let paramFilterValue = '';
          if (value) {
            value.forEach((val, index) => {
              paramFilterValue += `${val.value}`;
              if (index + 1 !== value.length) {
                paramFilterValue += ',';
              }
            });
          } else {
            paramFilterValue = undefined;
          }
          // value example: {label: "Fort Gordon", value: "Fort Gordon"}
          // value example to be converted to send back to react-query: { id: 'status', value:'New move' }
          setFilter(paramFilterValue || undefined);
        }}
        isMulti
        isSearchable={false}
        hideSelectedOptions={false}
        isClearable={false}
        components={{ DropdownIndicator, ValueContainer, MultiValueContainer, Option }}
      />
    </div>
  );
};

// Values come from react-table
MultiSelectCheckBoxFilter.propTypes = {
  options: PropTypes.arrayOf(PropTypes.object).isRequired,
  column: PropTypes.shape({
    filterValue: PropTypes.node,
    setFilter: PropTypes.func,
  }).isRequired,
};

export default MultiSelectCheckBoxFilter;
