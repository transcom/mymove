import React from 'react';

import FileUpload from './index';

export default {
  title: 'Components/FileUpload',
  component: FileUpload,
};

const mockCreateUploadSuccess = () => {
  return Promise.resolve();
};

const mockCreateUploadError = () => {
  return Promise.reject();
};

export const fileUploadSuccess = () => <FileUpload createUpload={mockCreateUploadSuccess} />;

export const fileUploadError = () => <FileUpload createUpload={mockCreateUploadError} />;
