import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEdit, faFile } from '@fortawesome/free-solid-svg-icons';

export const EditButton = ({ label, ...props }) => (
  /* eslint-disable-next-line react/jsx-props-no-spreading */
  <Button icon {...props}>
    <span className="icon">
      <FontAwesomeIcon icon={faEdit} />
    </span>
    <span>{label}</span>
  </Button>
);

EditButton.defaultProps = {
  label: 'Edit',
};

EditButton.propTypes = {
  label: PropTypes.string,
};

export const DocsButton = ({ label, ...props }) => (
  /* eslint-disable-next-line react/jsx-props-no-spreading */
  <Button icon {...props}>
    <span className="icon">
      <FontAwesomeIcon icon={faFile} />
    </span>
    <span>{label}</span>
  </Button>
);

DocsButton.propTypes = {
  label: PropTypes.string.isRequired,
};
