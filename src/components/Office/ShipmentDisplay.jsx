/* eslint-disable jsx-a11y/label-has-associated-control */
/* eslint-disable react/jsx-wrap-multilines */
import React from 'react';
import classNames from 'classnames/bind';
import * as PropTypes from 'prop-types';
import { Checkbox } from '@trussworks/react-uswds';

import { ReactComponent as ChevronDown } from '../../shared/icon/chevron-down.svg';

import ShipmentContainer from './ShipmentContainer';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import styles from 'components/Office/ShipmentDisplay.module.scss';
import { formatDate } from 'shared/dates';

const cx = classNames.bind(styles);

const ShipmentDisplay = ({ shipmentType, displayInfo, onChange, shipmentId }) => {
  return (
    <div className={`${cx('shipment-display')}`} data-cy="shipment-display">
      <ShipmentContainer className={`${cx('shipment-display__container')}`} shipmentType={shipmentType}>
        <table className={`${cx('table--small')}`} data-cy="shipment-display-table">
          <thead>
            <tr>
              <th className={`${cx('shipment-display__header-checkbox')}`}>
                <Checkbox
                  id={`shipment-display-checkbox-${shipmentId}`}
                  data-cy="shipment-display-checkbox"
                  onChange={onChange}
                  name="shipments"
                  label=""
                  value={shipmentId}
                />
              </th>
              <th>
                <h3 className={`${cx('shipment-display__heading')}`}>{displayInfo.heading}</h3>
              </th>
              <th> </th>
              <th className={`${cx('shipment-display__header-chevron-down')}`}>
                <ChevronDown />
              </th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td />
              <td className={`${cx('shipment-display__label')}`}>Requested move date</td>
              <td>{formatDate(displayInfo.requestedMoveDate, 'DD MMM YYYY')}</td>
              <td />
            </tr>
            <tr>
              <td />
              <td className={`${cx('shipment-display__label')}`}>Current address</td>
              <td>
                {displayInfo.currentAddress.street_address_1}
                <br />
                {`${displayInfo.currentAddress.city}, ${displayInfo.currentAddress.state} ${displayInfo.currentAddress.postal_code}`}
              </td>
              <td />
            </tr>
            <tr>
              <td />
              <td className={`${cx('shipment-display__label')}`}>Destination address</td>
              <td>
                {displayInfo.destinationAddress.street_address_1}
                <br />
                {`${displayInfo.destinationAddress.city}, ${displayInfo.destinationAddress.state} ${displayInfo.destinationAddress.postal_code}`}
              </td>
              <td />
            </tr>
          </tbody>
        </table>
      </ShipmentContainer>
    </div>
  );
};

ShipmentDisplay.propTypes = {
  onChange: PropTypes.func.isRequired,
  shipmentId: PropTypes.string.isRequired,
  shipmentType: PropTypes.oneOf([
    SHIPMENT_OPTIONS.HHG,
    SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.NTS,
  ]),
  displayInfo: PropTypes.shape({
    heading: PropTypes.string.isRequired,
    requestedMoveDate: PropTypes.string.isRequired,
    currentAddress: PropTypes.shape({
      street_address_1: PropTypes.string.isRequired,
      city: PropTypes.string.isRequired,
      state: PropTypes.string.isRequired,
      postal_code: PropTypes.string.isRequired,
    }).isRequired,
    destinationAddress: PropTypes.shape({
      street_address_1: PropTypes.string.isRequired,
      city: PropTypes.string.isRequired,
      state: PropTypes.string.isRequired,
      postal_code: PropTypes.string.isRequired,
    }).isRequired,
  }).isRequired,
};

ShipmentDisplay.defaultProps = {
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

export default ShipmentDisplay;
