import React from 'react';
import { string, arrayOf, func, shape, number, bool } from 'prop-types';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './BillableWeightCard.module.scss';

import ExternalVendorWeightSummary from 'components/Office/ExternalVendorWeightSummary/ExternalVendorWeightSummary';
import ShipmentList from 'components/ShipmentList/ShipmentList';
import { formatWeight } from 'utils/formatters';

export default function BillableWeightCard({
  maxBillableWeight,
  weightRequested,
  weightAllowance,
  totalBillableWeight,
  shipments,
  onReviewWeights,
  secondaryReviewWeightsBtn,
}) {
  return (
    <div className={classnames(styles.cardContainer, 'container')}>
      <div className={styles.cardHeader}>
        <div>
          <h2>Billable weights</h2>
          {totalBillableWeight > maxBillableWeight && (
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
        </div>
        <Button onClick={onReviewWeights} secondary={secondaryReviewWeightsBtn} style={{ maxWidth: '160px' }}>
          Review weights
        </Button>
      </div>
      <div className={styles.spaceBetween}>
        <div>
          <h5>Maximum billable weight</h5>
          <h4>{formatWeight(maxBillableWeight)}</h4>
          <h6>
            Weight requested<strong>{formatWeight(weightRequested)}</strong>
          </h6>
          <h6>
            Weight allowance<strong>{formatWeight(weightAllowance)}</strong>
          </h6>
          {shipments.some((s) => s.usesExternalVendor) && (
            <ExternalVendorWeightSummary shipments={shipments.filter((s) => s.usesExternalVendor)} />
          )}
        </div>
        <div className={styles.shipmentSection}>
          <h5>Total billable weight</h5>
          <h4>{formatWeight(totalBillableWeight)}</h4>
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
  totalBillableWeight: number,
  onReviewWeights: func.isRequired,
  secondaryReviewWeightsBtn: bool.isRequired,
  shipments: arrayOf(
    shape({
      id: string.isRequired,
      shipmentType: string.isRequired,
      reweigh: shape({ id: string.isRequired, weight: number }),
    }),
  ).isRequired,
};

BillableWeightCard.defaultProps = {
  weightRequested: null,
  totalBillableWeight: null,
};
