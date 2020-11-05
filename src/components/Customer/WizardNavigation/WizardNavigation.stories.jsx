import React from 'react';

import WizardNavigation from './WizardNavigation';

export default {
  title: 'Customer Components | Wizard Navigation',
  component: WizardNavigation,
};

export const backAndNext = () => <WizardNavigation />;

export const firstPage = () => <WizardNavigation isFirstPage />;

export const lastPage = () => <WizardNavigation isLastPage />;

export const nextDisabled = () => <WizardNavigation disableNext />;

export const showFinishLater = () => <WizardNavigation showFinishLater />;
