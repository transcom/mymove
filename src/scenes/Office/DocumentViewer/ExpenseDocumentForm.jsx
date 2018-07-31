import React, { Fragment } from 'react';
import { FormSection } from 'redux-form';
import PropTypes from 'prop-types';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const ExpenseDocumentForm = props => {
  const movingExpenseSchema = props.movingExpenseSchema;
  // const genericMoveDocSchema = props.genericMoveDocSchema;
  const reimbursementSchema = props.reimbursementSchema;
  console.log('expense doc schema', movingExpenseSchema);
  return (
    <Fragment>
      <SwaggerField
        title="Expense type"
        fieldName="moving_expense_type"
        swagger={movingExpenseSchema}
        required
      />
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
  // genericMoveDocSchema: PropTypes.object,
  movingExpenseSchema: PropTypes.object,
  reimbursementSchema: PropTypes.object,
};
export default ExpenseDocumentForm;
