import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import { useReviewShipmentWeightsQuery } from 'hooks/queries';
import WeightDisplay from 'components/Office/WeightDisplay/WeightDisplay';
import { useCalculatedEstimatedWeight, useCalculatedWeightRequested } from 'hooks/custom';

const ServicesCounselingReviewShipmentWeights = ({ moveCode }) => {
  const { orders, mtoShipments } = useReviewShipmentWeightsQuery(moveCode);
  const estimatedWeightTotal = useCalculatedEstimatedWeight(mtoShipments);
  const moveWeightTotal = useCalculatedWeightRequested(mtoShipments);
  const order = Object.values(orders)?.[0];

  return (
    <GridContainer>
      <Grid row>
        <h1>Review shipment weights</h1>
      </Grid>
      <Grid row>
        <WeightDisplay heading="Weight allowance" weightValue={order.entitlement.totalWeight} />
        <WeightDisplay heading="Estimated weight (total)" weightValue={estimatedWeightTotal} />
        <WeightDisplay heading="Max billable weight" weightValue={order.entitlement.authorizedWeight} />
        <WeightDisplay heading="Move weight (total)" weightValue={moveWeightTotal} />
      </Grid>
    </GridContainer>
  );
};

export default ServicesCounselingReviewShipmentWeights;
