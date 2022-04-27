import React from 'react';

import {
  HistoryLogValuesShape,
  dbFieldToDisplayName,
  dbWeightFields,
  HistoryLogContextShape,
} from 'constants/historyLogUIDisplayName';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatMoveHistoryFullAddress, formatMoveHistoryAgent } from 'utils/formatters';

const LabeledDetails = ({ changedValues, oldValues, context }) => {
  const backfilledChangedValues = {
    street_address_1: oldValues.street_address_1,
    street_address_2: oldValues.street_address_2,
    city: oldValues.city,
    state: oldValues.state,
    postal_code: oldValues.postal_code,
    email: oldValues.email,
    first_name: oldValues.first_name,
    last_name: oldValues.last_name,
    phone: oldValues.phone,
    ...changedValues,
  };

  const changedValuesWithFormattedItems = {
    ...changedValues,
    address: formatMoveHistoryFullAddress(backfilledChangedValues),
    agent: formatMoveHistoryAgent(backfilledChangedValues),
  };

  const dbFieldsToDisplay = Object.keys(dbFieldToDisplayName).filter((dbField) => {
    return changedValuesWithFormattedItems[dbField];
  });

  const retrieveTextToDisplay = (fieldName, value) => {
    let displayName = dbFieldToDisplayName[fieldName];
    let displayValue = value;

    if (displayName === dbFieldToDisplayName.storage_in_transit) {
      displayValue = `${displayValue} days`;
    } else if (dbWeightFields.includes(fieldName)) {
      displayValue = `${displayValue} lbs`;
    } else if (displayName === dbFieldToDisplayName.address) {
      const { addressType } = context.filter((contextObject) => contextObject.addressType)[0];

      if (addressType === 'pickupAddress') {
        displayName = 'Origin address';
      } else if (addressType === 'destinationAddress') {
        displayName = 'Destination address';
      }
    } else if (displayName === dbFieldToDisplayName.agent) {
      const agentType = changedValues.agent_type ?? oldValues.agent_type;

      if (agentType === 'RECEIVING_AGENT') {
        displayName = 'Receiving agent';
      } else if (agentType === 'RELEASING_AGENT') {
        displayName = 'Releasing agent';
      }
    }

    return {
      displayName,
      displayValue,
    };
  };

  return (
    <div>
      {dbFieldsToDisplay.map((modelField) => {
        const { displayName, displayValue } = retrieveTextToDisplay(
          modelField,
          changedValuesWithFormattedItems[modelField],
        );

        return (
          <div key={modelField} className={descriptionListStyles.row}>
            <b>{displayName}</b>: {displayValue}
          </div>
        );
      })}
    </div>
  );
};

LabeledDetails.propTypes = {
  changedValues: HistoryLogValuesShape,
  oldValues: HistoryLogValuesShape,
  context: HistoryLogContextShape,
};

LabeledDetails.defaultProps = {
  changedValues: {},
  oldValues: {},
  context: [],
};

export default LabeledDetails;
