import React from 'react';

import styles from './UpdateMTOShipmentAgent.module.scss';

import { formatMoveHistoryAgent } from 'utils/formatters';
import fieldMappings from 'constants/MoveHistory/Database/FieldMappings';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { oldValues, changedValues } = historyRecord;
  let newChangedValues = {
    email: oldValues.email,
    first_name: oldValues.first_name,
    last_name: oldValues.last_name,
    phone: oldValues.phone,
    agent_type: oldValues.agent_type,
    ...changedValues,
  };

  newChangedValues = {
    ...formatMoveHistoryAgent(newChangedValues),
    ...getMtoShipmentLabel(historyRecord),
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

const formatDeletedAgentRecord = (historyRecord) => {
  const mtoShipmentLabel = getMtoShipmentLabel(historyRecord);
  const historyLabel = `${mtoShipmentLabel.shipment_type} shipment #${mtoShipmentLabel.shipment_locator}`;
  const agentTypeFieldName = historyRecord.oldValues.agent_type.toLowerCase();
  const agentType = fieldMappings[agentTypeFieldName];
  const agent = formatMoveHistoryAgent(historyRecord.oldValues)[agentTypeFieldName];

  return (
    <>
      <span className={styles.shipmentType}>
        Deleted {agentType} on {historyLabel}
      </span>
      <div>
        <b>{agentType}</b>: {agent}
      </div>
    </>
  );
};

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.mto_agents,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: (historyRecord) => {
    if (historyRecord.changedValues.deleted_at) {
      return formatDeletedAgentRecord(historyRecord);
    }
    return <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />;
  },
};
