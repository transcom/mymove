import React from 'react';

import { SubmitMoveConfirmationModal } from './SubmitMoveConfirmationModal';

export default {
  title: 'Office Components/SubmitMoveConfirmationModal',
  component: SubmitMoveConfirmationModal,
};

export const Basic = () => <SubmitMoveConfirmationModal onClose={() => {}} onSubmit={() => {}} />;

export const AsApprovePPM = () => (
  <SubmitMoveConfirmationModal
    onClose={() => {}}
    onSubmit={() => {}}
    bodyText="You can't make changes after you approve the PPM."
  />
);
