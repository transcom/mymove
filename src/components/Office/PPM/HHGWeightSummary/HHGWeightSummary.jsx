import { React } from 'react';
import classnames from 'classnames';
import { Label } from '@trussworks/react-uswds';
import { PropTypes } from 'prop-types';

import headerSectionStyles from '../PPMHeaderSummary/HeaderSection.module.scss';

import styles from './HHGWeightSummary.module.scss';

import { DEFAULT_EMPTY_VALUE } from 'shared/constants';
import { ShipmentShape } from 'types/shipment';
import { formatWeight } from 'utils/formatters';

const getHHGShipments = (mtoShipments) => {
  return mtoShipments.filter((shipment) => shipment.shipmentType === 'HHG');
};

const HHGSummaryBlock = ({ hhgNumber, estimatedWeight, actualWeight }) => {
  return (
    <div className={styles.header}>
      <h3>HHG {hhgNumber}</h3>
      <section>
        <div className={classnames(headerSectionStyles.Details)}>
          <div>
            <Label>Estimated Weight</Label>
            <span className={headerSectionStyles.light}>
              {estimatedWeight ? formatWeight(estimatedWeight) : DEFAULT_EMPTY_VALUE}
            </span>
          </div>
          <div>
            <Label>Actual Weight</Label>
            <span className={headerSectionStyles.light}>
              {actualWeight ? formatWeight(actualWeight) : DEFAULT_EMPTY_VALUE}
            </span>
          </div>
        </div>
      </section>
    </div>
  );
};
export default function HHGWeightSummary({ mtoShipments }) {
  const hhgShipments = getHHGShipments(mtoShipments);
  if (hhgShipments.length === 0) return null;

  return (
    <header className={classnames(styles.HHGWeightSummary)}>
      <div className={classnames(styles.HHGContainer)}>
        {getHHGShipments(mtoShipments).map((shipment, idx) => {
          return (
            <HHGSummaryBlock
              key={`hhg-summary-${idx}`}
              hhgNumber={idx + 1}
              estimatedWeight={shipment.primeEstimatedWeight}
              actualWeight={shipment.primeActualWeight}
            />
          );
        })}
        <hr />
      </div>
    </header>
  );
}

HHGWeightSummary.propTypes = {
  mtoShipments: PropTypes.arrayOf(ShipmentShape),
};

HHGWeightSummary.defaultProps = {
  mtoShipments: [],
};
