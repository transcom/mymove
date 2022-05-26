import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updateMTOShipmentAgent';

describe('when given an mto shipment agents update with mto agents table history record', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.updateMTOShipment,
    tableName: 'mto_agents',
    detailsType: d.LABELED,
    changedValues: {
      email: 'grace@email.com',
      first_name: 'Grace',
      phone: '555-555-5555',
    },
    oldValues: {
      agent_type: 'RELEASING_AGENT',
      email: 'gracie@email.com',
      first_name: 'Gracie',
      last_name: 'Griffin',
      phone: '555-555-5551',
    },
    context: [{ shipment_type: 'HHG' }],
  };

  it('correctly matches the Update mto shipment agent event for releasing agents', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    // expect to have formatted the agent correctly
    expect(
      result.getDetailsLabeledDetails({
        changedValues: item.changedValues,
        oldValues: item.oldValues,
        context: item.context,
      }),
    ).toEqual({
      releasing_agent: 'Grace Griffin, 555-555-5555, grace@email.com',
      email: 'grace@email.com',
      first_name: 'Grace',
      phone: '555-555-5555',
      shipment_type: 'HHG',
    });
  });

  it('correctly matches the Update mto shipment agent event for receiving agents', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    // expect to have formatted the agent correctly
    expect(
      result.getDetailsLabeledDetails({
        changedValues: item.changedValues,
        oldValues: { ...item.oldValues, agent_type: 'RECEIVING_AGENT' },
        context: item.context,
      }),
    ).toEqual({
      receiving_agent: 'Grace Griffin, 555-555-5555, grace@email.com',
      email: 'grace@email.com',
      first_name: 'Grace',
      phone: '555-555-5555',
      shipment_type: 'HHG',
    });
  });
});
