import React from 'react';
import PropTypes from 'prop-types';

import PlainTextDetails from './PlainTextDetails';

import {
  eventNamesWithLabelledDetails,
  eventNamesWithServiceItemDetails,
  eventNamesWithPlainTextDetails,
  HistoryLogValuesShape,
} from 'constants/historyLogUIDisplayName';

const formatChangedValues = (values) => {
  return values
    ? values.map((value) => (
        <div key={`${value.columnName}-${value.columnValue}`}>
          {value.columnName}: {value.columnValue}
        </div>
      ))
    : '';
};

const MoveHistoryDetailsSelector = ({ eventName, oldValues, changedValues }) => {
  /**
   * Inside the component, we should map oldValues and changedValues into an object so the ordering can be consistent.
   */
  if (eventNamesWithLabelledDetails[eventName]) {
    return (
      <div>
        Labelled {eventName}
        <div>old Values {formatChangedValues(oldValues)}</div>
        <div>changed values {formatChangedValues(changedValues)}</div>
      </div>
    );
  }

  if (eventNamesWithServiceItemDetails[eventName]) {
    return (
      <div>
        Service Items {eventName}
        <div>old Values {formatChangedValues(oldValues)}</div>
        <div>changed values {formatChangedValues(changedValues)}</div>
      </div>
    );
  }

  if (eventNamesWithPlainTextDetails[eventName]) {
    return <PlainTextDetails eventName={eventName} changedValues={changedValues} />;
  }

  return (
    <div>
      - {eventName}
      <div>old Values {formatChangedValues(oldValues)}</div>
      <div>changed values {formatChangedValues(changedValues)}</div>
    </div>
  );
};

MoveHistoryDetailsSelector.propTypes = {
  eventName: PropTypes.string,
  oldValues: HistoryLogValuesShape,
  changedValues: HistoryLogValuesShape,
};

MoveHistoryDetailsSelector.defaultProps = {
  eventName: '',
  oldValues: [],
  changedValues: [],
};

export default MoveHistoryDetailsSelector;
