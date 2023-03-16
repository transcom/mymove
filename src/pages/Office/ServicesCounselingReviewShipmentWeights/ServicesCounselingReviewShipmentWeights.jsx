import React from 'react';
import PropTypes from 'prop-types';
import { Grid, GridContainer, Tag, Alert } from '@trussworks/react-uswds';
import { Link, generatePath } from 'react-router-dom';

import tabStyles from '../TXOMoveInfo/TXOTab.module.scss';

import styles from './ServicesCounselingReviewShipmentWeights.module.scss';

import { useReviewShipmentWeightsQuery } from 'hooks/queries';
import WeightDisplay from 'components/Office/WeightDisplay/WeightDisplay';
import { calculateEstimatedWeight, calculateWeightRequested } from 'hooks/custom';
import hasRiskOfExcess from 'utils/hasRiskOfExcess';
import { servicesCounselingRoutes } from 'constants/routes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const ServicesCounselingReviewShipmentWeights = ({ moveCode }) => {
  const { orders, mtoShipments, isLoading, isError } = useReviewShipmentWeightsQuery(moveCode);
  const estimatedWeightTotal = calculateEstimatedWeight(mtoShipments);
  const moveWeightTotal = calculateWeightRequested(mtoShipments);
  const externalVendorShipmentCount = mtoShipments?.length
    ? mtoShipments.filter((shipment) => shipment.usesExternalVendor).length
    : 0;
  const order = Object.values(orders)?.[0];

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <div className={tabStyles.tabContent}>
      <GridContainer>
        <Grid className={styles.alertContainer}>
          <Alert headingLevel="h4" slim type="warning">
            <span>This move has excess weight. Review PPM weight ticket documents to resolve.</span>
          </Alert>
        </Grid>
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

ServicesCounselingReviewShipmentWeights.propTypes = {
  moveCode: PropTypes.string.isRequired,
};

export default ServicesCounselingReviewShipmentWeights;
