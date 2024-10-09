import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/updateMTOShipment';
import * as fieldDefault from 'constants/MoveHistory/Database/FieldMappings';
import {
  formatCents,
  formatCustomerDate,
  formatWeight,
  formatYesNoMoveHistoryValue,
  toDollarString,
} from 'utils/formatters';
import { retrieveTextToDisplay } from 'pages/Office/MoveHistory/LabeledDetails';

const FieldMappings = fieldDefault.default;

describe('when given an mto shipment update with mto shipment table history record', () => {
  const changedValues = {
    destination_address_type: 'HOME_OF_SELECTION',
    requested_delivery_date: '2020-04-14',
    requested_pickup_date: '2020-03-23',
    actual_pickup_date: '2020-04-15',
    approved_date: '2020-04-21',
    billable_weight_cap: '10000',
    billable_weight_justification: 'Heavy items',
    counselor_remarks: 'Remarks',
    customer_remarks: 'Words',
    diversion: false,
    first_available_delivery_date: '2020-04-12',
    id: 'b1cb1428-6d65-40fc-addd-b57f4ea120f1',
    prime_actual_weight: '14000',
    prime_estimated_weight: '12000',
    sac_type: 'HHG',
    scheduled_delivery_date: '2020-04-13',
    scheduled_pickup_date: '2020-03-22',
    service_order_number: '767567576',
    status: 'SUBMITTED',
    tac_type: 'NTS',
    uses_external_vendor: true,

    advance_amount_requested: 100,
    destination_postal_code: '29102',
    estimated_incentive: 2252814,
    estimated_weight: 600,
    expected_departure_date: '2024-02-18',
    has_requested_advance: true,
    pro_gear_weight: 243,
    sit_estimated_cost: 11627,
    sit_estimated_departure_date: '2024-02-29',
    sit_estimated_entry_date: '2024-02-18',
    sit_estimated_weight: 524,
    spouse_pro_gear_weight: 257,
    distance: 400,
  };
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateMTOShipment',
    tableName: 'mto_shipments',
    changedValues,
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: 'b4b4b',
      },
    ],
  };
  it('correctly matches the Update mto shipment event', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
  });
  describe('it correctly renders the details component for Create MTO shipments', () => {
    it.each([
      [FieldMappings.status, changedValues.status],
      [FieldMappings.requested_delivery_date, formatCustomerDate(changedValues.requested_delivery_date)],
      [FieldMappings.requested_pickup_date, formatCustomerDate(changedValues.requested_pickup_date)],
      Object.values(retrieveTextToDisplay('destination_address_type', changedValues.destination_address_type)),
      [FieldMappings.actual_pickup_date, formatCustomerDate(changedValues.actual_pickup_date)],
      [FieldMappings.approved_date, formatCustomerDate(changedValues.approved_date)],
      [FieldMappings.billable_weight_cap, formatWeight(Number(changedValues.billable_weight_cap))],
      [FieldMappings.billable_weight_justification, changedValues.billable_weight_justification],
      [FieldMappings.counselor_remarks, changedValues.counselor_remarks],
      [FieldMappings.customer_remarks, changedValues.customer_remarks],
      [FieldMappings.diversion, formatYesNoMoveHistoryValue(changedValues.diversion)],
      [FieldMappings.first_available_delivery_date, formatCustomerDate(changedValues.first_available_delivery_date)],
      [FieldMappings.shipment_weight, formatWeight(changedValues.shipment_weight)],
      [FieldMappings.prime_estimated_weight, formatWeight(Number(changedValues.prime_estimated_weight))],
      [FieldMappings.sac_type, changedValues.sac_type],
      [FieldMappings.tac_type, changedValues.tac_type],
      [FieldMappings.scheduled_pickup_date, formatCustomerDate(changedValues.scheduled_pickup_date)],
      [FieldMappings.service_order_number, `${changedValues.service_order_number}`],
      [FieldMappings.uses_external_vendor, formatYesNoMoveHistoryValue(changedValues.uses_external_vendor)],
      Object.values(retrieveTextToDisplay('distance', changedValues.distance)),
      [FieldMappings.advance_amount_requested, toDollarString(changedValues.has_requested_advance)],
      [FieldMappings.destination_postal_code, changedValues.destination_postal_code],
      [FieldMappings.estimated_incentive, toDollarString(formatCents(changedValues.estimated_incentive))],
      [FieldMappings.estimated_weight, formatWeight(Number(changedValues.estimated_weight))],
      [FieldMappings.expected_departure_date, formatCustomerDate(changedValues.expected_departure_date)],
      [FieldMappings.has_requested_advance, formatYesNoMoveHistoryValue(changedValues.has_requested_advance)],
      [FieldMappings.pro_gear_weight, formatWeight(Number(changedValues.pro_gear_weight))],
      [FieldMappings.sit_estimated_cost, toDollarString(formatCents(changedValues.sit_estimated_cost))],
      [FieldMappings.sit_estimated_departure_date, formatCustomerDate(changedValues.sit_estimated_departure_date)],
      [FieldMappings.sit_estimated_entry_date, formatCustomerDate(changedValues.sit_estimated_entry_date)],
      [FieldMappings.sit_estimated_weight, formatWeight(Number(changedValues.sit_estimated_weight))],
      [FieldMappings.spouse_pro_gear_weight, formatWeight(Number(changedValues.spouse_pro_gear_weight))],
    ])('displays the correct details value for %s', async (label, value) => {
      const targetItem = Object.fromEntries(
        Object.entries(changedValues).filter(([key]) => FieldMappings[key] === label),
      );
      const history = { ...historyRecord, changedValues: { ...targetItem } };
      const result = getTemplate(history);
      render(result.getDetails(history));
      const displayingElements = screen.getAllByText(label);
      const displayingElement = displayingElements.find((element) => element.parentElement.textContent.includes(label));
      const parent = displayingElement.parentElement;
      expect(parent.textContent).toContain(label);
      expect(parent.textContent).toContain(value);
    });
    it('displays the correct label for shipment', () => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText('PPM shipment #RQ38D4-01')).toBeInTheDocument();
    });
  });
});
