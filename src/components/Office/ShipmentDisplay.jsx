/* eslint-disable jsx-a11y/label-has-associated-control */
/* eslint-disable react/jsx-wrap-multilines */
import React from 'react';
import * as PropTypes from 'prop-types';
import { Checkbox } from '@trussworks/react-uswds';

import { ReactComponent as ChevronDown } from '../../shared/icon/chevron-down.svg';

import ShipmentContainer from './ShipmentContainer';

import { SHIPMENT_TYPE } from 'shared/constants';
import styles from 'components/Office/ShipmentDisplay.module.scss';
import { formatDate } from 'shared/dates';
import { ReactComponent as CheckmarkIcon } from 'shared/icon/checkmark.svg';

const ShipmentDisplay = ({ shipmentType, displayInfo, onChange, shipmentId, isSubmitted }) => {
  return (
    <div className={styles['shipment-display']} data-cy="shipment-display">
      <ShipmentContainer className={styles['shipment-display__container']} shipmentType={shipmentType}>
        <table className="table--small" data-cy="shipment-display-table">
          <thead>
            <tr>
              <th className={styles['shipment-display__header-checkbox']}>
                {isSubmitted && (
                  <Checkbox
                    id={`shipment-display-checkbox-${shipmentId}`}
                    data-cy="shipment-display-checkbox"
                    onChange={onChange}
                    name="shipments"
                    label=""
                    value={shipmentId}
                  />
                )}
                {!isSubmitted && <CheckmarkIcon />}
              </th>
              <th>
                <h3 className={styles['shipment-display__heading']}>{displayInfo.heading}</h3>
              </th>
              <th> </th>
              <th className={styles['shipment-display__header-chevron-down']}>
                <ChevronDown />
              </th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td />
              <td className={styles['shipment-display__label']}>Requested move date</td>
              <td>{formatDate(displayInfo.requestedMoveDate, 'DD MMM YYYY')}</td>
              <td />
            </tr>
            <tr>
              <td />
              <td className={styles['shipment-display__label']}>Current address</td>
              <td>
                {displayInfo.currentAddress.street_address_1}
                <br />
                {`${displayInfo.currentAddress.city}, ${displayInfo.currentAddress.state} ${displayInfo.currentAddress.postal_code}`}
              </td>
              <td />
            </tr>
            <tr>
              <td />
              <td className={styles['shipment-display__label']}>Destination address</td>
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
  isSubmitted: PropTypes.bool.isRequired,
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
  onChange: () => {},
};

export default ShipmentDisplay;
