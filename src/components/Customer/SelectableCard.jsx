import React from 'react';
import { string, func, bool, node } from 'prop-types';
import { Radio, Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './SelectableCard.module.scss';

const SelectableCard = ({ id, label, name, value, cardText, onChange, disabled, checked, onHelpClick }) => {
  return (
    <div>
      <div className={styles.cardContainer}>
        <Radio
          id={id}
          label={label}
          value={value}
          name={name}
          onChange={onChange}
          checked={checked}
          disabled={disabled}
          labelDescription={cardText}
          data-testid="radio"
          tile
        />
        {onHelpClick && (
          <Button
            data-testid="helpButton"
            type="button"
            onClick={onHelpClick}
            unstyled
            className={styles.helpButton}
            aria-label="help"
          >
            <FontAwesomeIcon icon={['far', 'circle-question']} />
          </Button>
        )}
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
  onHelpClick: func,
};

SelectableCard.defaultProps = {
  cardText: '',
  checked: false,
  disabled: false,
  onHelpClick: null,
};

export default SelectableCard;
