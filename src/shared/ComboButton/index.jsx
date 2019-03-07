import React from 'react';
import PropTypes from 'prop-types';

import 'shared/shared.css';
import './index.css';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCaretDown from '@fortawesome/fontawesome-free-solid/faCaretDown';

const ComboButton = ({ buttonText, toolTipText, isDisabled }) => (
  <span className="combo-button tooltip">
    <button disabled={isDisabled}>
      {buttonText}
      <FontAwesomeIcon className="combo-button-icon" icon={faCaretDown} />
    </button>
    {toolTipText && <span className="tooltiptext tooltiptext-large">{toolTipText}</span>}
  </span>
);

ComboButton.propTypes = {
  buttonText: PropTypes.string,
  toolTipText: PropTypes.string,
  isDisabled: PropTypes.bool,
};

export default ComboButton;
