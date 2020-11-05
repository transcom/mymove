import React from 'react';
import { string, func, bool, node } from 'prop-types';
import { Radio, Button } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import faQuestionCircle from '@fortawesome/free-regular-svg-icons/faQuestionCircle';

import styles from './SelectableCard.module.scss';

const SelectableCard = ({ id, label, name, value, cardText, onChange, disabled, checked, onHelpClick }) => {
  return (
    <div className={classnames(styles.cardContainer, { [styles.selected]: checked })}>
      <div className={styles.cardTitle}>
        <Radio
          id={id}
          label={label}
          value={value}
          name={name}
          onChange={onChange}
          checked={checked}
          disabled={disabled}
        />
        {onHelpClick && (
          <Button data-testid="helpButton" type="button" onClick={onHelpClick} unstyled className={styles.helpButton}>
            <FontAwesomeIcon icon={faQuestionCircle} />
          </Button>
        )}
      </div>
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
  onHelpClick: func,
};

SelectableCard.defaultProps = {
  cardText: '',
  checked: false,
  disabled: false,
  onHelpClick: null,
};

export default SelectableCard;
