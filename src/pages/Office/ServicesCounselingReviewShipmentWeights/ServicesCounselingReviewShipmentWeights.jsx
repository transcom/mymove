import React, { useEffect, useState } from 'react';
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
import ReviewShipmentWeightsTable from 'components/Office/PPM/ReviewShipmentWeightsTable/ReviewShipmentWeightsTable';
import {
  PPMReviewWeightsTableConfig,
  PPMReviewWeightsTableConfigWithoutGunSafe,
  nonPPMReviewWeightsTableConfig,
} from 'components/Office/PPM/ReviewShipmentWeightsTable/helpers';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { FEATURE_FLAG_KEYS, SHIPMENT_OPTIONS } from 'shared/constants';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const sortShipments = (shipments) => {
  const ppmShipment = [];
  const hhgShipment = [];
  if (!shipments) {
    return null;
  }
  shipments.forEach((shipment) => {
    if (shipment.shipmentType === SHIPMENT_OPTIONS.PPM) {
      ppmShipment.push(shipment);
      return;
    }
    hhgShipment.push(shipment);
  });
  return { hhgShipment, ppmShipment };
};

const ServicesCounselingReviewShipmentWeights = ({ moveCode }) => {
  const [showExcessWeightAlert, setShowExcessWeightAlert] = useState(false);
  const { orders, mtoShipments, isLoading, isError } = useReviewShipmentWeightsQuery(moveCode);
  const estimatedWeightTotal = calculateEstimatedWeight(mtoShipments);
  const moveWeightTotal = calculateWeightRequested(mtoShipments);
  const externalVendorShipmentCount = mtoShipments?.length
    ? mtoShipments.filter((shipment) => shipment.usesExternalVendor).length
    : 0;
  const order = Object.values(orders)?.[0];

  useEffect(() => {
    setShowExcessWeightAlert(moveWeightTotal > order.entitlement.totalWeight);
  }, [moveWeightTotal, order.entitlement.totalWeight]);

  const [isGunSafeEnabled, setIsGunSafeEnabled] = useState(false);
  useEffect(() => {
    const fetchData = async () => {
      setIsGunSafeEnabled(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.GUN_SAFE));
    };
    fetchData();
  }, []);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;
  // sort shipments for the table
  const sortedShipments = sortShipments(mtoShipments);
  const hasProGear = Boolean(order.entitlement?.proGearWeight || order.entitlement?.spouseProGearWeight);
  const showWeightsMoved = Boolean(hasProGear || sortedShipments.hhgShipment);

  return (
    <div className={tabStyles.tabContent}>
      <GridContainer>
        <Grid className={styles.alertContainer}>
          {showExcessWeightAlert && (
            <Alert headingLevel="h4" slim type="warning">
              <span>This move has excess weight. Review PPM weight ticket documents to resolve.</span>
            </Alert>
          )}
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
        {sortedShipments.ppmShipment && (
          <div className={styles.weightMovedContainer} data-testid="ppmShipmentContainer">
            <h2 className={styles.weightMovedHeader}>Weight moved by customer</h2>
            <ReviewShipmentWeightsTable
              tableData={sortedShipments.ppmShipment}
              tableConfig={isGunSafeEnabled ? PPMReviewWeightsTableConfig : PPMReviewWeightsTableConfigWithoutGunSafe}
            />
          </div>
        )}
        {showWeightsMoved && (
          <div className={styles.weightMovedContainer}>
            <h2 className={styles.weightMovedHeader}>Weight moved</h2>
            {sortedShipments?.hhgShipment?.length > 0 && (
              <div className={styles.shipmentContainer} data-testid="nonPpmShipmentContainer">
                <h3 className={styles.shipmentHeader}>Shipments</h3>
                <ReviewShipmentWeightsTable
                  tableData={sortedShipments.hhgShipment}
                  tableConfig={nonPPMReviewWeightsTableConfig}
                />
              </div>
            )}
          </div>
        )}
      </GridContainer>
    </div>
  );
};

ServicesCounselingReviewShipmentWeights.propTypes = {
  moveCode: PropTypes.string.isRequired,
};

export default ServicesCounselingReviewShipmentWeights;
