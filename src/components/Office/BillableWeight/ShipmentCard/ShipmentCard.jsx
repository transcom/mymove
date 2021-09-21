import React from 'react';
import { string, number, shape } from 'prop-types';
import classnames from 'classnames';

import EditBillableWeight from '../EditBillableWeight/EditBillableWeight';

import styles from './ShipmentCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatWeight, formatAddressShort, formatDateFromIso } from 'shared/formatters';
import { shipmentIsOverweight } from 'utils/shipmentWeights';

export default function ShipmentCard({
  billableWeight,
  dateReweighRequested,
  departedDate,
  pickupAddress,
  destinationAddress,
  estimatedWeight,
  originalWeight,
  reweighRemarks,
  reweighWeight,
}) {
  return (
    <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.HHG} className={styles.container}>
      <header>
        <h2>HHG</h2>
        <section>
          <span>
            <strong>Departed</strong>
            <span data-testid="departureDate">{formatDateFromIso(departedDate, 'DD MMM YYYY')}</span>
          </span>
          <span>
            <strong>From</strong> {pickupAddress && formatAddressShort(pickupAddress)}
          </span>
          <span>
            <strong>To</strong> {destinationAddress && formatAddressShort(destinationAddress)}
          </span>
        </section>
      </header>
      <div className={styles.weights}>
        <div
          className={classnames(styles.field, {
            [styles.missing]: !estimatedWeight,
          })}
        >
          <strong>Estimated weight</strong>
          <span data-testid="estimatedWeight">
            {estimatedWeight ? formatWeight(estimatedWeight) : <strong>Missing</strong>}
          </span>
        </div>
        <div
          className={classnames(styles.field, {
            [styles.missing]: !shipmentIsOverweight(estimatedWeight, billableWeight),
          })}
        >
          <strong>Original weight</strong>
          <span data-testid="originalWeight">{formatWeight(originalWeight)}</span>
        </div>
        {dateReweighRequested && (
          <div>
            <div
              className={classnames(styles.field, {
                [styles.missing]: !reweighWeight,
              })}
            >
              <strong>Reweigh weight</strong>
              <span data-testid="reweighWeight">
                {reweighWeight ? formatWeight(reweighWeight) : <strong>Missing</strong>}
              </span>
            </div>
            <div className={styles.field}>
              <strong>Date reweigh requested</strong>
              <span data-testid="dateReweighRequested">{formatDateFromIso(dateReweighRequested, 'DD MMM YYYY')}</span>
            </div>
            <div className={classnames(styles.field, styles.remarks)}>
              <strong>Reweigh remarks</strong>
              <span data-testid="reweighRemarks">{reweighRemarks}</span>
            </div>
          </div>
        )}
      </div>
      <footer>
        <EditBillableWeight
          title="Billable weight"
          billableWeight={billableWeight}
          originalWeight={originalWeight}
          estimatedWeight={estimatedWeight}
        />
      </footer>
    </ShipmentContainer>
  );
}

ShipmentCard.propTypes = {
  billableWeight: number.isRequired,
  dateReweighRequested: string.isRequired,
  departedDate: string.isRequired,
  destinationAddress: shape({
    city: string.isRequired,
    state: string.isRequired,
    postal_code: string.isRequired,
  }).isRequired,
  estimatedWeight: number.isRequired,
  originalWeight: number.isRequired,
  pickupAddress: shape({
    city: string.isRequired,
    state: string.isRequired,
    postal_code: string.isRequired,
  }).isRequired,
  reweighRemarks: string.isRequired,
  reweighWeight: number,
};

ShipmentCard.defaultProps = {
  reweighWeight: null,
};
