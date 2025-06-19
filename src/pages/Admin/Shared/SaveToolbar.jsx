import React from 'react';
import { DeleteButton, SaveButton, Toolbar } from 'react-admin';

import adminStyles from '../adminStyles.module.scss';

const SaveToolbar = ({ showDeleteBtn }) => {
  return (
    <Toolbar className={adminStyles.flexRight} sx={{ gap: '10px' }}>
      {showDeleteBtn ? <DeleteButton /> : null}
      <SaveButton />
    </Toolbar>
  );
};

export default SaveToolbar;
