import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.counselingUpdateAllowance,
  tableName: t.service_members,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated service member',
};
