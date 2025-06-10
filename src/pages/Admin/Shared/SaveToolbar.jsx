import React from 'react';
import { DeleteButton, SaveButton, Toolbar } from 'react-admin';

import adminStyles from '../adminStyles.module.scss';

const SaveToolbar = ({ showDeleteBtn }) => {
  return (
    <Toolbar className={adminStyles.flexRight} sx={{ gap: '10px' }}>
      {showDeleteBtn ? (
        <DeleteButton
          sx={{
            backgroundColor: '#e1400a !important',
            width: 120,
            '&:hover': {
              opacity: '0.8',
            },
          }}
        />
      ) : null}
      <SaveButton />
    </Toolbar>
  );
};

export default SaveToolbar;
