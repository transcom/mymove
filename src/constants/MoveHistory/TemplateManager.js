import undefinedEvent from 'constants/MoveHistory/EventTemplates/NullEvent/undefined';
import * as eventTemplates from 'constants/MoveHistory/EventTemplates';

const allMoveHistoryEventTemplates = [];

const registerTemplate = ({
  action,
  eventName,
  tableName,
  getEventNameDisplay = () => {
    return 'Undefined event type';
  },
  getDetails = null,
}) => {
  const eventType = {};
  eventType.action = action;
  eventType.eventName = eventName;
  eventType.tableName = tableName;
  eventType.getEventNameDisplay = getEventNameDisplay;
  eventType.getDetails = getDetails;

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
};

// Iterate on all the Event Templates pulled into eventTemplates.
Object.values(eventTemplates).forEach((o) => registerTemplate(o));

export default (historyRecord) => {
  return allMoveHistoryEventTemplates.find((eventType) => eventType.matches(historyRecord)) || undefinedEvent;
};
