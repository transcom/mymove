/*  react/jsx-props-no-spreading */
import React from 'react';

import DocsUploaded from '.';

const files = [{ filename: 'File 1' }, { filename: 'File 2' }, { filename: 'File 3' }];
export const Basic = () => (
  <div className="grid-container">
    <DocsUploaded files={files} />
  </div>
);

export default {
  title: 'Customer Components / DocsUploaded',
};
