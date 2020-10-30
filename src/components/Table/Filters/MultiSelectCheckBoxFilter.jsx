/* eslint-disable react/jsx-props-no-spreading,react/prop-types */
import React, { useMemo } from 'react';
import PropTypes from 'prop-types';
import Select, { components } from 'react-select';
import { Checkbox } from '@trussworks/react-uswds';

import styles from './MultiSelectCheckBoxFilter.module.scss';

import { ReactComponent as ArrowBoth } from 'icons/arrow-both.svg';

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
      <ArrowBoth />
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

const MultiValueContainer = ({ data: { value } }) => {
  return <span>{value}</span>;
};

const MultiSelectCheckBoxFilter = ({ options, column: { filterValue, setFilter } }) => {
  const selectFilterValue = useMemo(() => {
    return filterValue
      ? filterValue.split(',').map((val) => ({
          label: val,
          value: val,
        }))
      : [];
  }, [filterValue]);
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
        components={{ DropdownIndicator, ValueContainer, MultiValueContainer, Option }}
        isClearable
      />
    </div>
  );
};

// Values come from react-table
MultiSelectCheckBoxFilter.propTypes = {
  options: PropTypes.arrayOf(PropTypes.object).isRequired,
  column: PropTypes.shape({
    filterValue: PropTypes.node,
    preFilteredRows: PropTypes.array,
    setFilter: PropTypes.func,
  }).isRequired,
};

export default MultiSelectCheckBoxFilter;
