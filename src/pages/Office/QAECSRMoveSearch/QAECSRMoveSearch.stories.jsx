import React from 'react';
import { MemoryRouter } from 'react-router';

import QAECSRMoveSearch from './QAECSRMoveSearch';

export default {
  title: 'Office Components/QAECSRMoveSearch',
  component: QAECSRMoveSearch,
};

const defaultProps = {};

export const MoveSearch = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
    <MemoryRouter initialEntries={['/']}>
      <QAECSRMoveSearch {...defaultProps} />
    </MemoryRouter>
  </div>
);
