import React from 'react';
import PropTypes, { string } from 'prop-types';
import { get, isEmpty } from 'lodash';

import ReviewSection from '../../components/Customer/ReviewSection';

import Address from 'scenes/Review/Address';
import { formatDateSM } from 'shared/formatters';
import { getFullAgentName } from 'utils/moveSetupFlow';
import { MTOAgentType } from 'shared/constants';

import 'scenes/Review/Review.css';

export default function HHGShipmentSummary(props) {
  const { mtoShipment, newDutyStationPostalCode } = props;

  const requestedPickupDate = get(mtoShipment, 'requestedPickupDate', '');
  const pickupLocation = get(mtoShipment, 'pickupAddress', {});
  const agents = get(mtoShipment, 'agents', {});
  const releasingAgent = Object.values(agents).find((agent) => agent.agentType === MTOAgentType.RELEASING);
  const releasingAgentFullName = getFullAgentName(releasingAgent);
  const requestedDeliveryDate = get(mtoShipment, 'requestedDeliveryDate', '');
  const dropoffLocation = get(mtoShipment, 'destinationAddress', {});
  const receivingAgent = Object.values(agents).find((agent) => agent.agentType === MTOAgentType.RECEIVING);
  const receivingAgentFullName = getFullAgentName(receivingAgent);
  const remarks = get(mtoShipment, 'customerRemarks', '');

  const hhgPickupLocation = <Address address={pickupLocation} />;

  const destination = isEmpty(dropoffLocation) ? newDutyStationPostalCode : <Address address={dropoffLocation} />;

  // make ReviewSection component a storybook component

  // CSS
  // heading 'edit' should be right aligned
  // line under each field needs to go all the way across
  // field labels should be bold
  // Move setup heading should match Orders and Profile
  // box around each shipment
  // blue line at the top of the shipment box
  // reference CSS module styles
  // style shipment title using an h4

  const hhgShipmentData = [
    { label: 'Requested pickup date', value: formatDateSM(requestedPickupDate) },
    { label: 'Pickup location', value: hhgPickupLocation },
    { label: 'Releasing agent', value: isEmpty(releasingAgent) ? '–' : releasingAgentFullName },
    { label: 'Requested delivery date', value: formatDateSM(requestedDeliveryDate) },
    { label: 'Destination', value: destination },
    { label: 'Receiving agent', value: isEmpty(receivingAgent) ? '–' : receivingAgentFullName },
    { label: 'Remarks', value: !remarks ? '–' : remarks },
  ];

  // update title number when we can support multiple shipments
  // add shipment locator as shipment subheading when it exists
  return (
    <div data-testid="hhg-summary" className="review-content">
      <ReviewSection fieldData={hhgShipmentData} title="Shipment 1: HHG" editLink="" useH4 />
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
  newDutyStationPostalCode: PropTypes.string.isRequired,
};

HHGShipmentSummary.defaultProps = {
  mtoShipment: {},
};
