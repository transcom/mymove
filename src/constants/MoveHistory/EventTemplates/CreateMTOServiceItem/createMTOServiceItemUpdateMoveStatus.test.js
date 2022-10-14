import getTemplate from 'constants/MoveHistory/TemplateManager';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import e from 'constants/MoveHistory/EventTemplates/CreateMTOServiceItem/createMTOServiceItemUpdateMoveStatus';

describe('when given a move status update with create mto service item history record', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.createMTOServiceItem,
    tableName: 'moves',
    detailsType: d.LABELED,
    changedValues: {
      status: 'STATUS CHANGED',
    },
  };
  it('correctly matches the create MTO service item event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
  });
});
