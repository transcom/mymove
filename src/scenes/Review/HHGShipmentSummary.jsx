import React from 'react';
// import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';
import { get, isEmpty } from 'lodash';

// import truckIcon from 'shared/icon/truck-black.svg';

import Address from './Address';
import { formatDateSM } from 'shared/formatters';
import { MTOAgentType } from 'shared/constants';

import './Review.css';

export default function HHGShipmentSummary(props) {
  const { mtoShipment } = props;

  const requestedPickupDate = get(mtoShipment, 'requestedPickupDate', '');
  const pickupLocation = get(mtoShipment, 'pickupAddress', '');
  const agents = get(mtoShipment, 'agents', {});
  const releasingAgent = Object.values(agents).find((agent) => agent.agentType === MTOAgentType.RELEASING);
  console.log(releasingAgent);
  const requestedDeliveryDate = get(mtoShipment, 'requestedDeliveryDate', '');
  // const mtoShipment = Object.values(state.entities.mtoShipments).find(
  //   (mtoShipment) => mtoShipment.move_task_order_id === moveTaskOrderId,
  // );

  //   Dates and Locations

  // Requested pickup date: req

  // Pickup location: req

  // Releasing agent:

  // Requested delivery date: req

  // Drop-off location: (new duty station)

  // Receiving agent:

  // Remarks:
  return (
    <div className="usa-grid-full ppm-container hhg-shipment-summary">
      <h3>Shipment - Government moves all of your stuff (HHG)</h3>
      <div className="usa-width-one-half review-section ppm-review-sections">
        <div className="hhg-dates">
          <p className="heading">Dates and Locations</p>

          <table>
            {!isEmpty(mtoShipment) && (
              <tbody>
                <tr>
                  <td>Requested Pickup Date: </td>
                  <td>{formatDateSM(requestedPickupDate)}</td>
                </tr>
                <tr>
                  <td>Pickup Location: </td>
                  <td>
                    <Address address={pickupLocation} />
                  </td>
                </tr>
                {releasingAgent && (
                  <tr>
                    <td>Releasing Agent:</td>
                    <td>{releasingAgent.firstName}</td>
                  </tr>
                )}
                {requestedDeliveryDate !== '' && (
                  <tr>
                    <td>Requested Delivery Date:</td>
                    <td>{formatDateSM(requestedDeliveryDate)}</td>
                  </tr>
                )}
              </tbody>
            )}
          </table>
        </div>

        <p className="heading">
          Your Stuff
          <span className="not-implemented edit-section-link"> {/* <Link to={editWeightsPath}>Edit</Link> */}</span>
        </p>
      </div>

      <div className="usa-width-one-half review-section ppm-review-section">
        <p className="heading">
          Pickup &amp; Delivery Locations
          <span className="not-implemented edit-section-link"> {/* <Link to={editLocationsPath}>Edit</Link> */}</span>
        </p>
        <table>
          <tbody>
            <tr className="pickup-address">
              <td>Pickup Address:</td>
              <td>
                <Address address={get(mtoShipment, 'pickup_address', '')} />
              </td>
            </tr>
            {mtoShipment && mtoShipment.has_secondary_pickup_address && (
              <tr className="secondary-pickup-address">
                <td>Additional Pickup:</td>
                <td>
                  <Address address={get(mtoShipment, 'secondary_pickup_address', '')} />
                </td>
              </tr>
            )}
            {/* <tr className="delivery-address">
              <td>Delivery Address:</td>
              <td>
                {mtoShipment.has_delivery_address ? <Address address={mtoShipment.delivery_address} /> : 'Not entered'}
              </td>
            </tr> */}
          </tbody>
        </table>
        {!mtoShipment.has_delivery_address && (
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
  mtoShipment: PropTypes.object.isRequired,
  entitlements: PropTypes.object.isRequired,
};
