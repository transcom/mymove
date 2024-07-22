import React from 'react';
import { func, string, number, bool, node, oneOfType } from 'prop-types';
import classnames from 'classnames';

import EditBillableWeight from '../EditBillableWeight/EditBillableWeight';

import styles from './ShipmentCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import { formatWeight, formatDateFromIso, formatAddressShort } from 'utils/formatters';
import { shipmentIsOverweight } from 'utils/shipmentWeights';
import { MandatorySimpleAddressShape, SimpleAddressShape } from 'types/address';
import { ShipmentOptionsOneOf } from 'types/shipment';
import { shipmentTypeLabels } from 'content/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import returnLowestValue from 'utils/returnLowestValue';

const ShipmentCardDetailRow = ({ display, rowTestId, className, title, content, contentTestId }) => {
  if (display) {
    return (
      <div data-testid={rowTestId} className={className}>
        <strong>{title}</strong>
        <span data-testid={contentTestId}>{content}</span>
      </div>
    );
  }

  return null;
};

ShipmentCardDetailRow.propTypes = {
  display: bool,
  rowTestId: string,
  className: string,
  title: string,
  content: oneOfType([node, string]),
  contentTestId: string,
};

ShipmentCardDetailRow.defaultProps = {
  display: true,
  rowTestId: '',
  className: '',
  title: '',
  content: '',
  contentTestId: '',
};

export default function ShipmentCard({
  billableWeight,
  billableWeightJustification,
  dateReweighRequested,
  departedDate,
  editEntity,
  pickupAddress,
  destinationAddress,
  estimatedWeight,
  primeActualWeight,
  originalWeight,
  adjustedWeight,
  reweighRemarks,
  reweighWeight,
  maxBillableWeight,
  totalBillableWeight,
  shipmentType,
  storageFacilityAddress,
}) {
  let showOriginalWeightHighlight = false;
  let showReweighWeightHighlight = false;

  // shipment weight exceeds 110% estimated weight
  // no need to show yellow highlight if adjusted weight (billable weight cap) exists
  if (estimatedWeight && !adjustedWeight) {
    // reweigh and original weight available
    // determine if yellow highlight needs to show if reweigh weight is over weight
    if (reweighWeight && originalWeight && reweighWeight <= originalWeight) {
      // reweigh weight is the shipment weight
      showReweighWeightHighlight = shipmentIsOverweight(estimatedWeight, reweighWeight);
    } else {
      // original weight is the shipment weight
      showOriginalWeightHighlight = shipmentIsOverweight(estimatedWeight, originalWeight);
    }
  }

  // reweigh requested and missing weight, show yellow highlight
  if (dateReweighRequested && !reweighWeight) {
    showReweighWeightHighlight = true;
  }

  const shipmentIsNTS = shipmentType === SHIPMENT_OPTIONS.NTS;
  const shipmentIsNTSR = shipmentType === SHIPMENT_OPTIONS.NTSR;
  const dateText = shipmentIsNTSR ? 'Delivered' : 'Departed';

  const lowestWeight = formatWeight(returnLowestValue(primeActualWeight, reweighWeight));

  const originAddress = shipmentIsNTSR ? storageFacilityAddress : pickupAddress;
  const deliveryAddress = shipmentIsNTS ? storageFacilityAddress : destinationAddress;

  return (
    <ShipmentContainer shipmentType={shipmentType} className={styles.container}>
      <header>
        <h2>{shipmentTypeLabels[shipmentType]}</h2>
        <section>
          <span>
            <strong>{dateText}</strong>
            <span data-testid="departureDate"> {formatDateFromIso(departedDate, 'DD MMM YYYY')}</span>
          </span>
          <span>
            <strong>From</strong> {originAddress && formatAddressShort(originAddress)}
          </span>
          <span>
            <strong>To</strong> {deliveryAddress && formatAddressShort(deliveryAddress)}
          </span>
        </section>
      </header>
      <div className={styles.weights}>
        <ShipmentCardDetailRow
          rowTestId="estimatedWeightContainer"
          className={classnames(styles.field, {
            [styles.warning]: !estimatedWeight,
          })}
          title="Estimated weight"
          contentTestId="estimatedWeight"
          content={estimatedWeight ? formatWeight(estimatedWeight) : <strong>Missing</strong>}
        />

        <ShipmentCardDetailRow
          rowTestId="actualWeightContainer"
          className={classnames(styles.field, {
            [styles.warning]: !primeActualWeight,
          })}
          title="Actual weight"
          contentTestId="actualWeight"
          content={primeActualWeight ? lowestWeight : <strong>Missing</strong>}
        />

        <ShipmentCardDetailRow
          rowTestId="originalWeightContainer"
          className={classnames(styles.field, {
            [styles.warning]: showOriginalWeightHighlight,
          })}
          title="Original weight"
          contentTestId="originalWeight"
          content={formatWeight(originalWeight)}
        />

        <ShipmentCardDetailRow
          display={!!dateReweighRequested}
          rowTestId="reweighWeightContainer"
          className={classnames(styles.field, {
            [styles.warning]: showReweighWeightHighlight,
          })}
          title="Reweigh weight"
          contentTestId="reweighWeight"
          content={reweighWeight ? formatWeight(reweighWeight) : <strong>Missing</strong>}
        />

        <ShipmentCardDetailRow
          display={!!dateReweighRequested}
          className={reweighRemarks ? styles.field : classnames(styles.field, styles.lastRow)}
          title="Date reweigh requested"
          contentTestId="dateReweighRequested"
          content={formatDateFromIso(dateReweighRequested, 'DD MMM YYYY')}
        />

        <ShipmentCardDetailRow
          display={!!dateReweighRequested && !!reweighRemarks}
          className={classnames(styles.field, styles.remarks, styles.lastRow)}
          title="Reweigh remarks"
          contentTestId="reweighRemarks"
          content={reweighRemarks}
        />
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
          isNTSRShipment={shipmentType === SHIPMENT_OPTIONS.NTSR}
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
  destinationAddress: MandatorySimpleAddressShape.isRequired,
  editEntity: func.isRequired,
  estimatedWeight: number,
  primeActualWeight: number,
  originalWeight: number.isRequired,
  adjustedWeight: number,
  pickupAddress: MandatorySimpleAddressShape.isRequired,
  reweighRemarks: string,
  reweighWeight: number,
  maxBillableWeight: number.isRequired,
  totalBillableWeight: number,
  shipmentType: ShipmentOptionsOneOf.isRequired,
  storageFacilityAddress: SimpleAddressShape,
};

ShipmentCard.defaultProps = {
  billableWeight: 0,
  billableWeightJustification: '',
  dateReweighRequested: '',
  estimatedWeight: 0,
  primeActualWeight: 0,
  adjustedWeight: null,
  reweighWeight: null,
  reweighRemarks: '',
  totalBillableWeight: 0,
  storageFacilityAddress: {},
};
