import e from 'constants/MoveHistory/EventTemplates/UpdateMTOStatusServiceCounselingCompleted/updateMTOStatusServiceCounselingCompleted';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';

describe('When given a completed services counseling for a move', () => {
  const item = {
    action: a.UPDATE,
    eventName: o.updateMTOStatusServiceCounselingCompleted,
    tableName: t.moves,
  };
  it('correctly matches the update mto status services counseling completed event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getDetailsPlainText(item)).toEqual('Counseling Completed');
  });
});
