import React from 'react';
import PropTypes from 'prop-types';
import getColor from 'helpers/colors';

const ImportantShipmentDates = ({ requestedPickupDate, scheduledPickupDate }) => {
  return (
    <div
      style={{ backgroundColor: getColor('base-baselightest'), padding: '0 1.6rem 1rem 1.6rem' }}
      className="container maxw-tablet"
    >
      <div style={{ display: 'flex' }}>
        <p
          className="display-inline-block margin-bottom-0"
          style={{ fontSize: '0.87rem', fontWeight: 'bold', color: getColor('base-basedark'), width: '50%' }}
        >
          Customer requested pick up date
        </p>
        <p
          className="display-inline-block margin-bottom-0"
          style={{ fontSize: '0.87rem', fontWeight: 'bold', color: getColor('base-basedark') }}
        >
          Scheduled pick up date
        </p>
      </div>

      <hr style={{ border: `1px solid ${getColor('base-baselighter')}` }} />

      <div style={{ display: 'flex' }}>
        <p className="display-inline-block" style={{ width: '50%', marginTop: '0' }}>
          {requestedPickupDate}
        </p>
        <p className="display-inline-block" style={{ width: '50%', marginTop: '0' }}>
          {scheduledPickupDate}
        </p>
      </div>
    </div>
  );
};

ImportantShipmentDates.defaultProps = {
  requestedPickupDate: '',
  scheduledPickupDate: '',
};

ImportantShipmentDates.propTypes = {
  requestedPickupDate: PropTypes.string,
  scheduledPickupDate: PropTypes.string,
};

export default ImportantShipmentDates;
