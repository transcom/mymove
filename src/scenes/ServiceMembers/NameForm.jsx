import React from 'react';
import { Field, reduxForm } from 'redux-form';
import validator from 'shared/JsonSchemaForm/validator';

import './NameForm.css';

function NameForm(props) {
  return (
    <form>
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
        Middle name <span className="label-optional">Optional</span>
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
        Suffix <span className="label-optional">Optional</span>
        <Field name="suffix" component="input" type="text" />
      </label>
    </form>
  );
}

export default reduxForm({
  form: 'service_member_name',
})(NameForm);
