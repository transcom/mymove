import React from 'react';
import { Field, reduxForm } from 'redux-form';
import { getFields } from './fields';

import './DD1299.css';

const renderGroupOrField = (name, fields) => {
  /*TODO:
   handle enums
   telephone numbers
   textbox vs textarea (e.g for addresses)
   dates look wonky
   styling in accordance with USWDS
   formatting of labels?
  */
  const field = fields[name];
  if (field.type !== 'group') return renderField(name, fields);
  return (
    <fieldset key={name}>
      <legend htmlFor={name}>{field.label}</legend>
      {Object.keys(field.fields).map(f => renderGroupOrField(f, field.fields))}
    </fieldset>
  );
};
const renderField = (name, fields) => {
  const field = fields[name];
  return (
    <div key={name}>
      <label htmlFor={name}>
        {field.label || name}
        <Field name={name} component="input" type={field.format} />
      </label>
    </div>
  );
};
const DD1299 = props => {
  const { handleSubmit } = props;
  const fields = getFields();
  return (
    <form className="dd1299" onSubmit={handleSubmit}>
      <h1>Application For Shipment And Or Storage Of Personal Property</h1>
      {Object.keys(fields).map(f => renderGroupOrField(f, fields))}

      <button type="submit">Submit</button>
    </form>
  );
};

export default reduxForm({
  // a unique name for the form
  form: 'DD1299',
})(DD1299);
