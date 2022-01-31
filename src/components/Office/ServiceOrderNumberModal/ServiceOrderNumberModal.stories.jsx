import React from 'react';

import ServiceOrderNumberModal from './ServiceOrderNumberModal';

export default {
  title: 'Office Components / ServiceOrderNumberModal',
  component: ServiceOrderNumberModal,
};

export const standard = () => {
  return (
    <div className="officeApp">
      <ServiceOrderNumberModal isOpen serviceOrderNumber="AB123456" />;
    </div>
  );
};
