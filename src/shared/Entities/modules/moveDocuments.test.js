import { MOVE_DOC_TYPE, MOVE_DOC_STATUS } from 'shared/constants';
import { findPendingWeightTickets } from './moveDocuments';

describe('Move Document utility functions', () => {
  it('findingPendingWeightTickets finds weight ticket sets that are not OK status', () => {
    const moveDocs = [
      { status: MOVE_DOC_STATUS.OK, move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET },
      { status: MOVE_DOC_STATUS.AWAITING_REVIEW, move_document_type: MOVE_DOC_TYPE.GBL },
      { status: MOVE_DOC_STATUS.HAS_ISSUE, move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET },
      { status: MOVE_DOC_STATUS.AWAITING_REVIEW, move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET },
    ];
    expect(findPendingWeightTickets(moveDocs).length).toBe(2);
  });

  it('findPendingWeightTickets does not find any weight ticket sets that are not OK status', () => {
    const moveDocs = [
      { status: MOVE_DOC_STATUS.OK, move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET },
      { status: MOVE_DOC_STATUS.AWAITING_REVIEW, move_document_type: MOVE_DOC_TYPE.GBL },
      { status: MOVE_DOC_STATUS.HAS_ISSUE, move_document_type: MOVE_DOC_TYPE.EXPENSE },
    ];
    expect(findPendingWeightTickets(moveDocs).length).toBe(0);
  });
});
