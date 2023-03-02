import React, { useEffect, useState } from 'react';
import { Grid, GridContainer, Tag } from '@trussworks/react-uswds';

import tabStyles from '../TXOMoveInfo/TXOTab.module.scss';

import styles from './ServicesCounselingReviewShipmentWeights.module.scss';

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
    <div className={tabStyles.tabContent}>
      <GridContainer>
        <Grid row>
          <h1>Review shipment weights</h1>
        </Grid>
        <div className={styles.weightHeader} id="move-weights">
          <WeightDisplay heading="Weight allowance" weightValue={order.entitlement.totalWeight} />
          <WeightDisplay heading="Estimated weight (total)" weightValue={estimatedWeightTotal} />
          <WeightDisplay heading="Max billable weight" weightValue={order.entitlement.authorizedWeight}>
            {hasRiskOfExcess(estimatedWeightTotal, order.entitlement.totalWeight) && <Tag>Risk of excess</Tag>}
          </WeightDisplay>
          <WeightDisplay heading="Move weight (total)" weightValue={moveWeightTotal} />
        </div>
      </GridContainer>
    </div>
  );
};

export default ServicesCounselingReviewShipmentWeights;
