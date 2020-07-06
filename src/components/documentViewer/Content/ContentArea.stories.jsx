import React from 'react';

import DocViewerContent from './Content';

export default {
  title: 'Components|Document Viewer',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/8f32f4ab-cbe5-45f7-a0df-4ad19c0902a8?mode=design',
    },
  },
};

export const ContentArea = () => <DocViewerContent />;
