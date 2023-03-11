import React, { useEffect, useState } from 'react';
import { Grid, GridContainer, Tag } from '@trussworks/react-uswds';
import { Link, generatePath } from 'react-router-dom';

import tabStyles from '../TXOMoveInfo/TXOTab.module.scss';

import styles from './ServicesCounselingReviewShipmentWeights.module.scss';

import { useReviewShipmentWeightsQuery } from 'hooks/queries';
import WeightDisplay from 'components/Office/WeightDisplay/WeightDisplay';
import { calculateEstimatedWeight, calculateWeightRequested } from 'hooks/custom';
import hasRiskOfExcess from 'utils/hasRiskOfExcess';
import { servicesCounselingRoutes } from 'constants/routes';

const ServicesCounselingReviewShipmentWeights = ({ moveCode }) => {
  const { orders, mtoShipments, documents } = useReviewShipmentWeightsQuery(moveCode);
  const [estimatedWeightTotal, setEstimatedWeightTotal] = useState(null);
  const [externalVendorShipmentCount, setExternalVendorShipmentCount] = useState(0);
  const [moveWeightTotal, setMoveWeightTotal] = useState(null);
  const order = Object.values(orders)?.[0];

  useEffect(() => {
    setEstimatedWeightTotal(calculateEstimatedWeight(mtoShipments));
  }, [mtoShipments]);

  // documents are a dependency of this useEffect because they are appended to the mtoShipments
  // to calculate the net weight, but this doesn't trigger the useEffect automatically
  useEffect(() => {
    setMoveWeightTotal(calculateWeightRequested(mtoShipments));
  }, [mtoShipments, documents]);

  useEffect(() => {
    if (mtoShipments) {
      const externalVendorShipments = mtoShipments?.length
        ? mtoShipments.filter((shipment) => shipment.usesExternalVendor).length
        : 0;
      setExternalVendorShipmentCount(externalVendorShipments);
    }
  }, [mtoShipments]);

  return (
    <div className={tabStyles.tabContent}>
      <GridContainer>
        <Grid row>
          <h1>Review shipment weights</h1>
        </Grid>
        <div className={styles.weightHeader} id="move-weights">
          <WeightDisplay heading="Weight allowance" weightValue={order.entitlement.totalWeight} />
          <WeightDisplay heading="Estimated weight (total)" weightValue={estimatedWeightTotal}>
            {hasRiskOfExcess(estimatedWeightTotal, order.entitlement.totalWeight) && <Tag>Risk of excess</Tag>}
            {hasRiskOfExcess(estimatedWeightTotal, order.entitlement.totalWeight) &&
              externalVendorShipmentCount > 0 && <br />}
            {externalVendorShipmentCount > 0 && (
              <small>
                {externalVendorShipmentCount} shipment{externalVendorShipmentCount > 1 && 's'} not moved by GHC prime.{' '}
                <Link className="usa-link" to={generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode })}>
                  View move details
                </Link>
              </small>
            )}
          </WeightDisplay>
          <WeightDisplay heading="Max billable weight" weightValue={order.entitlement.authorizedWeight} />
          <WeightDisplay heading="Move weight (total)" weightValue={moveWeightTotal} />
        </div>
      </GridContainer>
    </div>
  );
};

export default ServicesCounselingReviewShipmentWeights;
