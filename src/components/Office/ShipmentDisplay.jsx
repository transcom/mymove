/* eslint-disable jsx-a11y/label-has-associated-control */
/* eslint-disable react/jsx-wrap-multilines */
import React from 'react';
import * as PropTypes from 'prop-types';
import { Checkbox } from '@trussworks/react-uswds';

import { ReactComponent as ChevronDown } from '../../shared/icon/chevron-down.svg';

import ShipmentContainer from './ShipmentContainer';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import styles from 'components/Office/ShipmentDisplay.module.scss';
import { formatDate } from 'shared/dates';
import { ReactComponent as CheckmarkIcon } from 'shared/icon/checkbox--unchecked.svg';

const ShipmentDisplay = ({ shipmentType, displayInfo, onChange, shipmentId, isSubmitted }) => {
  return (
    <div className={styles['shipment-display']} data-testid="shipment-display">
      <ShipmentContainer className={styles['shipment-display__container']} shipmentType={shipmentType}>
        <table className="table--small" data-testid="shipment-display-table">
          <thead>
            <tr>
              <th className={styles['shipment-display__header-checkbox']}>
                {isSubmitted && (
                  <Checkbox
                    id={`shipment-display-checkbox-${shipmentId}`}
                    data-testid="shipment-display-checkbox"
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
                ,&nbsp;
                {`${displayInfo.currentAddress.city}, ${displayInfo.currentAddress.state} ${displayInfo.currentAddress.postal_code}`}
              </td>
              <td />
            </tr>
            <tr>
              <td />
              <td className={styles['shipment-display__label']}>Destination address</td>
              <td>
                {displayInfo.destinationAddress.street_address_1}
                , &nbsp;
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
  onChange: () => {},
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

export default ShipmentDisplay;
