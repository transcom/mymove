import React from 'react';
import { string, arrayOf, func, shape, number, bool } from 'prop-types';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './BillableWeightCard.module.scss';

import ShipmentModificationTag from 'components/ShipmentModificationTag/ShipmentModificationTag';
import ExternalVendorWeightSummary from 'components/Office/ExternalVendorWeightSummary/ExternalVendorWeightSummary';
import ShipmentList from 'components/ShipmentList/ShipmentList';
import { formatWeight } from 'utils/formatters';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import { shipmentModificationTypes } from 'constants/shipments';

export default function BillableWeightCard({
  maxBillableWeight,
  weightRequested,
  weightAllowance,
  actualBillableWeight,
  shipments,
  onReviewWeights,
  secondaryReviewWeightsBtn,
  isMoveLocked,
}) {
  const includesDivertedShipment = shipments.filter((s) => s.diversion).length > 0;

  return (
    <div className={classnames(styles.cardContainer, 'container')}>
      <div className={styles.cardHeader}>
        <div className={styles.cardTitleContainer}>
          <h2>Billable weights</h2>
          {actualBillableWeight > maxBillableWeight && (
            <div>
              <FontAwesomeIcon icon="exclamation-circle" className={styles.errorFlag} />
              <span
                data-testid="maxBillableWeightErrorText"
                className={classnames(styles.errorText, 'usa-error-message')}
              >
                Move exceeds max billable weight
              </span>
            </div>
          )}
          {includesDivertedShipment && (
            <ShipmentModificationTag shipmentModificationType={shipmentModificationTypes.DIVERSION} />
          )}
        </div>
        <Restricted to={permissionTypes.updateMaxBillableWeight}>
          <Button
            onClick={onReviewWeights}
            secondary={secondaryReviewWeightsBtn}
            style={{ maxWidth: '240px' }}
            disabled={!shipments.length > 0 || isMoveLocked}
          >
            Review shipment weights
          </Button>
        </Restricted>
      </div>
      <div className={styles.spaceBetween}>
        <div>
          <h5>Maximum billable weight</h5>
          <h4>{formatWeight(maxBillableWeight)}</h4>
          <h6>
            Actual weight<strong>{formatWeight(weightRequested)}</strong>
          </h6>
          <h6>
            Weight allowance<strong>{formatWeight(weightAllowance)}</strong>
          </h6>
          {shipments.some((s) => s.usesExternalVendor) && (
            <ExternalVendorWeightSummary shipments={shipments.filter((s) => s.usesExternalVendor)} />
          )}
        </div>
        <div className={styles.shipmentSection}>
          <h5>Actual billable weight</h5>
          <h4>{formatWeight(actualBillableWeight)}</h4>
          <div className={styles.shipmentList}>
            <ShipmentList shipments={shipments} showShipmentWeight moveSubmitted />
          </div>
        </div>
      </div>
    </div>
  );
}

BillableWeightCard.propTypes = {
  maxBillableWeight: number.isRequired,
  weightRequested: number,
  weightAllowance: number.isRequired,
  actualBillableWeight: number,
  onReviewWeights: func.isRequired,
  secondaryReviewWeightsBtn: bool.isRequired,
  shipments: arrayOf(
    shape({
      id: string.isRequired,
      shipmentType: string.isRequired,
      reweigh: shape({ id: string.isRequired, weight: number }),
    }),
  ).isRequired,
  isMoveLocked: bool,
};

BillableWeightCard.defaultProps = {
  weightRequested: null,
  actualBillableWeight: null,
  isMoveLocked: false,
};
