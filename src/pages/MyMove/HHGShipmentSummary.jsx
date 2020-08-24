import React from 'react';
import PropTypes, { string } from 'prop-types';
import { get, isEmpty } from 'lodash';
import { GridContainer, Grid } from '@trussworks/react-uswds';

import Address from 'scenes/Review/Address';
import { formatDateSM } from 'shared/formatters';
import { getFullAgentName } from 'utils/moveSetupFlow';
import { MTOAgentType } from 'shared/constants';

import 'scenes/Review/Review.css';

export default function HHGShipmentSummary(props) {
  const { mtoShipment } = props;

  const requestedPickupDate = get(mtoShipment, 'requestedPickupDate', '');
  const pickupLocation = get(mtoShipment, 'pickupAddress', {});
  const agents = get(mtoShipment, 'agents', {});
  const releasingAgent = Object.values(agents).find((agent) => agent.agentType === MTOAgentType.RELEASING);

  const requestedDeliveryDate = get(mtoShipment, 'requestedDeliveryDate', '');
  const dropoffLocation = get(mtoShipment, 'destinationAddress', {});
  const receivingAgent = Object.values(agents).find((agent) => agent.agentType === MTOAgentType.RECEIVING);
  const remarks = get(mtoShipment, 'customerRemarks', '');

  // <ReviewSection fieldData={hhgShipmentData} title="Shipment [x]: HHG" editLink='' />

  /*
const hhgShipmentData = [
  { label: 'Expected departure', value: },
  { label: 'Starting ZIP', value: },
  { label: 'Storage (SIT)', value: }, // val: "No" || "Yes, [x] days"
  { label: 'Destination ZIP', value: },
  { label: 'HHG shipment weight' }, // subheading, needs to be editable, but has the same styling as the rest...
  { label: 'Estimated weight', value: },
  { label: 'Estimated incentive', value: }, // val: [x] lbs
  { label: 'Requested pickup date', value: },
  { label: 'Pickup location', value: },
  { label: 'Releasing agent', value: },
  { label: 'Requested delivery date', value: },
  { label: 'Destination', value: },
  { label: 'Receiving agent', value: },
  { label: 'Remarks', value: },
  { label: '', value: },
];
*/

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
              <div className="review-section">
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
