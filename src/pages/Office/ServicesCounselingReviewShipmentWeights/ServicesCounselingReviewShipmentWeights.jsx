import React, { useEffect, useState } from 'react';
import { Grid, GridContainer, Tag } from '@trussworks/react-uswds';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import { useReviewShipmentWeightsQuery } from 'hooks/queries';
import WeightDisplay from 'components/Office/WeightDisplay/WeightDisplay';
import { calculateEstimatedWeight, useCalculatedWeightRequested } from 'hooks/custom';
import hasRiskOfExcess from 'utils/hasRiskOfExcess';

const ServicesCounselingReviewShipmentWeights = ({ moveCode }) => {
  const { orders, mtoShipments } = useReviewShipmentWeightsQuery(moveCode);
  const [estimatedWeightTotal, setEstimatedWeightTotal] = useState(null);
  const moveWeightTotal = useCalculatedWeightRequested(mtoShipments);
  const order = Object.values(orders)?.[0];

  useEffect(() => {
    setEstimatedWeightTotal(calculateEstimatedWeight(mtoShipments));
  }, [mtoShipments]);

  return (
    <div className={styles.tabContent}>
      <GridContainer>
        <Grid row>
          <h1>Review shipment weights</h1>
        </Grid>
        <Grid row>
          <WeightDisplay heading="Weight allowance" weightValue={order.entitlement.totalWeight} />
          <WeightDisplay heading="Estimated weight (total)" weightValue={estimatedWeightTotal} />
          <WeightDisplay heading="Max billable weight" weightValue={order.entitlement.authorizedWeight}>
            {hasRiskOfExcess(estimatedWeightTotal, order.entitlement.totalWeight) && <Tag>Risk of excess</Tag>}
          </WeightDisplay>
          <WeightDisplay heading="Move weight (total)" weightValue={moveWeightTotal} />
        </Grid>
      </GridContainer>
    </div>
  );
};

export default ServicesCounselingReviewShipmentWeights;
