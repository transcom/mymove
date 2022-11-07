import getTemplate from 'constants/MoveHistory/TemplateManager';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/createMTOShipmentAgent';

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
    context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
  };

  const item2 = {
    action: 'INSERT',
    eventName: o.createMTOShipment,
    tableName: 'mto_agents',
    detailsType: d.LABELED,
    changedValues: {
      email: 'catalina@email.com',
      first_name: 'Catalina',
      last_name: 'Washington',
      phone: '999-999-9999',
      agent_type: 'RELEASING_AGENT',
    },
    context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
  };

  it('correctly matches the insert mto shipment agent event for releasing agents when shipment updated', () => {
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
      shipment_id_display: 'A1B2C',
    });
  });

  it('correctly matches the insert mto shipment agent event for releasing agents when shipment created', () => {
    const result = getTemplate(item2);
    expect(result).toMatchObject(e);
    // expect to have formatted the agent correctly
    expect(
      result.getDetailsLabeledDetails({
        changedValues: item2.changedValues,
        context: item2.context,
      }),
    ).toEqual({
      releasing_agent: 'Catalina Washington, 999-999-9999, catalina@email.com',
      email: 'catalina@email.com',
      first_name: 'Catalina',
      last_name: 'Washington',
      phone: '999-999-9999',
      agent_type: 'RELEASING_AGENT',
      shipment_type: 'HHG',
      shipment_id_display: 'A1B2C',
    });
  });

  it('correctly matches the insert mto shipment agent event for receiving agents when shipment updated', () => {
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
      shipment_id_display: 'A1B2C',
    });
  });

  it('correctly matches the insert mto shipment agent event for receiving agents when shipment created', () => {
    const result = getTemplate(item2);
    expect(result).toMatchObject(e);
    // expect to have formatted the agent correctly
    expect(
      result.getDetailsLabeledDetails({
        changedValues: { ...item2.changedValues, agent_type: 'RECEIVING_AGENT' },
        context: item2.context,
      }),
    ).toEqual({
      receiving_agent: 'Catalina Washington, 999-999-9999, catalina@email.com',
      email: 'catalina@email.com',
      first_name: 'Catalina',
      last_name: 'Washington',
      phone: '999-999-9999',
      agent_type: 'RECEIVING_AGENT',
      shipment_type: 'HHG',
      shipment_id_display: 'A1B2C',
    });
  });
});
