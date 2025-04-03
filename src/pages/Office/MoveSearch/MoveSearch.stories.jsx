import React from 'react';
import { MemoryRouter } from 'react-router-dom';

import MoveSearch from './MoveSearch';

export default {
  title: 'Office Components/MoveSearch',
  component: MoveSearch,
};

const defaultProps = {};

export const StoryMoveSearch = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
    <MemoryRouter>
      <MoveSearch {...defaultProps} />
    </MemoryRouter>
  </div>
);
