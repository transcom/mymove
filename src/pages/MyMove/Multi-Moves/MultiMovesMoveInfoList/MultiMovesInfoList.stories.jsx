import React from 'react';

import { mockMovesPCS, mockMovesRetirement, mockMovesSeparation } from '../MultiMovesTestData';

import MultiMovesMoveInfoList from './MultiMovesMoveInfoList';

export default {
  title: 'Customer Components / MultiMovesMoveInfoList',
};

const currentMovePCS = mockMovesPCS.currentMove[0];
const currentMoveRetirement = mockMovesRetirement.currentMove[0];
const currentMoveSeparation = mockMovesSeparation.currentMove[0];

export const MultiMoveHeaderPCS = () => <MultiMovesMoveInfoList move={currentMovePCS} />;

export const MultiMoveHeaderRetirement = () => <MultiMovesMoveInfoList move={currentMoveRetirement} />;

export const MultiMoveHeaderSeparation = () => <MultiMovesMoveInfoList move={currentMoveSeparation} />;
