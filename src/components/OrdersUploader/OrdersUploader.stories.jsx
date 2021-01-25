import React from 'react';

import OrdersUploader from './index';

export default {
  title: 'Customer Components/Orders Uploader',
  component: OrdersUploader,
};

const mockCreateUploadSuccess = () => {
  return Promise.resolve();
};

const mockCreateUploadError = () => {
  return Promise.reject();
};

export const ordersUploaderSuccess = () => <OrdersUploader createUpload={mockCreateUploadSuccess} />;

export const ordersUploaderError = () => <OrdersUploader createUpload={mockCreateUploadError} />;
