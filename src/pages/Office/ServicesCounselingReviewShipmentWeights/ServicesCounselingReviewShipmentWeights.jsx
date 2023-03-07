import React from 'react';

import { useReviewShipmentWeightsQuery } from 'hooks/queries';

const ServicesCounselingReviewShipmentWeights = ({ moveCode }) => {
  useReviewShipmentWeightsQuery(moveCode);
  return <h1>Review shipment weights</h1>;
};

export default ServicesCounselingReviewShipmentWeights;
