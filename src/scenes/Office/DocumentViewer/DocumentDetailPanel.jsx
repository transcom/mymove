import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, omit, cloneDeep } from 'lodash';
import { reduxForm, getFormValues } from 'redux-form';

import { convertDollarsToCents } from 'shared/utils';
import { formatCents } from 'utils/formatters';
import { editablePanelify } from 'shared/EditablePanel';
import { selectMoveDocument, updateMoveDocument } from 'shared/Entities/modules/moveDocuments';
import { selectActivePPMForMove } from 'shared/Entities/modules/ppms';
import { isMovingExpenseDocument } from 'shared/Entities/modules/movingExpenseDocuments';

import DocumentDetailDisplay from './DocumentDetailDisplay';
import DocumentDetailEdit from './DocumentDetailEdit';

const formName = 'move_document_viewer';

let DocumentDetailPanel = editablePanelify(DocumentDetailDisplay, DocumentDetailEdit);

DocumentDetailPanel = reduxForm({ form: formName })(DocumentDetailPanel);

function mapStateToProps(state, props) {
  const { moveId, moveDocumentId } = props;
  const moveDocument = selectMoveDocument(state, moveDocumentId);
  const isExpenseDocument = isMovingExpenseDocument(moveDocument);
  const isWeightTicketDocument = get(moveDocument, 'move_document_type') === 'WEIGHT_TICKET_SET';
  const isStorageExpenseDocument =
    get(moveDocument, 'move_document_type') === 'EXPENSE' && get(moveDocument, 'moving_expense_type') === 'STORAGE';
  // Convert cents to collars - make a deep clone copy to not modify moveDocument itself
  const initialMoveDocument = cloneDeep(moveDocument);
  const requested_amount = get(initialMoveDocument, 'requested_amount_cents');
  if (requested_amount) {
    initialMoveDocument.requested_amount_cents = formatCents(requested_amount);
  }

  return {
    // reduxForm
    initialValues: {
      moveDocument: initialMoveDocument,
    },
    isExpenseDocument,
    isWeightTicketDocument,
    isStorageExpenseDocument,
    formValues: getFormValues(formName)(state),
    moveDocSchema: get(state, 'swaggerInternal.spec.definitions.MoveDocumentPayload', {}),
    hasError: false,
    isUpdating: false,
    moveDocument,

    // editablePanelify
    getUpdateArgs: function () {
      // Make a copy of values to not modify moveDocument
      let values = cloneDeep(getFormValues(formName)(state));
      values.moveDocument.personally_procured_move_id = selectActivePPMForMove(state, props.moveId).id;
      if (
        get(values.moveDocument, 'move_document_type', '') !== 'EXPENSE' &&
        get(values.moveDocument, 'payment_method', false)
      ) {
        values.moveDocument = omit(values.moveDocument, ['payment_method', 'requested_amount_cents']);
      }
      if (get(values.moveDocument, 'move_document_type', '') === 'EXPENSE') {
        values.moveDocument.requested_amount_cents = convertDollarsToCents(values.moveDocument.requested_amount_cents);
      }
      return [moveId, moveDocumentId, values.moveDocument];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updateMoveDocument,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(DocumentDetailPanel);
