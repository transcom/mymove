import React from 'react';
import PropTypes, { string, number } from 'prop-types';
import { get, isEmpty } from 'lodash';

import ReviewSection from '../../components/Customer/ReviewSection';
import SectionWrapper from '../../components/Customer/SectionWrapper';

import Address from 'scenes/Review/Address';
import { formatDateSM } from 'shared/formatters';
import { getFullAgentName } from 'utils/moveSetupFlow';
import { MTOAgentType } from 'shared/constants';

import 'scenes/Review/Review.css';

export default function HHGShipmentSummary(props) {
  const { mtoShipment, movePath, newDutyStationPostalCode, shipmentNumber } = props;
  const editShipmentPath = `${movePath}/mto-shipments/${mtoShipment?.id}/edit-shipment?shipmentNumber=${shipmentNumber}`;

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

  const hhgShipmentData = [
    { label: 'Requested pickup date', value: formatDateSM(requestedPickupDate) },
    { label: 'Pickup location', value: hhgPickupLocation },
    { label: 'Releasing agent', value: isEmpty(releasingAgent) ? '–' : releasingAgentFullName },
    { label: 'Requested delivery date', value: formatDateSM(requestedDeliveryDate) },
    { label: 'Destination', value: destination },
    { label: 'Receiving agent', value: isEmpty(receivingAgent) ? '–' : receivingAgentFullName },
    { label: 'Remarks', value: !remarks ? '–' : remarks },
  ];

  // add shipment locator as shipment subheading when it exists
  return (
    <div data-testid="hhg-summary" className="review-content">
      <SectionWrapper>
        <ReviewSection
          fieldData={hhgShipmentData}
          title={`Shipment ${shipmentNumber}: HHG`}
          editLink={editShipmentPath}
          useH4
        />
      </SectionWrapper>
    </div>
  );
}

HHGShipmentSummary.propTypes = {
  movePath: string.isRequired,
  mtoShipment: PropTypes.shape({
    id: string.isRequired,
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
  shipmentNumber: number.isRequired,
};

HHGShipmentSummary.defaultProps = {
  mtoShipment: {},
};
