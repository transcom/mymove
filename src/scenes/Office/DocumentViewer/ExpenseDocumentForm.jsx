import React, { Fragment } from 'react';
import { FormSection } from 'redux-form';
import PropTypes from 'prop-types';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const ExpenseDocumentForm = props => {
  const movingExpenseDocumentSchema = props.movingExpenseDocumentSchema;
  const reimbursementSchema = props.reimbursementSchema;

  return (
    <Fragment>
      <FormSection name="movingExpenseDocument">
        <SwaggerField
          title="Expense type"
          fieldName="moving_expense_type"
          swagger={movingExpenseDocumentSchema}
          required
        />
      </FormSection>
      <FormSection name="reimbursement">
        <SwaggerField
          title="Amount"
          fieldName="requested_amount"
          swagger={reimbursementSchema}
          required
        />
        <SwaggerField
          title="Method of Payment"
          fieldName="method_of_receipt"
          swagger={reimbursementSchema}
          required
        />
      </FormSection>
    </Fragment>
  );
};
ExpenseDocumentForm.propTypes = {
  movingExpenseDocumentSchema: PropTypes.object,
  reimbursementSchema: PropTypes.object,
};
export default ExpenseDocumentForm;
