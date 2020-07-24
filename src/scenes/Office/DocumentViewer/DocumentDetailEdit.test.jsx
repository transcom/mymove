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

      const title = documentForm.find('[data-testid="title"]');
      const moveDocumentType = documentForm.find('[data-testid="move-document-type"]');
      const weightTicketSetType = documentForm.find('[data-testid="weight-ticket-set-type"]');
      const make = documentForm.find('[data-testid="vehicle-make"]');
      const model = documentForm.find('[data-testid="vehicle-model"]');
      const vehicleNickname = documentForm.find('[data-testid="vehicle-nickname"]');
      const status = documentForm.find('[data-testid="status"]');
      const notes = documentForm.find('[data-testid="notes"]');

      expect(title.props()).toHaveProperty('fieldName', 'title');
      expect(moveDocumentType.props()).toHaveProperty('fieldName', 'move_document_type');
      expect(weightTicketSetType.props()).toHaveProperty('fieldName', 'weight_ticket_set_type');
      expect(make.props()).toHaveProperty('fieldName', 'vehicle_make');
      expect(model.props()).toHaveProperty('fieldName', 'vehicle_model');
      expect(vehicleNickname.length).toBeFalsy();
      expect(status.props()).toHaveProperty('fieldName', 'status');
      expect(notes.props()).toHaveProperty('fieldName', 'notes');
    });

    it('shows all form fields for a car+trailer', () => {
      const formValues = {
        moveDocument: {
          move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
          weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.CAR_TRAILER,
        },
      };

      const documentForm = renderDocumentDetailEdit({ formValues, moveDocSchema });

      const title = documentForm.find('[data-testid="title"]');
      const moveDocumentType = documentForm.find('[data-testid="move-document-type"]');
      const weightTicketSetType = documentForm.find('[data-testid="weight-ticket-set-type"]');
      const make = documentForm.find('[data-testid="vehicle-make"]');
      const model = documentForm.find('[data-testid="vehicle-model"]');
      const vehicleNickname = documentForm.find('[data-testid="vehicle-nickname"]');
      const status = documentForm.find('[data-testid="status"]');
      const notes = documentForm.find('[data-testid="notes"]');

      expect(title.props()).toHaveProperty('fieldName', 'title');
      expect(moveDocumentType.props()).toHaveProperty('fieldName', 'move_document_type');
      expect(weightTicketSetType.props()).toHaveProperty('fieldName', 'weight_ticket_set_type');
      expect(make.props()).toHaveProperty('fieldName', 'vehicle_make');
      expect(model.props()).toHaveProperty('fieldName', 'vehicle_model');
      expect(vehicleNickname.length).toBeFalsy();
      expect(status.props()).toHaveProperty('fieldName', 'status');
      expect(notes.props()).toHaveProperty('fieldName', 'notes');
    });

    it('shows all form fields for a boxtruck', () => {
      const formValues = {
        moveDocument: {
          move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
          weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.BOX_TRUCK,
        },
      };

      const documentForm = renderDocumentDetailEdit({ formValues, moveDocSchema });
      const title = documentForm.find('[data-testid="title"]');
      const moveDocumentType = documentForm.find('[data-testid="move-document-type"]');
      const weightTicketSetType = documentForm.find('[data-testid="weight-ticket-set-type"]');
      const make = documentForm.find('[data-testid="vehicle-make"]');
      const model = documentForm.find('[data-testid="vehicle-model"]');
      const vehicleNickname = documentForm.find('[data-testid="vehicle-nickname"]');
      const emptyWeight = documentForm.find('[data-testid="empty-weight"]');
      const fullWeight = documentForm.find('[data-testid="full-weight"]');
      const status = documentForm.find('[data-testid="status"]');
      const notes = documentForm.find('[data-testid="notes"]');

      expect(title.props()).toHaveProperty('fieldName', 'title');
      expect(moveDocumentType.props()).toHaveProperty('fieldName', 'move_document_type');
      expect(weightTicketSetType.props()).toHaveProperty('fieldName', 'weight_ticket_set_type');
      expect(make.length).toBeFalsy();
      expect(model.length).toBeFalsy();
      expect(vehicleNickname.props()).toHaveProperty('fieldName', 'vehicle_nickname');
      expect(emptyWeight.props()).toHaveProperty('fieldName', 'empty_weight');
      expect(fullWeight.props()).toHaveProperty('fieldName', 'full_weight');
      expect(status.props()).toHaveProperty('fieldName', 'status');
      expect(notes.props()).toHaveProperty('fieldName', 'notes');
    });
    it('shows all form fields for progear', () => {
      const formValues = {
        moveDocument: {
          move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET,
          weight_ticket_set_type: WEIGHT_TICKET_SET_TYPE.PRO_GEAR,
        },
      };

      const documentForm = renderDocumentDetailEdit({ formValues, moveDocSchema });
      const title = documentForm.find('[data-testid="title"]');
      const moveDocumentType = documentForm.find('[data-testid="move-document-type"]');
      const weightTicketSetType = documentForm.find('[data-testid="weight-ticket-set-type"]');
      const make = documentForm.find('[data-testid="vehicle-make"]');
      const model = documentForm.find('[data-testid="vehicle-model"]');
      const progearType = documentForm.find('[data-testid="progear-type"]');
      const emptyWeight = documentForm.find('[data-testid="empty-weight"]');
      const fullWeight = documentForm.find('[data-testid="full-weight"]');
      const status = documentForm.find('[data-testid="status"]');
      const notes = documentForm.find('[data-testid="notes"]');

      expect(title.props()).toHaveProperty('fieldName', 'title');
      expect(moveDocumentType.props()).toHaveProperty('fieldName', 'move_document_type');
      expect(weightTicketSetType.props()).toHaveProperty('fieldName', 'weight_ticket_set_type');
      expect(make.length).toBeFalsy();
      expect(model.length).toBeFalsy();
      expect(progearType.props()).toHaveProperty('fieldName', 'vehicle_nickname');
      expect(emptyWeight.props()).toHaveProperty('fieldName', 'empty_weight');
      expect(fullWeight.props()).toHaveProperty('fieldName', 'full_weight');
      expect(status.props()).toHaveProperty('fieldName', 'status');
      expect(notes.props()).toHaveProperty('fieldName', 'notes');
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
      const title = documentForm.find('[data-testid="title"]');
      const moveDocumentType = documentForm.find('[data-testid="move-document-type"]');
      const storageStartDate = documentForm.find('[data-testid="storage-start-date"]');
      const storageEndDate = documentForm.find('[data-testid="storage-end-date"]');
      const status = documentForm.find('[data-testid="status"]');
      const notes = documentForm.find('[data-testid="notes"]');
      const expenseDocumentForm = documentForm.find('ExpenseDocumentForm');

      expect(title.props()).toHaveProperty('fieldName', 'title');
      expect(moveDocumentType.props()).toHaveProperty('fieldName', 'move_document_type');
      expect(storageStartDate.props()).toHaveProperty('fieldName', 'storage_start_date');
      expect(storageEndDate.props()).toHaveProperty('fieldName', 'storage_end_date');
      expect(status.props()).toHaveProperty('fieldName', 'status');
      expect(notes.props()).toHaveProperty('fieldName', 'notes');
      expect(expenseDocumentForm.props()).toBeTruthy();
    });
  });
});
