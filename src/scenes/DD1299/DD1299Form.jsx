import React from 'react';
import { Field, FormSection, reduxForm } from 'redux-form';
import { getFields } from './fields';

import './DD1299.css';

const renderGroupOrField = (name, fields) => {
  const field = fields[name];
  if (field.type !== 'group') return renderField(name, fields);
  return (
    <section key={name}>
      <h4 htmlFor={name}>{field.label}</h4>
      {Object.keys(field.fields).map(f => renderField(f, field.fields))}
    </section>
  );
};
const renderField = (name, fields) => {
  const field = fields[name];
  return (
    <div key={name}>
      <label htmlFor={name}>
        {field.label || name}
        <Field
          name={name}
          component="input"
          type={field.format}
          placeholder={field.example}
        />
      </label>
    </div>
  );
};
const DD1299 = props => {
  const { handleSubmit } = props;
  const fields = getFields();
  return (
    <form onSubmit={handleSubmit}>
      {Object.keys(fields).map(f => renderGroupOrField(f, fields))}

      <button type="submit">Submit</button>
    </form>
  );
};

export default reduxForm({
  // a unique name for the form
  form: 'DD1299',
})(DD1299);
