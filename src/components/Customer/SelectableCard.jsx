import React from 'react';
import { string, func, bool } from 'prop-types';
import { Radio } from '@trussworks/react-uswds';

import styles from './SelectableCard.module.scss';

const SelectableCard = ({ id, label, name, value, cardText, onChange, checked }) => {
  return (
    <div className={styles.cardContainer}>
      <Radio id={id} label={label} value={value} name={name} onChange={onChange} checked={checked} />
      <div className={styles.cardText}>{cardText}</div>
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
