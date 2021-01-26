import React from 'react';

import WizardNavigation from './WizardNavigation';

export default {
  title: 'Customer Components / Wizard Navigation',
  component: WizardNavigation,
};

export const stepStartPages = () => <WizardNavigation isFirstPage />;

export const stepStartPagesWithFinishLater = () => <WizardNavigation isFirstPage showFinishLater />;

export const midStepPages = () => <WizardNavigation />;

export const midStepPagesWithFinishLater = () => <WizardNavigation showFinishLater />;

export const fullFlowCompletion = () => <WizardNavigation isLastPage />;

export const fullFlowCompletionWithFinishLater = () => <WizardNavigation isLastPage showFinishLater />;

export const nextDisabled = () => <WizardNavigation disableNext />;

export const editPages = () => <WizardNavigation editMode />;

export const postSubmissionPages = () => <WizardNavigation readOnly />;
