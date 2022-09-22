import getTemplate from 'constants/MoveHistory/TemplateManager';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import updateAllowanceCounseling from 'constants/MoveHistory/EventTemplates/updateAllowanceCounseling';

describe('When a service counselor updates shipping allowances', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.counselingUpdateAllowance,
    tableName: t.entitlements,
    detailsType: d.LABELED,
    eventNameDisplay: 'Updated allowances',
    changedValues: {
      dependents_authorized: 'false',
      organizational_clothing_and_individual_equipment: 'false',
      pro_gear_weight: '1999',
      pro_gear_weight_spouse: '49g',
      required_medical_equipment_weight: '99g',
      storage_in_transit: 'gg',
    },
  };
  it('correctly matches the update allowances event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(updateAllowanceCounseling);
    expect(result.getEventNameDisplay()).toMatch(item.eventNameDisplay);
    expect(result.getDetailsLabeledDetails(item)).toMatchObject({
      dependents_authorized: 'false',
      organizational_clothing_and_individual_equipment: 'false',
      pro_gear_weight: '1999',
      pro_gear_weight_spouse: '49g',
      required_medical_equipment_weight: '99g',
      storage_in_transit: 'gg',
    });
  });
});
