import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Label, Button, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './HeaderSection.module.scss';

import { formatDate, formatCentsTruncateWhole } from 'utils/formatters';

export const sectionTypes = {
  incentives: 'incentives',
  shipmentInfo: 'shipmentInfo',
  gcc: 'gcc',
};

const getSectionTitle = (sectionInfo) => {
  switch (sectionInfo.type) {
    case sectionTypes.incentives:
      return `Incentives/Costs`;
    case sectionTypes.shipmentInfo:
      return `Shipment Info`;
    case sectionTypes.gcc:
      return `GCC Factors`;
    default:
      return `Error getting section title!`;
  }
};

// Returns the markup needed for a specific section
const getSectionMarkup = (sectionInfo) => {
  switch (sectionInfo.type) {
    case sectionTypes.incentives:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label className={styles.headerLabel}>Gross Incentive</Label>
            <span className={styles.light}>{sectionInfo.grossIncentive ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Government Constructive Cost (GCC)</Label>
            <span className={styles.light}>{sectionInfo.gcc ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Advanced Operating Allowance</Label>
            <span className={styles.light}>{sectionInfo.aoa ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Remaining Reimbursement Owed to Customer</Label>
            <span className={styles.light}>{sectionInfo.remainingReimbursement ?? `TEST VAL`}</span>
          </div>
        </div>
      );
    case sectionTypes.shipmentInfo:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label className={styles.headerLabel}>Planned Move Start Date</Label>
            <span className={styles.light}>
              {sectionInfo.plannedMoveDate ? formatDate(sectionInfo.plannedMoveDate, null, 'DD-MMM-YYYY') : `TEST DATE`}
            </span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Actual Move Start Date</Label>
            <span className={styles.light}>{formatDate(sectionInfo.actualMoveDate, null, 'DD-MMM-YYYY')}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Starting ZIP</Label>
            <span className={styles.light}>{sectionInfo.actualPickupPostalCode}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Ending ZIP</Label>
            <span className={styles.light}>{sectionInfo.actualDestinationPostalCode}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Miles</Label>
            <span className={styles.light}>{sectionInfo.miles ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Estimated Net Weight</Label>
            <span className={styles.light}>{sectionInfo.estimatedNetWeight ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Actual Net Weight</Label>
            <span className={styles.light}>{sectionInfo.actualNetWeight ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Advance received</Label>
            <span className={styles.light}>
              {sectionInfo.hasReceivedAdvance
                ? `Yes, $${formatCentsTruncateWhole(sectionInfo.advanceAmountReceived)}`
                : 'No'}
            </span>
          </div>
        </div>
      );
    case sectionTypes.gcc:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label className={styles.headerLabel}>Base Linehaul</Label>
            <span className={styles.light}>{sectionInfo.baseLinehaul ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Origin Linehaul Factor</Label>
            <span className={styles.light}>{sectionInfo.originLinehaulFactor ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Destination Linehaul Factor</Label>
            <span className={styles.light}>{sectionInfo.destinationLinehaulFactor ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Linehaul Adjustment</Label>
            <span className={styles.light}>{sectionInfo.linehaulAdjustment ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>ShortHaul Charge</Label>
            <span className={styles.light}>{sectionInfo.shorthaulCharge ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Transportation Cost</Label>
            <span className={styles.light}>{sectionInfo.transportationCost ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Linehaul Fuel Surcharge</Label>
            <span className={styles.light}>{sectionInfo.linehaulFuelSurcharge ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Fuel Surcharge Percent</Label>
            <span className={styles.light}>{sectionInfo.fuelSurchargePercent ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Origin Service Area Fee</Label>
            <span className={styles.light}>{sectionInfo.originServiceAreaFee ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Origin Factor</Label>
            <span className={styles.light}>{sectionInfo.originFactor ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Destination Service Area Fee</Label>
            <span className={styles.light}>{sectionInfo.destinationServiceAreaFee ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Destination Factor</Label>
            <span className={styles.light}>{sectionInfo.destinationFactor ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Full Pack/Unpack Charge</Label>
            <span className={styles.light}>{sectionInfo.fullPackUnpackCharge ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>PPM Factor</Label>
            <span className={styles.light}>{sectionInfo.ppmFactor ?? `TEST VAL`}</span>
          </div>
        </div>
      );
    default:
      return <Alert>An error occured while getting section markup!</Alert>;
  }
};

export default function PPMHeaderSummary({ sectionInfo }) {
  const [showDetails, setShowDetails] = useState(false);
  const showRequestDetailsButton = true;
  const handleToggleDetails = () => setShowDetails((prevState) => !prevState);
  const showDetailsChevron = showDetails ? 'chevron-up' : 'chevron-down';
  const showDetailsText = showDetails ? 'Hide details' : 'Show details';

  return (
    <section className={classnames(styles.HeaderSection)}>
      <header>
        <h4>{getSectionTitle(sectionInfo)}</h4>
      </header>
      <div className={styles.toggleDrawer}>
        {showRequestDetailsButton && (
          <Button
            aria-expanded={showDetails}
            data-testid="showRequestDetailsButton"
            type="button"
            unstyled
            onClick={handleToggleDetails}
          >
            <FontAwesomeIcon icon={showDetailsChevron} /> {showDetailsText}
          </Button>
        )}
      </div>
      {showDetails && getSectionMarkup(sectionInfo)}
    </section>
  );
}
