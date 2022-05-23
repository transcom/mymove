import o from 'constants/MoveHistory/UIDisplay/Operations';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updateBillableWeightRemarksAsTIO';

describe('when given an update billable weight remarks as tio history record', () => {
  const item = {
    action: 'UPDATE',
    changedValues: { tio_remarks: 'New max billable weight' },
    eventName: o.updateBillableWeightAsTIO,
    tableName: 'moves',
  };
  it('correctly matches the update billable weight remarks as tio event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(item)).toEqual('Updated move');
  });
});
