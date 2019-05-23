import React, { Fragment } from 'react';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import './AccessCode.css';

import { createServiceMember } from 'scenes/ServiceMembers/ducks';

let AccessCode = props => {
  const { history, schema, createServiceMember } = props;
  return (
    <Fragment>
      <div className="usa-grid">
        <h3 className="title">Welcome to MilMove</h3>
        <p>Please enter your MilMove access code in the field below.</p>
        <SwaggerField fieldName="claim_access_code" swagger={schema} required />
        <button
          className="usa-button-primary"
          onClick={() => {
            createServiceMember({}).then(response => {
              const serviceMemberId = response.payload.id;
              history.push(`/service-member/${serviceMemberId}/create`);
            });
          }}
        >
          Continue
        </button>
        <br />No code? Go to DPS to schedule your move.
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

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createServiceMember }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(AccessCode);
