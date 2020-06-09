import React from 'react';
import classNames from 'classnames/bind';
import { SHIPMENT_TYPE } from 'shared/constants';
import * as PropTypes from 'prop-types';
import styles from 'components/Office/ShipmentDisplay.module.scss';
import { formatDate } from 'shared/dates';
import { ReactComponent as ChevronDown } from '../../shared/icon/chevron-down.svg';
import ShipmentContainer from './ShipmentContainer';

const cx = classNames.bind(styles);

const ShipmentDisplay = ({ shipmentType, checkboxId, displayInfo, onChange }) => {
  return (
    <div className={`${cx('shipment-display')}`} data-cy="shipment-display">
      <ShipmentContainer className={`${cx('shipment-display__container')}`} shipmentType={shipmentType}>
        <table className={`${cx('table--small')}`} data-cy="shipment-display-table">
          <thead>
            <tr>
              <th className={`${cx('shipment-display__header-checkbox')}`}>
                <div className="usa-checkbox">
                  <input
                    id={checkboxId || `shipment-display-checkbox-${shipmentType.toLowerCase()}`}
                    type="checkbox"
                    className="usa-checkbox__input"
                    data-cy="shipment-display-checkbox"
                    onChange={onChange}
                  />
                  {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
                  <label
                    className={checkboxId || `usa-checkbox__label ${cx('shipment-display__checkbox__label')}`}
                    htmlFor={`shipment-display-checkbox-${shipmentType.toLowerCase()}`}
                  />
                </div>
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
  checkboxId: PropTypes.string,
  onChange: PropTypes.func,
  shipmentType: PropTypes.oneOf([SHIPMENT_TYPE.HHG, SHIPMENT_TYPE.NTS]),
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
