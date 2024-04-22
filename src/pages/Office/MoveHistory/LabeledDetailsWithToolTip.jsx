import React from 'react';

import styles from './LabeledDetails.module.scss';

import dateFields from 'constants/MoveHistory/Database/DateFields';
import fieldMappings from 'constants/MoveHistory/Database/FieldMappings';
import weightFields from 'constants/MoveHistory/Database/WeightFields';
import { shipmentTypes } from 'constants/shipments';
import { HistoryLogRecordShape } from 'constants/MoveHistory/UIDisplay/HistoryLogShape';
import optionFields from 'constants/MoveHistory/Database/OptionFields';
import { formatCustomerDate, formatWeight, formatYesNoMoveHistoryValue } from 'utils/formatters';
import ToolTip from 'shared/ToolTip/ToolTip';

const retrieveTextToDisplay = (fieldName, value) => {
  const emptyValue = '—';
  const displayName = fieldMappings[fieldName];
  let displayValue = value;

  if (displayName === fieldMappings.storage_in_transit) {
    displayValue = `${displayValue} days`;
  } else if (weightFields[fieldName]) {
    // turn string value into number so it can be formatted correctly
    displayValue = formatWeight(Number(displayValue));
  } else if (optionFields[displayValue]) {
    displayValue = optionFields[displayValue];
  } else if (dateFields[fieldName]) {
    displayValue = formatCustomerDate(displayValue);
  } else if (displayName === fieldMappings.dependents_authorized || displayName === fieldMappings.has_dependents) {
    displayValue = formatYesNoMoveHistoryValue(displayValue);
  }

  if (!displayValue) {
    displayValue = emptyValue;
  }

  return {
    displayName,
    displayValue,
  };
};

const LabeledDetailsWithToolTip = ({ historyRecord, toolTipText, toolTipColor, toolTipTextPosition, toolTipIcon }) => {
  const changedValuesToUse = historyRecord.changedValues;
  let shipmentDisplay = '';

  // Check for shipment_type to use it as a header for the row
  if ('shipment_type' in changedValuesToUse) {
    shipmentDisplay = shipmentTypes[changedValuesToUse.shipment_type];
    shipmentDisplay += ` shipment #${changedValuesToUse.shipment_locator}`;
    delete changedValuesToUse.shipment_type;
  }

  // check for service item and add it to the header
  if ('service_item_name' in changedValuesToUse) {
    shipmentDisplay += `, ${changedValuesToUse.service_item_name}`;
    delete changedValuesToUse.service_item_name;
  }

  /* Filter out empty values unless they used to be non-empty
     These values may be non-nullish in oldValues and nullish in changedValues */
  const dbFieldsToDisplay = Object.keys(fieldMappings).filter((dbField) => {
    return (
      changedValuesToUse[dbField] ||
      (dbField in changedValuesToUse && historyRecord.oldValues && historyRecord.oldValues[dbField])
    );
  });

  return (
    <>
      <span className={styles.shipmentType}>{shipmentDisplay}</span>
      {dbFieldsToDisplay.map((modelField) => {
        const { displayName, displayValue } = retrieveTextToDisplay(modelField, changedValuesToUse[modelField]);

        return (
          <div key={modelField}>
            <b>{displayName}</b>: {displayValue}{' '}
            <ToolTip text={toolTipText} color={toolTipColor} position={toolTipTextPosition} icon={toolTipIcon} />
          </div>
        );
      })}
    </>
  );
};

LabeledDetailsWithToolTip.propTypes = {
  historyRecord: HistoryLogRecordShape,
};

LabeledDetailsWithToolTip.defaultProps = {
  historyRecord: {},
};

export default LabeledDetailsWithToolTip;
