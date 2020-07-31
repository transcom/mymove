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
          <SwaggerField data-testid="title" fieldName="title" swagger={moveDocSchema} required />
          <SwaggerField
            data-testid="move-document-type"
            fieldName="move_document_type"
            swagger={moveDocSchema}
            required
          />
          {isExpenseDocument && <ExpenseDocumentForm moveDocSchema={moveDocSchema} />}
          {isWeightTicketDocument && (
            <>
              <div>
                <SwaggerField
                  data-testid="weight-ticket-set-type"
                  fieldName="weight_ticket_set_type"
                  swagger={moveDocSchema}
                  required
                />
              </div>
              {isWeightTicketTypeBoxTruck && (
                <SwaggerField
                  data-testid="vehicle-nickname"
                  fieldName="vehicle_nickname"
                  swagger={moveDocSchema}
                  required
                />
              )}
              {isWeightTicketTypeProGear && (
                <SwaggerField
                  data-testid="progear-type"
                  fieldName="vehicle_nickname"
                  title="Pro-gear type (ex. 'My pro-gear', 'Spouse pro-gear', 'Both')"
                  swagger={moveDocSchema}
                  required
                />
              )}
              {isWeightTicketTypeCarOrTrailer && (
                <>
                  <SwaggerField data-testid="vehicle-make" fieldName="vehicle_make" swagger={moveDocSchema} required />
                  <SwaggerField
                    data-testid="vehicle-model"
                    fieldName="vehicle_model"
                    swagger={moveDocSchema}
                    required
                  />
                </>
              )}
              <SwaggerField
                data-testid="empty-weight"
                className="short-field"
                title="Empty weight"
                fieldName="empty_weight"
                swagger={moveDocSchema}
                required
              />{' '}
              <span className="field-with-units">lbs</span>
              <SwaggerField
                data-testid="full-weight"
                className="short-field"
                title="Full weight"
                fieldName="full_weight"
                swagger={moveDocSchema}
                required
              />{' '}
              <span className="field-with-units">lbs</span>
            </>
          )}
          {isStorageExpenseDocument && (
            <>
              <SwaggerField
                data-testid="storage-start-date"
                title="Start date"
                fieldName="storage_start_date"
                required
                swagger={moveDocSchema}
              />
              <SwaggerField
                data-testid="storage-end-date"
                title="End date"
                fieldName="storage_end_date"
                required
                swagger={moveDocSchema}
              />
            </>
          )}
          <SwaggerField
            data-testid="status"
            label="Document status"
            fieldName="status"
            swagger={moveDocSchema}
            required
          />
          <SwaggerField data-testid="notes" fieldName="notes" swagger={moveDocSchema} />
        </FormSection>
      </div>
    </Fragment>
  );
};
const { object, shape, string, arrayOf } = PropTypes;

DocumentDetailEdit.propTypes = {
  moveDocSchema: shape({
    properties: object.isRequired,
    required: arrayOf(string).isRequired,
    type: string.isRequired,
  }).isRequired,
};

export default DocumentDetailEdit;
