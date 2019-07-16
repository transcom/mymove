import { MOVE_DOC_TYPE, MOVE_DOC_STATUS } from 'shared/constants';
import { calcWeightTicketNetWeight, findPendingWeightTickets } from './moveDocuments';

describe('Move Document utility functions', () => {
  it('calcWeightTicketNetWeight should return 0 if there are no weight ticket sets', () => {
    const moveDocs = [
      { status: MOVE_DOC_STATUS.OK, move_document_type: MOVE_DOC_TYPE.GBL },
      { stauts: MOVE_DOC_STATUS.OK, move_document_type: MOVE_DOC_TYPE.EXPENSE },
    ];
    expect(calcWeightTicketNetWeight([])).toBe(0);
    expect(calcWeightTicketNetWeight(moveDocs)).toBe(0);
  });

  it('calcWeightTicketNetWeight should return net weight of all OKed weight ticket sets', () => {
    const moveDocs = [
      {
        status: MOVE_DOC_STATUS.OK,
        move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
        full_weight: 3000,
        empty_weight: 1000,
      },
      {
        status: MOVE_DOC_STATUS.HAS_ISSUE,
        move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
        full_weight: 9000,
        empty_weight: 3000,
      },
      {
        status: MOVE_DOC_STATUS.AWAITING_REVIEW,
        move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
        full_weight: 4000,
        empty_weight: 1000,
      },
      {
        status: MOVE_DOC_STATUS.OK,
        move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
        full_weight: 4000,
        empty_weight: 1000,
      },
    ];

    expect(calcWeightTicketNetWeight(moveDocs)).toBe(5000);
  });

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
