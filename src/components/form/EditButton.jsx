import React from 'react';

import { Button } from '@trussworks/react-uswds';
import { ReactComponent as EditIcon } from '../../shared/images/edit-24px.svg';

export const EditButton = (props) => (
  /* eslint-disable-next-line react/jsx-props-no-spreading */
  <Button icon {...props}>
    <span className="icon">
      <EditIcon />
    </span>
    <span>Edit</span>
  </Button>
);

export default EditButton;
