import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

export const EditButton = ({ label, ...props }) => (
  /* eslint-disable-next-line react/jsx-props-no-spreading */
  <Button {...props}>
    <span className="icon">
      <FontAwesomeIcon data-testid="edit-icon" icon="pen" />
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
  <Button {...props}>
    <span className="icon">
      <FontAwesomeIcon data-testid="docs-icon" icon="file" />
    </span>
    <span>{label}</span>
  </Button>
);

DocsButton.propTypes = {
  label: PropTypes.string.isRequired,
};

export const ReviewButton = ({ label, ...props }) => (
  /* eslint-disable-next-line react/jsx-props-no-spreading */
  <Button {...props}>
    <span>{label}</span>
  </Button>
);

ReviewButton.defaultProps = {
  label: 'Review Documents',
};

ReviewButton.propTypes = {
  label: PropTypes.string,
};
