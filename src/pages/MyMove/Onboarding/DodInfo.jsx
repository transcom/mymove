/* eslint-disable react/forbid-prop-types */
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { get, pick } from 'lodash';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues, Field } from 'redux-form';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { normalizeSSN } from 'shared/JsonSchemaForm/reduxFieldNormalizer';
import SSNField from 'components/form/fields/SSNInput';

const subsetOfFields = ['affiliation', 'edipi', 'social_security_number', 'rank'];

const updateServiceMemberAction = () => {};

const validateDodForm = (values, form) => {
  // Everything is taken care of except for SSN
  const errors = {};
  const ssn = values.social_security_number;
  const hasSSN = form.ssnOnServer;

  const validSSNPattern = RegExp('^\\d{3}-\\d{2}-\\d{4}$');
  const validSSN = validSSNPattern.test(ssn);
  const ssnPresent = ssn !== '' && ssn !== undefined;

  if (hasSSN) {
    if (ssnPresent && !validSSN) {
      errors.social_security_number = 'SSN must have 9 digits';
    }
  } else if (!ssnPresent) {
    errors.social_security_number = 'Required';
  } else if (!validSSN) {
    errors.social_security_number = 'SSN must have 9 digits';
  }

  return errors;
};

const formName = 'service_member_dod_info';
const DodWizardForm = reduxifyWizardForm(formName, validateDodForm);

export class DodInfo extends Component {
  handleSubmit = () => {
    const { values, updateServiceMember } = this.props;
    const pendingValues = values;
    if (pendingValues) {
      const patch = pick(pendingValues, subsetOfFields);
      return updateServiceMember(patch);
    }

    return null;
  };

  render() {
    const { pages, pageKey, error, currentServiceMember, schema } = this.props;
    const initialValues = currentServiceMember ? pick(currentServiceMember, subsetOfFields) : null;

    const ssnOnServer = currentServiceMember ? currentServiceMember.has_social_security_number : false;

    return (
      <DodWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error}
        initialValues={initialValues}
        ssnOnServer={ssnOnServer}
      >
        <div className="grid-row">
          <div className="grid-col-12">
            <h1 className="sm-heading">Create your profile</h1>
            <p>Before we can schedule your move, we need to know a little more about you.</p>
            <SwaggerField fieldName="affiliation" swagger={schema} required />
            <SwaggerField fieldName="edipi" swagger={schema} required />
            <Field
              name="social_security_number"
              component={SSNField}
              ssnOnServer={ssnOnServer}
              normalize={normalizeSSN}
            />
            <SwaggerField fieldName="rank" swagger={schema} required />
          </div>
        </div>
      </DodWizardForm>
    );
  }
}

DodInfo.propTypes = {
  pageKey: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  values: PropTypes.object,
  pages: PropTypes.arrayOf(PropTypes.string).isRequired,
};

DodInfo.defaultProps = {
  currentServiceMember: null,
  error: null,
  values: null,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateServiceMember: updateServiceMemberAction }, dispatch);
}
function mapStateToProps(state) {
  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    values: getFormValues(formName)(state),
    ...state.serviceMember,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(DodInfo);
