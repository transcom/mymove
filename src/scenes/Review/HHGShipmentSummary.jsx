import React from 'react';
import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';

import truckIcon from 'shared/icon/truck-black.svg';

import Address from './Address';
import HHGWeightSummary from './HHGWeightSummary';
import HHGWeightWarning from './HHGWeightWarning';
import { displayDateRange } from 'shared/formatters';

import './Review.css';

export default function HHGShipmentSummary(props) {
  const { movePath, shipment, entitlements } = props;

  const editDatePath = movePath + '/edit-hhg-date';
  const editWeightsPath = movePath + '/edit-hhg-weights';
  const editLocationsPath = movePath + '/edit-hhg-locations';

  return (
    <div className="usa-grid-full ppm-container">
      <h3>
        <img src={truckIcon} alt="PPM shipment" /> Shipment - Government moves all of your stuff (HHG)
      </h3>
      <div className="usa-width-one-half review-section ppm-review-section">
        <p className="heading">
          Move Dates
          <span className="edit-section-link">
            {' '}
            <Link to={editDatePath}>Edit</Link>
          </span>
        </p>

        <table>
          <tbody>
            <tr>
              <td>Movers Packing: </td>
              <td>{displayDateRange(shipment.move_dates_summary.pack)}</td>
            </tr>
            <tr>
              <td>Loading Truck: </td>
              <td>{displayDateRange(shipment.move_dates_summary.pickup)}</td>
            </tr>
            <tr>
              <td>Move in Transit:</td>
              <td>{displayDateRange(shipment.move_dates_summary.transit)}</td>
            </tr>
            <tr>
              <td>Delivery:</td>
              <td>{displayDateRange(shipment.move_dates_summary.delivery, 'condensed')}</td>
            </tr>
          </tbody>
        </table>

        <p className="notice">
          Move dates are subject to change. Your mover will confirm final dates after your pre-move survey.
        </p>

        <p className="heading">
          Your Stuff
          <span className="edit-section-link">
            {' '}
            <Link to={editWeightsPath}>Edit</Link>
          </span>
        </p>

        <table>
          <tbody>
            <tr>
              <td>Weight Estimate:</td>
              <td>
                <HHGWeightSummary shipment={shipment} entitlements={entitlements} />
              </td>
            </tr>
            <tr>
              <td colSpan="2">
                <HHGWeightWarning shipment={shipment} entitlements={entitlements} />
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="usa-width-one-half review-section ppm-review-section">
        <p className="heading">
          Pickup &amp; Delivery Locations
          <span className="edit-section-link">
            {' '}
            <Link to={editLocationsPath}>Edit</Link>
          </span>
        </p>
        <table>
          <tbody>
            <tr className="pickup-address">
              <td>Pickup Address:</td>
              <td>
                <Address address={shipment.pickup_address} />
              </td>
            </tr>
            {shipment.has_secondary_pickup_address && (
              <tr className="secondary-pickup-address">
                <td>Additional Pickup:</td>
                <td>
                  <Address address={shipment.secondary_pickup_address} />
                </td>
              </tr>
            )}
            <tr className="delivery-address">
              <td>Delivery Address:</td>
              <td>{shipment.has_delivery_address ? <Address address={shipment.delivery_address} /> : 'Not entered'}</td>
            </tr>
          </tbody>
        </table>
        {!shipment.has_delivery_address && (
          <p className="notice delivery-notice">
            Note: If you don't have a delivery address before the movers arrive at your destination or you can't meet
            the delivery truck, the movers will put your stuff in temporary storage.
          </p>
        )}
      </div>
    </div>
  );
}

HHGShipmentSummary.propTypes = {
  shipment: PropTypes.object.isRequired,
  movePath: PropTypes.string.isRequired,
  entitlements: PropTypes.object.isRequired,
};
