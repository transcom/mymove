import React from 'react';
import { Field, reduxForm } from 'redux-form';
import './SignatureForm.css';
import validator from 'shared/JsonSchemaForm/validator';

function SignatureForm(props) {
  return (
    <div>
      <h3>SIGNATURE</h3>
      <p>I agree that I have read and understand the above notifications.</p>
      <div className="signing_box">
        <label id="name_field">
          Your name
          <Field
            name="signature"
            component="input"
            placeholder="Joseph Snuffy"
            type="text"
            validate={validator.isRequired}
          />
        </label>
        <label id="date_field">
          Today's date
          <Field name="date" component="input" type="text" readOnly />
        </label>
      </div>
    </div>
  );
}

export default reduxForm({
  form: 'certification_signature', // a unique identifier for this form
})(SignatureForm);
