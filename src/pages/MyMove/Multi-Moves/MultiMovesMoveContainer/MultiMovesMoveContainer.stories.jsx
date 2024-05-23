import React from 'react';

import { mockMovesPCS, mockMovesRetirement, mockMovesSeparation } from '../MultiMovesTestData';

import MultiMovesMoveContainer from './MultiMovesMoveContainer';

import { MockProviders } from 'testUtils';

export default {
  title: 'Customer Components / MultiMovesContainer',
};

export const PCSCurrentMove = () => (
  <MockProviders>
    <MultiMovesMoveContainer moves={mockMovesPCS.currentMove} />
  </MockProviders>
);

export const PCSPreviousMoves = () => (
  <MockProviders>
    <MultiMovesMoveContainer moves={mockMovesPCS.previousMoves} />
  </MockProviders>
);

export const RetirementCurrentMove = () => (
  <MockProviders>
    <MultiMovesMoveContainer moves={mockMovesRetirement.currentMove} />
  </MockProviders>
);

export const RetirementPreviousMoves = () => (
  <MockProviders>
    <MultiMovesMoveContainer moves={mockMovesRetirement.previousMoves} />
  </MockProviders>
);

export const SeparationCurrentMove = () => (
  <MockProviders>
    <MultiMovesMoveContainer moves={mockMovesSeparation.currentMove} />
  </MockProviders>
);

export const SeparationPreviousMoves = () => (
  <MockProviders>
    <MultiMovesMoveContainer moves={mockMovesSeparation.previousMoves} />
  </MockProviders>
);
