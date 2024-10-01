import React from 'react';
import { PropTypes } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './CancelMoveButton.module.scss';

function CancelMoveButton({ onClick, isMoveLocked }) {
  return (
    <div>
      <Button
        type="Button"
        className={classnames(styles.CancelMoveButton, ['usa-button--unstyled'])}
        onClick={onClick}
        disabled={isMoveLocked}
      >
        Cancel move
      </Button>
    </div>
  );
}

CancelMoveButton.propTypes = {
  onClick: PropTypes.func.isRequired,
};

export default CancelMoveButton;
