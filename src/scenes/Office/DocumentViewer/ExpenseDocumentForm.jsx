import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const ExpenseDocumentForm = props => {
  const moveDocSchema = props.moveDocSchema;
  return (
    <Fragment>
      <SwaggerField title="Expense type" fieldName="moving_expense_type" swagger={moveDocSchema} required />
      <SwaggerField title="Amount" fieldName="requested_amount_cents" swagger={moveDocSchema} required />
      <SwaggerField title="Payment Method" fieldName="payment_method" swagger={moveDocSchema} required />
    </Fragment>
  );
};
ExpenseDocumentForm.propTypes = {
  moveDocSchema: PropTypes.object,
};
export default ExpenseDocumentForm;
