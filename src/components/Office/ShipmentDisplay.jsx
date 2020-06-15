/* eslint-disable jsx-a11y/label-has-associated-control */
/* eslint-disable react/jsx-wrap-multilines */
import React from 'react';
import classNames from 'classnames/bind';
import { SHIPMENT_TYPE } from 'shared/constants';
import * as PropTypes from 'prop-types';
import styles from 'components/Office/ShipmentDisplay.module.scss';
import { formatDate } from 'shared/dates';
import { Checkbox } from '@trussworks/react-uswds';
import { ReactComponent as ChevronDown } from '../../shared/icon/chevron-down.svg';
import ShipmentContainer from './ShipmentContainer';

const cx = classNames.bind(styles);

const ShipmentDisplay = ({ shipmentType, displayInfo, onChange, checked, shipmentId }) => {
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
                  checked={checked}
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
  onChange: PropTypes.func,
  shipmentId: PropTypes.string.isRequired,
  checked: PropTypes.bool,
  shipmentType: PropTypes.oneOf([
    SHIPMENT_TYPE.HHG,
    SHIPMENT_TYPE.HHG_SHORTHAUL_DOMESTIC,
    SHIPMENT_TYPE.HHG_LONGHAUL_DOMESTIC,
    SHIPMENT_TYPE.NTS,
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
  shipmentType: SHIPMENT_TYPE.HHG,
};

export default ShipmentDisplay;
