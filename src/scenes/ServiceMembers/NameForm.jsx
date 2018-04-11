import React from 'react';
import { Field, reduxForm } from 'redux-form';
import validator from 'shared/JsonSchemaForm/validator';

function NameForm(props) {
  const { handleSubmit, pristine, invalid, submitting } = props;
  return (
    <form onSubmit={handleSubmit}>
      <h2>Name</h2>
      <label id="firstName">
        First name
        <Field
          name="first_name"
          component="input"
          type="text"
          validate={validator.isRequired}
        />
      </label>
      <label id="middleInitial">
        Middle name Optional
        <Field name="middle_initial" component="input" type="text" />
      </label>
      <label id="lastName">
        Last name
        <Field
          name="last_name"
          component="input"
          type="text"
          validate={validator.isRequired}
        />
      </label>
      <label id="suffix">
        Suffix Optional
        <Field name="suffix" component="input" type="text" />
      </label>
      <button type="submit" disabled={pristine || submitting || invalid}>
        Submit
      </button>
    </form>
  );
}

export default reduxForm({
  form: 'service_member_name',
})(NameForm);
