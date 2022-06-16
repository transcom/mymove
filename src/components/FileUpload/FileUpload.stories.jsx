import React from 'react';

import FileUpload from './FileUpload';

export default {
  title: 'Components/FileUpload',
  component: FileUpload,
};

const mockCreateUploadSuccess = () => {
  return Promise.resolve({ id: '1234' });
};

const mockCreateUploadError = () => {
  return Promise.reject();
};
const Template = (args) => <FileUpload {...args} />;

export const FileUploadSuccess = Template.bind({});
FileUploadSuccess.args = {
  createUpload: mockCreateUploadSuccess,
};

export const FileUploadWithExtendedAcceptedFileTypes = Template.bind({});
FileUploadWithExtendedAcceptedFileTypes.args = {
  createUpload: mockCreateUploadSuccess,
  acceptedFileTypes: [
    'image/jpeg',
    'image/png',
    'application/pdf',
    '.csv',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.ms-excel',
  ],
};

export const FileUploadWithLimitOneFile = Template.bind({});
FileUploadWithLimitOneFile.args = {
  createUpload: mockCreateUploadSuccess,
  allowMultiple: false,
};

export const FileUploadWithNoParallelUploads = Template.bind({});
FileUploadWithNoParallelUploads.args = {
  createUpload: mockCreateUploadSuccess,
  maxParallelUploads: 1,
};

export const FileUploadError = Template.bind({});
FileUploadError.args = {
  createUpload: mockCreateUploadError,
};
