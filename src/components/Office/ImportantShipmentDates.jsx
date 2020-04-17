import React from 'react';
import classNames from 'classnames/bind';
import * as PropTypes from 'prop-types';
import styles from './importantShipmentDates.module.scss';

const cx = classNames.bind(styles);

const ImportantShipmentDates = ({ requestedPickupDate, scheduledPickupDate }) => {
  return (
    <div className={`container maxw-tablet ${cx('container-override')}`}>
      <div className={`${cx('flex')}`}>
        <p className={`${cx('header')} display-inline-block margin-bottom-0`}>Customer requested pick up date</p>
        <p className={`${cx('header')} display-inline-block margin-bottom-0`}>Scheduled pick up date</p>
      </div>

      <hr className={`${cx('divider')}`} />

      <div className={`${cx('flex')}`}>
        <p className={`${cx('date')} display-inline-block`}>{requestedPickupDate}</p>
        <p className={`${cx('date')} display-inline-block`}>{scheduledPickupDate}</p>
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
