import React from 'react';

import { mockMovesPCS, mockMovesRetirement, mockMovesSeparation } from '../MultiMovesTestData';

import MultiMovesMoveContainer from './MultiMovesMoveContainer';

export default {
  title: 'Customer Components / MultiMovesContainer',
};

export const PCSCurrentMove = () => <MultiMovesMoveContainer moves={mockMovesPCS.currentMove} />;

export const PCSPreviousMoves = () => <MultiMovesMoveContainer moves={mockMovesPCS.previousMoves} />;

export const RetirementCurrentMove = () => <MultiMovesMoveContainer moves={mockMovesRetirement.currentMove} />;

export const RetirementPreviousMoves = () => <MultiMovesMoveContainer moves={mockMovesRetirement.previousMoves} />;

export const SeparationCurrentMove = () => <MultiMovesMoveContainer moves={mockMovesSeparation.currentMove} />;

export const SeparationPreviousMoves = () => <MultiMovesMoveContainer moves={mockMovesSeparation.previousMoves} />;
