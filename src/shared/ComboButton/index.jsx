import React, { Fragment } from 'react';
import PropTypes from 'prop-types';

import './index.css';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCaretDown from '@fortawesome/fontawesome-free-solid/faCaretDown';

const ComboButton = ({ toolTipText, isDisabled }) => (
  <Fragment>
    <span className="button-tooltip tooltip">
      <button disabled={isDisabled}>
        Approve&nbsp;&nbsp;&nbsp;
        <FontAwesomeIcon className="icon" icon={faCaretDown} />
      </button>
      <span className="tooltiptext">{toolTipText}</span>
    </span>
  </Fragment>
);

ComboButton.propTypes = {
  toolTipText: PropTypes.string,
  isDisabled: PropTypes.bool,
};

export default ComboButton;
