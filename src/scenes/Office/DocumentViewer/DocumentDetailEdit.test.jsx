import { shallow } from 'enzyme';
import DocumentDetailEdit from './DocumentDetailEdit';
import React from 'react';
import { MOVE_DOC_TYPE, WEIGHT_TICKET_SET_TYPE } from '../../../shared/constants';

describe('DocumentDetailEdit', () => {
  const renderDocumentDetailEdit = ({ moveDocSchema = {}, formValues = { moveDocument: {} } }) =>
    shallow(<DocumentDetailEdit formValues={formValues} moveDocSchema={moveDocSchema} />);

  const moveDocSchema = {
    properties: {},
    required: [],
    type: 'string type',
  };

  describe('weight ticket document edit', () => {
    it('shows all form fields for a car', () => {
      const formValues = {
        moveDocument: {
          move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
          weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.CAR,
        },
      };

      const documentForm = renderDocumentDetailEdit({ formValues, moveDocSchema });
      const title = documentForm.find('[data-cy="document-title-field"]');
      const moveDocumentType = documentForm.find('[data-cy="move-document-type"]');
      const weightTicketSetType = documentForm.find('[data-cy="weight-ticket-set-type"]');
      const make = documentForm.find('[data-cy="vehicle-make"]');
      const model = documentForm.find('[data-cy="vehicle-model"]');
      const vehicleNickname = documentForm.find('[data-cy="vehicle-nickname"]');
      const status = documentForm.find('[data-cy="status"]');
      const notes = documentForm.find('[data-cy="notes"]');

      expect(title.length).toEqual(1);
      expect(moveDocumentType.length).toEqual(1);
      expect(weightTicketSetType.length).toEqual(1);
      expect(weightTicketSetType.length).toEqual(1);
      expect(make.length).toEqual(1);
      expect(model.length).toEqual(1);
      expect(vehicleNickname.length).toEqual(0);
      expect(status.length).toEqual(1);
      expect(notes.length).toEqual(1);
    });

    it('shows all form fields for a boxtruck', () => {
      const formValues = {
        moveDocument: {
          move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
          weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.BOX_TRUCK,
        },
      };

      const documentForm = renderDocumentDetailEdit({ formValues, moveDocSchema });
      const title = documentForm.find('[data-cy="title"]');
      const moveDocumentType = documentForm.find('[data-cy="move-document-type"]');
      const weightTicketSetType = documentForm.find('[data-cy="weight-ticket-set-type"]');
      const make = documentForm.find('[data-cy="vehicle-make"]');
      const model = documentForm.find('[data-cy="vehicle-model"]');
      const vehicleNickname = documentForm.find('[data-cy="vehicle-nickname"]');
      const emptyWeight = documentForm.find('[data-cy="empty-weight"]');
      const fullWeight = documentForm.find('[data-cy="full-weight"]');
      const status = documentForm.find('[data-cy="status"]');
      const notes = documentForm.find('[data-cy="notes"]');

      expect(title.length).toEqual(1);
      expect(moveDocumentType.length).toEqual(1);
      expect(weightTicketSetType.length).toEqual(1);
      expect(make.length).toEqual(0);
      expect(model.length).toEqual(0);
      expect(vehicleNickname.length).toEqual(1);
      expect(emptyWeight.length).toEqual(1);
      expect(fullWeight.length).toEqual(1);
      expect(status.length).toEqual(1);
      expect(notes.length).toEqual(1);
    });
    it('shows all form fields for progear', () => {
      const formValues = {
        moveDocument: {
          move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
          weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.PRO_GEAR,
        },
      };

      const documentForm = renderDocumentDetailEdit({ formValues, moveDocSchema });
      const title = documentForm.find('[data-cy="title"]');
      const moveDocumentType = documentForm.find('[data-cy="move-document-type"]');
      const weightTicketSetType = documentForm.find('[data-cy="weight-ticket-set-type"]');
      const make = documentForm.find('[data-cy="vehicle-make"]');
      const model = documentForm.find('[data-cy="vehicle-model"]');
      const progearType = documentForm.find('[data-cy="progear-type"]');
      const emptyWeight = documentForm.find('[data-cy="empty-weight"]');
      const fullWeight = documentForm.find('[data-cy="full-weight"]');
      const status = documentForm.find('[data-cy="status"]');
      const notes = documentForm.find('[data-cy="notes"]');

      expect(title.length).toEqual(1);
      expect(moveDocumentType.length).toEqual(1);
      expect(weightTicketSetType.length).toEqual(1);
      expect(make.length).toEqual(0);
      expect(model.length).toEqual(0);
      expect(progearType.length).toEqual(1);
      expect(emptyWeight.length).toEqual(1);
      expect(fullWeight.length).toEqual(1);
      expect(status.length).toEqual(1);
      expect(notes.length).toEqual(1);
    });
  });
  describe('expense document type', () => {
    it('shows all form fields for storage expense document type', () => {
      const formValues = {
        moveDocument: {
          move_document_type: MOVE_DOC_TYPE.EXPENSE,
          moving_expense_type: 'STORAGE',
        },
      };

      const documentForm = renderDocumentDetailEdit({ formValues, moveDocSchema });
      console.log(documentForm.find('ExpenseDocumentForm').debug());
      const title = documentForm.find('[data-cy="title"]');
      const moveDocumentType = documentForm.find('[data-cy="move-document-type"]');
      const storageStartDate = documentForm.find('[data-cy="storage-start-date"]');
      const storageEndDate = documentForm.find('[data-cy="storage-end-date"]');
      const status = documentForm.find('[data-cy="status"]');
      const notes = documentForm.find('[data-cy="notes"]');
      const expenseDocumentForm = documentForm.find('ExpenseDocumentForm');

      expect(title.length).toEqual(1);
      expect(moveDocumentType.length).toEqual(1);
      expect(storageStartDate.length).toEqual(1);
      expect(storageEndDate.length).toEqual(1);
      expect(status.length).toEqual(1);
      expect(notes.length).toEqual(1);
      expect(expenseDocumentForm.length).toEqual(1);
    });
  });
});
