import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { ReactComponent as EditIcon } from '../../shared/icon/edit.svg';
import { ReactComponent as DocsIcon } from '../../shared/icon/documents.svg';

export const EditButton = ({ label, ...props }) => (
  /* eslint-disable-next-line  react/jsx-props-no-spreading */
  <Button icon {...props}>
    <span className="icon">
      <EditIcon />
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
  /* eslint-disable react/jsx-props-no-spreading */
  <Button icon {...props}>
    <span className="icon">
      <DocsIcon />
    </span>
    <span>{label}</span>
  </Button>
);

DocsButton.propTypes = {
  label: PropTypes.string.isRequired,
};
