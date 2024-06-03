import React from 'react';

import RequestShipmentDiversionModal from './RequestShipmentDiversionModal';

const shipmentInfo = {
  shipmentID: '123456',
  ifMatchEtag: 'string',
  shipmentLocator: '123456-01',
};

export default {
  title: 'Office Components/RequestShipmentDiversionModal',
  component: RequestShipmentDiversionModal,
};

export const Basic = () => (
  <RequestShipmentDiversionModal onClose={() => {}} onSubmit={() => {}} shipmentInfo={shipmentInfo} />
);
