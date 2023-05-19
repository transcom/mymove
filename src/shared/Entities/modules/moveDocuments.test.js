import { WEIGHT_TICKET_SET_TYPE, MOVE_DOC_TYPE, MOVE_DOC_STATUS } from '../../constants';

import { findPendingWeightTickets, findOKedVehicleWeightTickets, findOKedProgearWeightTickets } from './moveDocuments';

describe('Move Document utility functions', () => {
  const assortedMoveDocuments = [
    {
      status: MOVE_DOC_STATUS.AWAITING_REVIEW,
      move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
      weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.CAR,
    },
    {
      status: MOVE_DOC_STATUS.AWAITING_REVIEW,
      move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
      weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.PRO_GEAR,
    },
    {
      status: MOVE_DOC_STATUS.EXCLUDE,
      move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
      weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.PRO_GEAR,
    },
    { status: MOVE_DOC_STATUS.AWAITING_REVIEW, move_document_type: MOVE_DOC_TYPE.GBL },
    {
      status: MOVE_DOC_STATUS.OK,
      move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
      weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.CAR,
    },
    {
      status: MOVE_DOC_STATUS.OK,
      move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
      weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.PRO_GEAR,
    },
  ];

  it('findOKedProgearWeightTickets finds pending pro-gear weight ticket sets', () => {
    expect(findOKedProgearWeightTickets(assortedMoveDocuments)).toEqual([
      {
        status: MOVE_DOC_STATUS.OK,
        move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
        weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.PRO_GEAR,
      },
    ]);
  });

  it('findOKedVehicleWeightTickets finds pending vehicle weight ticket sets', () => {
    expect(findOKedVehicleWeightTickets(assortedMoveDocuments)).toEqual([
      {
        status: MOVE_DOC_STATUS.OK,
        move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
        weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.CAR,
      },
    ]);
  });

  it('findPendingWeightTickets finds pending weight ticket sets', () => {
    expect(findPendingWeightTickets(assortedMoveDocuments)).toEqual([
      {
        status: MOVE_DOC_STATUS.AWAITING_REVIEW,
        move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
        weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.CAR,
      },
      {
        status: MOVE_DOC_STATUS.AWAITING_REVIEW,
        move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
        weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.PRO_GEAR,
      },
    ]);
  });
});
