import React, { Fragment } from 'react';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

let AccessCode = props => {
  const { schema } = props;
  return (
    <Fragment>
      <div className="usa-grid">
        <h3 className="title">Welcome to MilMove</h3>
        <p>Please enter your MilMove access code in the field below.</p>
        <SwaggerField fieldName="claim_access_code" swagger={schema} required />
        <button className="usa-button-secondary">Continue</button>
        <p>No code? Go to DPS to schedule your move.</p>
      </div>
    </Fragment>
  );
};
const formName = 'claim_access_code_field';
AccessCode = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(AccessCode);

AccessCode.propTypes = {
  schema: PropTypes.object.isRequired,
};

function mapStateToProps(state) {
  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.ClaimAccessCodePayload', {}),
  };
  return props;
}

export default connect(mapStateToProps)(AccessCode);
