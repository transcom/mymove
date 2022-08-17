import React from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';

import SectionWrapper from 'components/Customer/SectionWrapper';
import { formatCentsTruncateWhole, formatDaysInTransit, formatWeight } from 'utils/formatters';
import { LOCATION_TYPES } from 'types/sitStatusShape';

const SITCostDetails = ({ cost, weight, originZip, destinationZip, sitLocation, departureDate, entryDate }) => {
  const days = 1 + moment(departureDate).diff(moment(entryDate), 'days');
  const displaySitLocation = sitLocation.toLowerCase();
  const displayZip = sitLocation === LOCATION_TYPES.DESTINATION ? destinationZip : originZip;
  return (
    <SectionWrapper>
      <h2>Storage in transit (SIT)</h2>
      <h3 className="margin-bottom-0">{`Government constructed cost: $${formatCentsTruncateWhole(cost)}`}</h3>
      <p>
        {`${formatWeight(weight)} of ${displaySitLocation} SIT
        at ${displayZip} for ${formatDaysInTransit(days)}.`}
      </p>
    </SectionWrapper>
  );
};

SITCostDetails.propTypes = {
  cost: PropTypes.number.isRequired,
  weight: PropTypes.number.isRequired,
  sitLocation: PropTypes.oneOf([LOCATION_TYPES.DESTINATION, LOCATION_TYPES.ORIGIN]),
  departureDate: PropTypes.string.isRequired,
  entryDate: PropTypes.string.isRequired,
};

SITCostDetails.defaultProps = {
  sitLocation: LOCATION_TYPES.DESTINATION,
};

export default SITCostDetails;
