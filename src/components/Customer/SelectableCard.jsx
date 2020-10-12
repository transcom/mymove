import React from 'react';
import { string, func, bool, node } from 'prop-types';
import { Radio } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './SelectableCard.module.scss';

const SelectableCard = ({ id, label, name, value, cardText, onChange, disabled, checked }) => {
  return (
    <div className={classnames(styles.cardContainer, { [styles.selected]: checked })}>
      <Radio
        id={id}
        label={label}
        value={value}
        name={name}
        onChange={onChange}
        checked={checked}
        disabled={disabled}
      />
      <div data-testid="selectableCardText" className={styles.cardText}>
        {cardText}
      </div>
    </div>
  );
};

SelectableCard.propTypes = {
  id: string.isRequired,
  label: string.isRequired,
  name: string.isRequired,
  value: string.isRequired,
  cardText: node,
  onChange: func.isRequired,
  checked: bool,
  disabled: bool,
};

SelectableCard.defaultProps = {
  cardText: '',
  checked: false,
  disabled: false,
};

export default SelectableCard;
