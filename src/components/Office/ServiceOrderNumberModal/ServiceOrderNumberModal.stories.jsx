import React from 'react';

import ServiceOrderNumberModal from './ServiceOrderNumberModal';

export default {
  title: 'Office Components / ServiceOrderNumberModal',
  component: ServiceOrderNumberModal,
};

export const standard = () => {
  return <ServiceOrderNumberModal isOpen serviceOrderNumber="AB123456" />;
};
