import React from 'react';
import { func, string, number, shape } from 'prop-types';
import classnames from 'classnames';

import EditBillableWeight from '../EditBillableWeight/EditBillableWeight';

import styles from './ShipmentCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatWeight, formatAddressShort, formatDateFromIso } from 'shared/formatters';
import { shipmentIsOverweight } from 'utils/shipmentWeights';

export default function ShipmentCard({
  billableWeight,
  billableWeightJustification,
  dateReweighRequested,
  departedDate,
  editEntity,
  pickupAddress,
  destinationAddress,
  estimatedWeight,
  originalWeight,
  reweighRemarks,
  reweighWeight,
  maxBillableWeight,
  totalBillableWeight,
}) {
  return (
    <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.HHG} className={styles.container}>
      <header>
        <h2>HHG</h2>
        <section>
          <span>
            <strong>Departed</strong>
            <span data-testid="departureDate"> {formatDateFromIso(departedDate, 'DD MMM YYYY')}</span>
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
            [styles.missing]: !shipmentIsOverweight(estimatedWeight, billableWeight) && estimatedWeight,
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
            <div className={reweighRemarks ? styles.field : classnames(styles.field, styles.lastRow)}>
              <strong>Date reweigh requested</strong>
              <span data-testid="dateReweighRequested">{formatDateFromIso(dateReweighRequested, 'DD MMM YYYY')}</span>
            </div>
            {reweighRemarks && (
              <div className={classnames(styles.field, styles.remarks, styles.lastRow)}>
                <strong>Reweigh remarks</strong>
                <span data-testid="reweighRemarks">{reweighRemarks}</span>
              </div>
            )}
          </div>
        )}
      </div>
      <footer>
        <EditBillableWeight
          title="Billable weight"
          billableWeight={billableWeight}
          billableWeightJustification={billableWeightJustification}
          originalWeight={originalWeight}
          estimatedWeight={estimatedWeight}
          editEntity={editEntity}
          maxBillableWeight={maxBillableWeight}
          totalBillableWeight={totalBillableWeight}
        />
      </footer>
    </ShipmentContainer>
  );
}

ShipmentCard.propTypes = {
  billableWeight: number,
  billableWeightJustification: string,
  dateReweighRequested: string,
  departedDate: string.isRequired,
  destinationAddress: shape({
    city: string.isRequired,
    state: string.isRequired,
    postal_code: string.isRequired,
  }).isRequired,
  editEntity: func.isRequired,
  estimatedWeight: number,
  originalWeight: number.isRequired,
  pickupAddress: shape({
    city: string.isRequired,
    state: string.isRequired,
    postal_code: string.isRequired,
  }).isRequired,
  reweighRemarks: string,
  reweighWeight: number,
  maxBillableWeight: number.isRequired,
  totalBillableWeight: number,
};

ShipmentCard.defaultProps = {
  billableWeight: 0,
  billableWeightJustification: '',
  dateReweighRequested: '',
  estimatedWeight: 0,
  reweighWeight: null,
  reweighRemarks: '',
  totalBillableWeight: 0,
};
