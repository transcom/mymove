import React from 'react';

import { useReviewShipmentWeightsQuery } from 'hooks/queries';

const ServicesCounselingReviewShipmentWeights = ({ moveCode }) => {
  // eslint-disable-next-line no-unused-vars
  const { move, orders, mtoShipments, weightTickets, isLoading, isError, isSuccess } =
    useReviewShipmentWeightsQuery(moveCode);
  return <div />;
};

export default ServicesCounselingReviewShipmentWeights;
