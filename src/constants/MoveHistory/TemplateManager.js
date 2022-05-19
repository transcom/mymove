import detailsTypes from 'constants/MoveHistory/DetailsColumn/Types';
import undefinedEvent from 'constants/MoveHistory/EventTemplates/undefinedEvent';

const allMoveHistoryEventTemplates = [];

export const buildTemplate = ({
  action = '*',
  eventName = '*',
  tableName = '*',
  detailsType = detailsTypes.PLAIN_TEXT,
  getEventNameDisplay = () => {
    return 'Undefined event type';
  },
  getDetailsPlainText = () => {
    return 'Undefined details';
  },
  getStatusDetails = () => {
    return 'Undefined status';
  },
  getDetailsLabeledDetails = null,
}) => {
  const eventType = {};
  eventType.action = action;
  eventType.eventName = eventName;
  eventType.tableName = tableName;
  eventType.detailsType = detailsType;
  eventType.getEventNameDisplay = getEventNameDisplay;
  eventType.getDetailsPlainText = getDetailsPlainText;
  eventType.getStatusDetails = getStatusDetails;
  eventType.getDetailsLabeledDetails = getDetailsLabeledDetails;

  // Used for matching properties on Events when building an Event Template
  function propertiesMatch(p1, p2) {
    return p1 === '*' || p2 === '*' || p1 === p2;
  }

  eventType.matches = (other) => {
    if (eventType === undefined || other === undefined) {
      return false;
    }
    return (
      propertiesMatch(eventType.action, other?.action) &&
      propertiesMatch(eventType.eventName, other?.eventName) &&
      propertiesMatch(eventType.tableName, other?.tableName)
    );
  };

  allMoveHistoryEventTemplates.push(eventType);

  return eventType;
};

export const getTemplate = (historyRecord) => {
  const foundTemplate = allMoveHistoryEventTemplates.find((eventType) => eventType.matches(historyRecord));
  return false || undefinedEvent();
};
