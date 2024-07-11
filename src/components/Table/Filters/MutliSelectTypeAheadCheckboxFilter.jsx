/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { bool, string, shape, node, func, arrayOf } from 'prop-types';
import AsyncSelect, { components } from 'react-select';
import { Checkbox } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './MultiSelectTypeAheadCheckBoxFilter.module.scss';

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

Option.propTypes = {
  isSelected: bool.isRequired,
  label: string.isRequired,
  innerProps: shape({
    id: string.isRequired,
  }).isRequired,
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

ValueContainer.propTypes = {
  children: node.isRequired,
};

const MultiValueContainer = ({ data: { label } }) => {
  return <span data-testid="multi-value-container">{label}</span>;
};

MultiValueContainer.propTypes = {
  data: shape({
    label: string.isRequired,
  }).isRequired,
};

const MultiSelectTypeAheadCheckBoxFilter = ({ options, placeholder, column: { filterValue, setFilter } }) => {
  const onChange = (value) => {
    let paramFilterValue = [];
    if (value) {
      value.forEach((val) => {
        paramFilterValue.push(`${val.value}`);
      });
    } else {
      paramFilterValue = undefined;
    }
    setFilter(paramFilterValue || undefined);
  };

  return (
    <div data-testid="MultiSelectTypeAheadCheckBoxFilter">
      <AsyncSelect
        classNamePrefix="MultiSelectTypeAheadCheckBoxFilter"
        className={styles.MultiSelectTypeAheadCheckBoxFilterWrapper}
        options={options}
        defaultValue={filterValue || undefined}
        onChange={onChange}
        hideSelectedOptions={false}
        isClearable={false}
        components={{ DropdownIndicator, ValueContainer, MultiValueContainer, Option }}
        placeholder={placeholder || 'Start typing...'}
        isMulti
      />
    </div>
  );
};

// Values come from react-table
MultiSelectTypeAheadCheckBoxFilter.propTypes = {
  options: arrayOf(
    shape({
      label: string.isRequired,
      value: string.isRequired,
    }).isRequired,
  ).isRequired,
  column: shape({
    filterValue: node,
    setFilter: func,
  }).isRequired,
};

export default MultiSelectTypeAheadCheckBoxFilter;
