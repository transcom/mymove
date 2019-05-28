import React, { Fragment } from 'react';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';
import { connect } from 'react-redux';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import './AccessCode.css';
import { validateAccessCode, claimAccessCode } from 'shared/Entities/modules/accessCodes';

class AccessCode extends React.Component {
  claimAccessCode = () => {
    const { serviceMemberId, formValues, validateAccessCode } = this.props;
    const accessCode = formValues.claim_access_code;
    validateAccessCode(accessCode)
      .then(response => {
        if (!response.valid) {
          // propagate error to reduxForm
          console.log('Form is not valid');
          return;
        }
        return claimAccessCode(accessCode, serviceMemberId);
      })
      .catch(() => {});
  };

  render() {
    const { schema } = this.props;
    return (
      <Fragment>
        <div className="usa-grid">
          <h3 className="title">Welcome to MilMove</h3>
          <p>Please enter your MilMove access code in the field below.</p>
          <SwaggerField fieldName="claim_access_code" swagger={schema} required />
          <button className="usa-button-primary" onClick={this.claimAccessCode}>
            Continue
          </button>
          <br />No code? Go to DPS to schedule your move.
        </div>
      </Fragment>
    );
  }
}

const formName = 'claim_access_code_form';
AccessCode = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(AccessCode);

AccessCode.propTypes = {
  schema: PropTypes.object.isRequired,
  serviceMemberId: PropTypes.string.isRequired,
};

function mapStateToProps(state) {
  const serviceMember = get(state, 'serviceMember.currentServiceMember');
  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.ClaimAccessCodePayload', {}),
    serviceMemberId: get(serviceMember, 'id'),
    formValues: getFormValues(formName)(state),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ validateAccessCode }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(AccessCode);
