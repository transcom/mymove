import React from 'react';
import PropTypes from 'prop-types';
import { get, isEmpty } from 'lodash';

import Address from './Address';
import { formatDateSM } from 'shared/formatters';
import { MTOAgentType } from 'shared/constants';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { getFullAgentName } from 'shared/formatters';
import './Review.css';

export default function HHGShipmentSummary(props) {
  const { mtoShipment } = props;

  const requestedPickupDate = get(mtoShipment, 'requestedPickupDate', '');
  const pickupLocation = get(mtoShipment, 'pickupAddress', '');
  const agents = get(mtoShipment, 'agents', {});
  const releasingAgent = Object.values(agents).find((agent) => agent.agentType === MTOAgentType.RELEASING);

  const requestedDeliveryDate = get(mtoShipment, 'requestedDeliveryDate', '');
  const dropoffLocation = get(mtoShipment, 'destinationAddress', ''); //OR newduty station
  const receivingAgent = Object.values(agents).find((agent) => agent.agentType === MTOAgentType.RECEIVING);
  const remarks = get(mtoShipment, 'customerRemarks', '');

  return (
    <div data-testid="hhg-summary" className="review-content">
      <GridContainer>
        <h3>Shipment - Government moves all of your stuff (HHG)</h3>
        <Grid row>
          <Grid tablet={{ col: true }}>
            <div className="review-section">
              <p className="heading">Pickup Dates & Locations</p>
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
                        <td>{getFullAgentName(releasingAgent)}</td>
                      </tr>
                    )}
                  </tbody>
                )}
              </table>
            </div>
          </Grid>
          <Grid tablet={{ col: true }}>
            <div className="review-section">
              <p className="heading">Delivery Dates & Locations</p>
              <table>
                {!isEmpty(mtoShipment) && (
                  <tbody>
                    <tr>
                      <td>Requested Delivery Date: </td>
                      <td>{formatDateSM(requestedDeliveryDate)}</td>
                    </tr>
                    <tr>
                      <td>Drop-off Location: </td>
                      <td>
                        <Address address={dropoffLocation} />
                      </td>
                    </tr>
                    {releasingAgent && (
                      <tr>
                        <td>Receiving Agent:</td>
                        <td>{getFullAgentName(receivingAgent)}</td>
                      </tr>
                    )}
                  </tbody>
                )}
              </table>
            </div>
          </Grid>
        </Grid>
        <Grid row>
          <Grid tablet={{ col: true }}>
            <div className="review-section">
              <p className="heading">Customer Remarks</p>
              <table>
                {!isEmpty(mtoShipment) && (
                  <tbody>
                    {remarks !== '' && (
                      <tr>
                        <td>Notes:</td>
                        <td>{remarks}</td>
                      </tr>
                    )}
                  </tbody>
                )}
              </table>
            </div>
          </Grid>
          <Grid tablet={{ col: true }}></Grid>
        </Grid>
      </GridContainer>
    </div>
  );
}

HHGShipmentSummary.propTypes = {
  mtoShipment: PropTypes.object.isRequired,
  entitlements: PropTypes.object.isRequired,
};
