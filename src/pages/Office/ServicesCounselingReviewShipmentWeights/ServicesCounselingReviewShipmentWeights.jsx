import React from 'react';

import { createHeader } from '../../../components/Table/utils';

import { useReviewShipmentWeightsQuery } from 'hooks/queries';
import { calculateTotalNetWeightForWeightTickets } from 'utils/ppmCloseout';

export const PPMReviewWeightsTableColumns = [
  createHeader('', (row) => row?.shipmentType, {
    id: 'shipmentType',
    isFilterable: false,
  }),
  // TODO get url for row
  createHeader('Weight ticket', (row) => <a href={row?.url}> Review Documents </a>, {
    id: 'weightTicket',
    isFilterable: false,
  }),
  createHeader('Pro-gear (lbs)', (row) => (row.ppmShipment.proGearWeight > 0 ? row.ppmShipment.proGearWeight : '-'), {
    id: 'proGear',
    isFilterable: false,
  }),
  createHeader(
    'Spouse pro-gear',
    (row) => (row.ppmShipment.spouseProGearWeight > 0 ? row.ppmShipment.spouseProGearWeight : '-'),
    {
      id: 'spouseProGear',
      isFilterable: false,
    },
  ),
  createHeader(
    'Estimated Weight',
    (row) => (row.ppmShipment.estimatedWeight > 0 ? row.ppmShipment.estimatedWeight : '-'),
    {
      id: 'estimatedWeight',
      isFilterable: false,
    },
  ),
  createHeader(
    'Net Weight',
    (row) => {
      const calculatedNetWeight = calculateTotalNetWeightForWeightTickets(row.ppmShipment?.weightTickets);
      return calculatedNetWeight > 0 ? calculatedNetWeight : '-';
    },
    {
      id: 'netWeight',
      isFilterable: false,
    },
  ),
];

export const ProGearTableColumns = [];

export const NonPPMTableColumns = [];

const ServicesCounselingReviewShipmentWeights = ({ moveCode }) => {
  useReviewShipmentWeightsQuery(moveCode);
  return <h1>Review shipment weights</h1>;
};

export default ServicesCounselingReviewShipmentWeights;
