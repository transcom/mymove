import React from 'react';

import CertificationText from './CertificationText';

import { completeCertificationText } from 'scenes/Legalese/legaleseText';

export default {
  title: 'Customer Components / Forms / Certification Text',
  component: CertificationText,
  argTypes: {
    certificationText: '',
    onScrollToBottom: { action: 'scroll to bottom' },
  },
};

export const DefaultState = (argTypes) => (
  <CertificationText onScrollToBottom={argTypes.onScrollToBottom} certificationText={completeCertificationText} />
);

export const WithServerError = (argTypes) => (
  <CertificationText onScrollToBottom={argTypes.onScrollToBottom} certificationText={completeCertificationText} />
);

export const LoadingCertificationText = (argTypes) => (
  <CertificationText onScrollToBottom={argTypes.onScrollToBottom} certificationText={completeCertificationText} />
);
