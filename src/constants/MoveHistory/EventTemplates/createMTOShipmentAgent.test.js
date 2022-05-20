import getTemplate from 'constants/MoveHistory/TemplateManager';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import e from 'constants/MoveHistory/EventTemplates/createMTOShipmentAgent';

describe('when given an mto shipment agents insert with mto agents table history record', () => {
  const item = {
    action: 'INSERT',
    eventName: o.updateMTOShipment,
    tableName: 'mto_agents',
    detailsType: d.LABELED,
    changedValues: {
      email: 'grace@email.com',
      first_name: 'Grace',
      last_name: 'Griffin',
      phone: '555-555-5555',
      agent_type: 'RELEASING_AGENT',
    },
    context: [{ shipment_type: 'HHG' }],
  };

  it('correctly matches the insert mto shipment agent event for releasing agents', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    // expect to have formatted the agent correctly
    expect(
      result.getDetailsLabeledDetails({
        changedValues: item.changedValues,
        context: item.context,
      }),
    ).toEqual({
      releasing_agent: 'Grace Griffin, 555-555-5555, grace@email.com',
      email: 'grace@email.com',
      first_name: 'Grace',
      last_name: 'Griffin',
      phone: '555-555-5555',
      agent_type: 'RELEASING_AGENT',
      shipment_type: 'HHG',
    });
  });

  it('correctly matches the insert mto shipment agent event for receiving agents', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    // expect to have formatted the agent correctly
    expect(
      result.getDetailsLabeledDetails({
        changedValues: { ...item.changedValues, agent_type: 'RECEIVING_AGENT' },
        context: item.context,
      }),
    ).toEqual({
      receiving_agent: 'Grace Griffin, 555-555-5555, grace@email.com',
      email: 'grace@email.com',
      first_name: 'Grace',
      last_name: 'Griffin',
      phone: '555-555-5555',
      agent_type: 'RECEIVING_AGENT',
      shipment_type: 'HHG',
    });
  });
});
