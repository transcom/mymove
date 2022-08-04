import React from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';

import styles from 'components/Office/SITCostDetails/SITCostDetails.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { formatCentsTruncateWhole, formatDaysInTransit, formatWeight } from 'utils/formatters';

const SITCostDetails = ({ cost, weight, location, sitLocation, departureDate, entryDate }) => {
  const days = moment(departureDate).diff(moment(entryDate), 'days');
  return (
    <SectionWrapper className={styles.SITCostDetails}>
      <h2>Storage in transit (SIT)</h2>
      <h3 className={styles.NoSpacing}>{`Government constructed cost: $${formatCentsTruncateWhole(cost)}`}</h3>
      <p>
        {`Maximum reimbursement for storing ${formatWeight(weight)} of ${sitLocation.toLowerCase()} SIT
        at ${location} for ${formatDaysInTransit(days)}.`}
      </p>
    </SectionWrapper>
  );
};

SITCostDetails.propTypes = {
  cost: PropTypes.number.isRequired,
  weight: PropTypes.number.isRequired,
  sitLocation: PropTypes.string.isRequired,
  location: PropTypes.string.isRequired,
  departureDate: PropTypes.string.isRequired,
  entryDate: PropTypes.string.isRequired,
};

export default SITCostDetails;
