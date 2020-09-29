import React from 'react';
import { string, func, bool } from 'prop-types';
import { Radio } from '@trussworks/react-uswds';

const SelectableCard = ({ id, label, name, value, cardText, onChange, checked }) => {
  return (
    <div>
      <Radio id={id} label={label} value={value} name={name} onChange={onChange} checked={checked} />
      <div>{cardText}</div>
    </div>
  );
};

SelectableCard.propTypes = {
  id: string.isRequired,
  label: string.isRequired,
  name: string.isRequired,
  value: string.isRequired,
  cardText: string,
  onChange: func.isRequired,
  checked: bool,
};

SelectableCard.defaultProps = {
  cardText: '',
  checked: false,
};

export default SelectableCard;
