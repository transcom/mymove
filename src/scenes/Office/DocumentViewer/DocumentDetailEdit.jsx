import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { get, isEmpty } from 'lodash';
import { FormSection } from 'redux-form';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { MOVE_DOC_TYPE, WEIGHT_TICKET_SET_TYPE } from 'shared/constants';

import ExpenseDocumentForm from 'scenes/Office/DocumentViewer/ExpenseDocumentForm';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

const DocumentDetailEdit = ({ formValues, moveDocSchema }) => {
  const isExpenseDocument = get(formValues.moveDocument, 'move_document_type') === MOVE_DOC_TYPE.EXPENSE;
  const isWeightTicketDocument = get(formValues.moveDocument, 'move_document_type') === MOVE_DOC_TYPE.WEIGHT_TICKET_SET;
  const isStorageExpenseDocument =
    get(formValues.moveDocument, 'move_document_type') === 'EXPENSE' &&
    get(formValues.moveDocument, 'moving_expense_type') === 'STORAGE';
  const isWeightTicketTypeCarOrTrailer =
    isWeightTicketDocument &&
    (formValues.moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.CAR ||
      formValues.moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.CAR_TRAILER);
  const isWeightTicketTypeBoxTruck =
    isWeightTicketDocument && formValues.moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.BOX_TRUCK;
  const isWeightTicketTypeProGear =
    isWeightTicketDocument && formValues.moveDocument.weight_ticket_set_type === WEIGHT_TICKET_SET_TYPE.PRO_GEAR;

  return isEmpty(formValues.moveDocument) ? (
    <LoadingPlaceholder />
  ) : (
    <Fragment>
      <div>
        <FormSection name="moveDocument">
          <SwaggerField fieldName="title" swagger={moveDocSchema} required />
          <SwaggerField data-cy="move-document-type" fieldName="move_document_type" swagger={moveDocSchema} required />
          {isExpenseDocument && <ExpenseDocumentForm moveDocSchema={moveDocSchema} />}
          {isWeightTicketDocument && (
            <>
              <div>
                <SwaggerField fieldName="weight_ticket_set_type" swagger={moveDocSchema} required />
              </div>
              {isWeightTicketTypeBoxTruck && (
                <SwaggerField fieldName="vehicle_nickname" swagger={moveDocSchema} required />
              )}
              {isWeightTicketTypeProGear && (
                <SwaggerField
                  fieldName="vehicle_nickname"
                  title="Pro-gear type (ex. 'My pro-gear', 'Spouse pro-gear', 'Both')"
                  swagger={moveDocSchema}
                  required
                />
              )}
              {isWeightTicketTypeCarOrTrailer && (
                <>
                  <SwaggerField fieldName="vehicle_make" swagger={moveDocSchema} required />
                  <SwaggerField fieldName="vehicle_model" swagger={moveDocSchema} required />
                </>
              )}
              <SwaggerField className="short-field" fieldName="empty_weight" swagger={moveDocSchema} required />{' '}
              <span className="field-with-units">lbs</span>
              <SwaggerField className="short-field" fieldName="full_weight" swagger={moveDocSchema} required />{' '}
              <span className="field-with-units">lbs</span>
            </>
          )}
          {isStorageExpenseDocument && (
            <>
              <SwaggerField title="Start date" fieldName="storage_start_date" required swagger={moveDocSchema} />
              <SwaggerField title="End date" fieldName="storage_end_date" required swagger={moveDocSchema} />
            </>
          )}
          <SwaggerField data-cy="status" fieldName="status" swagger={moveDocSchema} required />
          <SwaggerField data-cy="notes" fieldName="notes" swagger={moveDocSchema} />
        </FormSection>
      </div>
    </Fragment>
  );
};
const { bool, object, shape, string, arrayOf } = PropTypes;

DocumentDetailEdit.propTypes = {
  isExpenseDocument: bool.isRequired,
  isWeightTicketDocument: bool.isRequired,
  moveDocSchema: shape({
    properties: object.isRequired,
    required: arrayOf(string).isRequired,
    type: string.isRequired,
  }).isRequired,
};

export default DocumentDetailEdit;
