import React, { Fragment } from 'react';
import { FormSection } from 'redux-form';
import PropTypes from 'prop-types';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const ExpenseDocumentForm = props => {
  const moveDocSchema = props.moveDocSchema;
  const reimbursementSchema = props.reimbursementSchema;
  return (
    <Fragment>
      <SwaggerField
        title="Expense type"
        fieldName="moving_expense_type"
        swagger={moveDocSchema}
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
  moveDocSchema: PropTypes.object,
  reimbursementSchema: PropTypes.object,
};
export default ExpenseDocumentForm;
