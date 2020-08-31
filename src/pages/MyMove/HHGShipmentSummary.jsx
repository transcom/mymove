import React from 'react';
import { Link } from 'react-router-dom';
import PropTypes, { string } from 'prop-types';
import { get, isEmpty } from 'lodash';
import { GridContainer, Grid } from '@trussworks/react-uswds';

import Address from 'scenes/Review/Address';
import { formatDateSM } from 'shared/formatters';
import { getFullAgentName } from 'utils/moveSetupFlow';
import { MTOAgentType } from 'shared/constants';

import 'scenes/Review/Review.css';

export default function HHGShipmentSummary(props) {
  const { mtoShipment, movePath } = props;
  const editShipmentPath = `${movePath}/edit-shipment`;

  const requestedPickupDate = get(mtoShipment, 'requestedPickupDate', '');
  const pickupLocation = get(mtoShipment, 'pickupAddress', {});
  const agents = get(mtoShipment, 'agents', {});
  const releasingAgent = Object.values(agents).find((agent) => agent.agentType === MTOAgentType.RELEASING);

  const requestedDeliveryDate = get(mtoShipment, 'requestedDeliveryDate', '');
  const dropoffLocation = get(mtoShipment, 'destinationAddress', {});
  const receivingAgent = Object.values(agents).find((agent) => agent.agentType === MTOAgentType.RECEIVING);
  const remarks = get(mtoShipment, 'customerRemarks', '');

  return (
    <div data-testid="hhg-summary" className="review-content">
      <GridContainer>
        <h3>Shipment - Government moves all of your stuff (HHG)</h3>
        <span>
          <Link data-testid="edit-shipment" to={editShipmentPath}>
            Edit
          </Link>
        </span>
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
                    {!isEmpty(releasingAgent) && (
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
            {!isEmpty(mtoShipment) && (
              <div className="review-section">
                <p className="heading">Delivery Dates & Locations</p>
                <table>
                  <tbody>
                    <tr>
                      <td>Requested Delivery Date: </td>
                      {!isEmpty(requestedDeliveryDate) && <td>{formatDateSM(requestedDeliveryDate)}</td>}
                    </tr>

                    <tr>
                      <td>Drop-off Location: </td>
                      {!isEmpty(dropoffLocation) && (
                        <td>
                          <Address address={dropoffLocation} />
                        </td>
                      )}
                    </tr>

                    {!isEmpty(receivingAgent) && (
                      <tr>
                        <td>Receiving Agent:</td>
                        <td>{getFullAgentName(receivingAgent)}</td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            )}
          </Grid>
        </Grid>
        <Grid row>
          <Grid tablet={{ col: true }}>
            {!isEmpty(mtoShipment) && !isEmpty(remarks) && (
              <div data-testid="remarks" className="review-section">
                <p className="heading">Customer Remarks</p>
                <table>
                  <tbody>
                    <tr>
                      <td>Notes:</td>
                      <td>{remarks}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            )}
          </Grid>
          <Grid tablet={{ col: true }} />
        </Grid>
      </GridContainer>
    </div>
  );
}

HHGShipmentSummary.propTypes = {
  movePath: string.isRequired,
  mtoShipment: PropTypes.shape({
    agents: PropTypes.arrayOf(
      PropTypes.shape({
        firstName: string,
        lastName: string,
        agentType: string,
      }),
    ),
    customerRemarks: string,
    requestedPickupDate: string,
    requestedDeliveryDate: string,
    pickupAddress: PropTypes.shape({
      city: string,
      postal_code: string,
      state: string,
      street_address_1: string,
    }),
    destinationAddress: PropTypes.shape({
      city: string,
      postal_code: string,
      state: string,
      street_address_1: string,
    }),
  }),
};

HHGShipmentSummary.defaultProps = {
  mtoShipment: {},
};
