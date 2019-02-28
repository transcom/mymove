import React from 'react';
import PropTypes from 'prop-types';

import './index.css';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCaretDown from '@fortawesome/fontawesome-free-solid/faCaretDown';

const ComboButton = ({ toolTipText, isDisabled }) => (
  <span className="combo-button button-tooltip tooltip">
    <button disabled={isDisabled}>
      Approve
      <FontAwesomeIcon className="combo-button-icon" icon={faCaretDown} />
    </button>
    {toolTipText && <span className="tooltiptext">{toolTipText}</span>}
  </span>
);

ComboButton.propTypes = {
  toolTipText: PropTypes.string,
  isDisabled: PropTypes.bool,
};

export default ComboButton;
