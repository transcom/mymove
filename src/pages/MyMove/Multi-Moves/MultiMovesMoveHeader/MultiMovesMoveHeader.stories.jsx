import React from 'react';

import MultiMovesMoveHeader from './MultiMovesMoveHeader';

export default {
  title: 'Customer Components / MultiMovesMoveHeader',
};

export const MultiMoveCurrentMoveHeader = () => <MultiMovesMoveHeader title="Current Move" />;

export const MultiMovePreviousMoveHeader = () => <MultiMovesMoveHeader title="Previous Moves" />;
