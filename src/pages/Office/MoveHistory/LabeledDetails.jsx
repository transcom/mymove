/* eslint-disable camelcase */
import React from 'react';

import styles from './LabeledDetails.module.scss';

import booleanFields from 'constants/MoveHistory/Database/BooleanFields';
import dateFields from 'constants/MoveHistory/Database/DateFields';
import fieldMappings from 'constants/MoveHistory/Database/FieldMappings';
import distanceFields from 'constants/MoveHistory/Database/DistanceFields';
import timeUnitFields from 'constants/MoveHistory/Database/TimeUnitFields';
import weightFields from 'constants/MoveHistory/Database/WeightFields';
import monetaryFields from 'constants/MoveHistory/Database/MonetaryFields';
import { shipmentTypes } from 'constants/shipments';
import { HistoryLogRecordShape } from 'constants/MoveHistory/UIDisplay/HistoryLogShape';
import optionFields from 'constants/MoveHistory/Database/OptionFields';
import statusFields from 'constants/MoveHistory/Database/StatusFields';
import { expenseTypeLabels } from 'constants/ppmExpenseTypes.js';
import {
  formatCents,
  formatCustomerDate,
  formatDistanceUnitMiles,
  formatTimeUnitDays,
  formatWeight,
  formatYesNoMoveHistoryValue,
  toDollarString,
} from 'utils/formatters';

export const withMappings = () => {
  const self = {
    displayMappings: [],
  };

  const getResult = ([mapping, fn] = []) => ({
    mapping,
    fn,
  });

  const defaultField = [
    {},
    ({ value }) =>
      (value in optionFields && optionFields[value]) ||
      (value in statusFields && statusFields[value]) ||
      (value in expenseTypeLabels && expenseTypeLabels[value]) ||
      (`${value}` && value) ||
      '—',
  ];

  self.addNameMappings = (mappings) => {
    self.displayMappings = self.displayMappings.concat(mappings);
    return self;
  };
  self.getMappedDisplayName = (field) =>
    getResult(self.displayMappings.find(([mapping]) => field in mapping) || defaultField);
  return self;
};

export const { displayMappings, getMappedDisplayName } = withMappings().addNameMappings([
  [weightFields, ({ value }) => formatWeight(Number(value))],
  [dateFields, ({ value }) => formatCustomerDate(value)],
  [booleanFields, ({ value }) => formatYesNoMoveHistoryValue(value)],
  [monetaryFields, ({ value }) => toDollarString(formatCents(value))],
  [timeUnitFields, ({ value }) => formatTimeUnitDays(value)],
  [distanceFields, ({ value }) => formatDistanceUnitMiles(value)],
  [statusFields, ({ value }) => statusFields[value]],
  [optionFields, ({ value }) => optionFields[value]],
]);

export const retrieveTextToDisplay = (fieldName, value) => {
  const displayName = fieldMappings[fieldName];

  const { fn: valueFormatFn } = getMappedDisplayName(fieldName);
  const displayValue = valueFormatFn({ value });

  if (fieldName === 'has_received_advance') {
    return {
      displayName,
      displayValue: (!`${value}` && '—') || (value && displayValue) || 'No',
    };
  }

  return {
    displayName,
    displayValue: (!`${value}` && '—') || (value !== null && value !== '' && displayValue) || '—',
  };
};

// testable for code coverage //
export const createLineItemLabel = (
  shipmentType,
  shipmentLocator,
  serviceItemName,
  movingExpenseType,
  belongs_to_self,
) =>
  [
    shipmentType && `${shipmentTypes[shipmentType]} shipment #${shipmentLocator}`,
    serviceItemName,
    movingExpenseType && `${expenseTypeLabels[movingExpenseType]}`,
    belongs_to_self,
  ]
    .filter((e) => e)
    .join(', ');

// testable for code coverage //

// Filter out empty values unless they used to be non-empty
// These values may be non-nullish in oldValues and nullish in changedValues
// Use the existing keys in changed or old values to check against keys listed in fieldMappings
export const filterInLineItemValues = (changedValues, oldValues) =>
  Object.entries({ ...oldValues, ...changedValues }).filter(([theField]) => {
    if (!(fieldMappings[theField]?.length >= 0)) return false;

    const changed = changedValues || {};
    const old = oldValues || {};

    const isInChangedValues = theField in changedValues;

    if (isInChangedValues)
      switch (changed[theField]) {
        case undefined:
        case null:
          break;
        default:
          if (changed[theField] === '' && !(theField in old)) return false;
          return true;
      }

    if (isInChangedValues)
      switch (old[theField]) {
        case undefined:
        case null:
        case '':
          break;
        default:
          return true;
      }

    return false;
  });

const LabeledDetails = ({ historyRecord }) => {
  const { changedValues, oldValues = {} } = historyRecord;

  const { shipment_type, shipment_locator, service_item_name, ...changedValuesToUse } = changedValues;

  let belongs_to_self =
    oldValues?.belongs_to_self !== null && oldValues?.belongs_to_self !== ''
      ? oldValues?.belongs_to_self
      : changedValues.belongs_to_self;
  const moving_expense_type = oldValues?.moving_expense_type
    ? oldValues?.moving_expense_type
    : changedValues.moving_expense_type;

  if (belongs_to_self === true) belongs_to_self = 'Service member pro-gear';
  else if (belongs_to_self === false) belongs_to_self = 'Spouse pro-gear';

  // Check for shipment_type to use it as a header for the row
  const shipmentDisplay = createLineItemLabel(
    shipment_type,
    shipment_locator,
    service_item_name,
    moving_expense_type,
    belongs_to_self,
  );

  const lineItems = filterInLineItemValues(changedValuesToUse, oldValues).map(([label, value]) => {
    const { displayName, displayValue } = retrieveTextToDisplay(label, value);
    return (
      <div key={label}>
        <b>{displayName}</b>: {displayValue}
      </div>
    );
  });

  return (
    <>
      <span className={styles.shipmentType}>{shipmentDisplay}</span>
      {lineItems}
    </>
  );
};

LabeledDetails.propTypes = {
  historyRecord: HistoryLogRecordShape,
};

LabeledDetails.defaultProps = {
  historyRecord: {},
};

export default LabeledDetails;
